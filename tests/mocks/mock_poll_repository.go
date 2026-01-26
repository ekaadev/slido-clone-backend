package mocks

import (
	"slido-clone-backend/internal/entity"

	"gorm.io/gorm"
)

// MockPollRepository mock implementation of PollRepository
type MockPollRepository struct {
	Polls      map[uint]*entity.Poll
	NextID     uint
	ShouldFail bool
}

// NewMockPollRepository create new mock
func NewMockPollRepository() *MockPollRepository {
	return &MockPollRepository{
		Polls:  make(map[uint]*entity.Poll),
		NextID: 1,
	}
}

// Create mock create
func (r *MockPollRepository) Create(db *gorm.DB, poll *entity.Poll) error {
	if r.ShouldFail {
		return gorm.ErrInvalidDB
	}
	poll.ID = r.NextID
	r.NextID++
	r.Polls[poll.ID] = poll
	return nil
}

// FindByIdWithOptions mock find by id
func (r *MockPollRepository) FindByIdWithOptions(db *gorm.DB, id uint) (*entity.Poll, error) {
	if r.ShouldFail {
		return nil, gorm.ErrInvalidDB
	}
	poll, ok := r.Polls[id]
	if !ok {
		return nil, nil
	}
	return poll, nil
}

// FindActiveByRoomID mock find active by room
func (r *MockPollRepository) FindActiveByRoomID(db *gorm.DB, roomID uint) ([]entity.Poll, error) {
	if r.ShouldFail {
		return nil, gorm.ErrInvalidDB
	}
	var polls []entity.Poll
	for _, poll := range r.Polls {
		if poll.RoomID == roomID && poll.Status == "active" {
			polls = append(polls, *poll)
		}
	}
	return polls, nil
}

// MockPollOptionRepository mock implementation
type MockPollOptionRepository struct {
	Options    map[uint]*entity.PollOption
	NextID     uint
	ShouldFail bool
}

// NewMockPollOptionRepository create new mock
func NewMockPollOptionRepository() *MockPollOptionRepository {
	return &MockPollOptionRepository{
		Options: make(map[uint]*entity.PollOption),
		NextID:  1,
	}
}

// CreateBatch mock create batch
func (r *MockPollOptionRepository) CreateBatch(db *gorm.DB, options []entity.PollOption) error {
	if r.ShouldFail {
		return gorm.ErrInvalidDB
	}
	for i := range options {
		options[i].ID = r.NextID
		r.NextID++
		r.Options[options[i].ID] = &options[i]
	}
	return nil
}

// GetTotalVotesByPollID mock get total votes
func (r *MockPollOptionRepository) GetTotalVotesByPollID(db *gorm.DB, pollID uint) (int, error) {
	if r.ShouldFail {
		return 0, gorm.ErrInvalidDB
	}
	total := 0
	for _, opt := range r.Options {
		if opt.PollID == pollID {
			total += opt.VoteCount
		}
	}
	return total, nil
}

// ValidateOptionBelongsToPoll mock validate
func (r *MockPollOptionRepository) ValidateOptionBelongsToPoll(db *gorm.DB, optionID uint, pollID uint) (bool, error) {
	if r.ShouldFail {
		return false, gorm.ErrInvalidDB
	}
	opt, ok := r.Options[optionID]
	if !ok {
		return false, nil
	}
	return opt.PollID == pollID, nil
}

// MockPollResponseRepository mock implementation
type MockPollResponseRepository struct {
	Responses  map[uint]*entity.PollResponse
	NextID     uint
	ShouldFail bool
}

// NewMockPollResponseRepository create new mock
func NewMockPollResponseRepository() *MockPollResponseRepository {
	return &MockPollResponseRepository{
		Responses: make(map[uint]*entity.PollResponse),
		NextID:    1,
	}
}

// HasVoted mock has voted check
func (r *MockPollResponseRepository) HasVoted(db *gorm.DB, pollID uint, participantID uint) (bool, error) {
	if r.ShouldFail {
		return false, gorm.ErrInvalidDB
	}
	for _, resp := range r.Responses {
		if resp.PollID == pollID && resp.ParticipantID == participantID {
			return true, nil
		}
	}
	return false, nil
}

// Create mock create
func (r *MockPollResponseRepository) Create(db *gorm.DB, response *entity.PollResponse) error {
	if r.ShouldFail {
		return gorm.ErrInvalidDB
	}
	response.ID = r.NextID
	r.NextID++
	r.Responses[response.ID] = response
	return nil
}
