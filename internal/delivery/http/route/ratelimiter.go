package route

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// redisStorage implements fiber.Storage using an existing go-redis client.
type redisStorage struct {
	rdb *redis.Client
}

func (s *redisStorage) Get(key string) ([]byte, error) {
	val, err := s.rdb.Get(context.Background(), key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return val, err
}

func (s *redisStorage) Set(key string, val []byte, exp time.Duration) error {
	return s.rdb.Set(context.Background(), key, val, exp).Err()
}

func (s *redisStorage) Delete(key string) error {
	return s.rdb.Del(context.Background(), key).Err()
}

// Reset is a no-op — flushing a shared Redis instance is dangerous.
func (s *redisStorage) Reset() error { return nil }

// Close is a no-op — client lifecycle is managed externally.
func (s *redisStorage) Close() error { return nil }

// memoryEntry holds a value with its expiry time.
type memoryEntry struct {
	val []byte
	exp time.Time
}

// memoryStorage implements fiber.Storage using a sync.Map with lazy expiry.
type memoryStorage struct {
	m sync.Map
}

func (s *memoryStorage) Get(key string) ([]byte, error) {
	v, ok := s.m.Load(key)
	if !ok {
		return nil, nil
	}
	entry := v.(memoryEntry)
	if !entry.exp.IsZero() && time.Now().After(entry.exp) {
		s.m.Delete(key)
		return nil, nil
	}
	return entry.val, nil
}

func (s *memoryStorage) Set(key string, val []byte, exp time.Duration) error {
	var expTime time.Time
	if exp > 0 {
		expTime = time.Now().Add(exp)
	}
	s.m.Store(key, memoryEntry{val: val, exp: expTime})
	return nil
}

func (s *memoryStorage) Delete(key string) error {
	s.m.Delete(key)
	return nil
}

func (s *memoryStorage) Reset() error {
	s.m.Range(func(k, _ any) bool {
		s.m.Delete(k)
		return true
	})
	return nil
}

func (s *memoryStorage) Close() error { return nil }

// FallbackStorage tries Redis first; on repeated failure it opens a circuit and
// uses in-memory storage until Redis recovers.
type FallbackStorage struct {
	primary  *redisStorage
	fallback *memoryStorage
	log      *logrus.Logger

	// circuit breaker state (all accessed atomically)
	failures    atomic.Int32
	circuitOpen atomic.Bool
	openedAt    atomic.Int64 // unix seconds

	maxFailures int32
	cooldown    time.Duration
}

// NewFallbackStorage creates a FallbackStorage. The circuit opens after 3
// consecutive Redis failures and re-probes Redis every 10 seconds.
func NewFallbackStorage(rdb *redis.Client, log *logrus.Logger) *FallbackStorage {
	return &FallbackStorage{
		primary:     &redisStorage{rdb: rdb},
		fallback:    &memoryStorage{},
		log:         log,
		maxFailures: 3,
		cooldown:    10 * time.Second,
	}
}

// useFallback returns true when the circuit is open and the cooldown has not
// yet elapsed. When the cooldown has elapsed it allows one probe attempt
// by returning false.
func (s *FallbackStorage) useFallback() bool {
	if !s.circuitOpen.Load() {
		return false
	}
	elapsed := time.Since(time.Unix(s.openedAt.Load(), 0))
	return elapsed < s.cooldown
}

// recordSuccess closes the circuit and resets the failure counter.
func (s *FallbackStorage) recordSuccess() {
	if s.circuitOpen.Swap(false) {
		s.log.Info("rate limiter: Redis recovered, circuit closed")
	}
	s.failures.Store(0)
}

// recordFailure increments the failure counter and opens the circuit when the
// threshold is reached.
func (s *FallbackStorage) recordFailure(err error) {
	n := s.failures.Add(1)
	if n >= s.maxFailures && !s.circuitOpen.Swap(true) {
		s.openedAt.Store(time.Now().Unix())
		s.log.Warnf("rate limiter: Redis unavailable (%v), circuit opened — falling back to in-memory for %s", err, s.cooldown)
	}
}

func (s *FallbackStorage) Get(key string) ([]byte, error) {
	if s.useFallback() {
		return s.fallback.Get(key)
	}
	val, err := s.primary.Get(key)
	if err != nil {
		s.recordFailure(err)
		return s.fallback.Get(key)
	}
	s.recordSuccess()
	return val, nil
}

func (s *FallbackStorage) Set(key string, val []byte, exp time.Duration) error {
	if s.useFallback() {
		return s.fallback.Set(key, val, exp)
	}
	if err := s.primary.Set(key, val, exp); err != nil {
		s.recordFailure(err)
		return s.fallback.Set(key, val, exp)
	}
	s.recordSuccess()
	return nil
}

func (s *FallbackStorage) Delete(key string) error {
	if s.useFallback() {
		return s.fallback.Delete(key)
	}
	if err := s.primary.Delete(key); err != nil {
		s.recordFailure(err)
		return s.fallback.Delete(key)
	}
	s.recordSuccess()
	return nil
}

func (s *FallbackStorage) Reset() error {
	_ = s.primary.Reset()
	return s.fallback.Reset()
}

func (s *FallbackStorage) Close() error { return nil }
