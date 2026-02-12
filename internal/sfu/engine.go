package sfu

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

type SFUManager struct {
	rooms map[uint]*Room
	lock  sync.RWMutex
	log   *logrus.Logger
}

func NewSFUManager(log *logrus.Logger) *SFUManager {
	return &SFUManager{
		rooms: make(map[uint]*Room),
		log:   log,
	}
}

func (m *SFUManager) GetRoom(roomID uint) *Room {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.rooms[roomID]; !ok {
		m.rooms[roomID] = NewRoom(roomID, m.log)
	}
	return m.rooms[roomID]
}

func (m *SFUManager) CreatePeer(roomID uint, participantID string, signalFunc func(interface{})) (*Peer, error) {
	room := m.GetRoom(roomID)
	return room.Join(participantID, signalFunc)
}

func (m *SFUManager) HandleOffer(roomID uint, participantID string, offer webrtc.SessionDescription) error {
	room := m.GetRoom(roomID)
	peer := room.GetPeer(participantID)
	if peer == nil {
		return fmt.Errorf("peer not found")
	}
	return peer.HandleOffer(offer)
}

func (m *SFUManager) HandleAnswer(roomID uint, participantID string, answer webrtc.SessionDescription) error {
	room := m.GetRoom(roomID)
	peer := room.GetPeer(participantID)
	if peer == nil {
		return fmt.Errorf("peer not found")
	}
	return peer.HandleAnswer(answer)
}

func (m *SFUManager) HandleCandidate(roomID uint, participantID string, candidate webrtc.ICECandidateInit) error {
	room := m.GetRoom(roomID)
	peer := room.GetPeer(participantID)
	if peer == nil {
		return fmt.Errorf("peer not found")
	}
	return peer.HandleCandidate(candidate)
}

func (m *SFUManager) RemovePeer(roomID uint, participantID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if room, ok := m.rooms[roomID]; ok {
		room.Leave(participantID)
		if len(room.peers) == 0 {
			delete(m.rooms, roomID)
		}
	}
}
