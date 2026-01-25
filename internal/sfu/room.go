package sfu

import (
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

type Room struct {
	id     uint
	peers  map[string]*Peer
	tracks []*webrtc.TrackLocalStaticRTP
	lock   sync.RWMutex
	log    *logrus.Logger
}

func NewRoom(id uint, log *logrus.Logger) *Room {
	return &Room{
		id:     id,
		peers:  make(map[string]*Peer),
		tracks: make([]*webrtc.TrackLocalStaticRTP, 0),
		log:    log,
	}
}

func (r *Room) Join(participantID string, signalFunc func(interface{})) (*Peer, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Clean up existing if any
	if existing, ok := r.peers[participantID]; ok {
		existing.Close()
	}

	peer, err := NewPeer(participantID, r.log, signalFunc)
	if err != nil {
		return nil, err
	}

	// Handle incoming tracks from this peer
	peer.OnTrack = func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		r.log.WithField("participant_id", participantID).Info("Track received")
		r.BroadcastTrack(participantID, remoteTrack)
	}

	peer.OnRenegotiationNeeded = func() {
		// Client implementation often handles this, or we might need to send Offer
		// For simplicity, we assume client can handle offers from us
		// But usually we just call Negotiate()
		if err := peer.Negotiate(); err != nil {
			r.log.Error("Negotiation failed", err)
		}
	}

	r.peers[participantID] = peer

	// Subscribe new peer to all existing tracks
	for _, track := range r.tracks {
		if err := peer.AddTrack(track); err != nil {
			r.log.Error("Failed to add existing track to new peer", err)
		}
	}

	// Trigger negotiation for the new tracks
	// (We do this after loop to batch it, but AddTrack might need immediate)
	// Actually NewPeer initializes connection. Real negotiation happens when we set remote description or AddTrack.

	return peer, nil
}

func (r *Room) Leave(participantID string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if peer, ok := r.peers[participantID]; ok {
		peer.Close()
		delete(r.peers, participantID)
	}
}

func (r *Room) GetPeer(participantID string) *Peer {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.peers[participantID]
}

func (r *Room) BroadcastTrack(sourceID string, remoteTrack *webrtc.TrackRemote) {
	// Create local track
	localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), remoteTrack.StreamID())
	if err != nil {
		r.log.Error("failed to create local track: ", err)
		return
	}

	r.lock.Lock()
	r.tracks = append(r.tracks, localTrack)
	peers := make([]*Peer, 0, len(r.peers))
	for pid, p := range r.peers {
		if pid != sourceID {
			peers = append(peers, p)
		}
	}
	r.lock.Unlock()

	// Add to other peers
	for _, p := range peers {
		if err := p.AddTrack(localTrack); err != nil {
			r.log.Error("failed to add track to peer: ", err)
		}
	}

	// Forward data
	go func() {
		buf := make([]byte, 1500)
		for {
			i, _, err := remoteTrack.Read(buf)
			if err != nil {
				return
			}
			if _, err := localTrack.Write(buf[:i]); err != nil {
				return
			}
		}
	}()
}
