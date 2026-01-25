package sfu

import (
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

// TrackInfo stores track with metadata for cleanup
type TrackInfo struct {
	Track    *webrtc.TrackLocalStaticRTP
	SourceID string // participantID yang mengirim track
	Done     chan struct{}
}

// ConferenceState represents the state of the conference
type ConferenceState struct {
	IsActive    bool             // apakah conference sudah dimulai
	HostID      string           // participantID host/presenter
	Speakers    map[string]bool  // participantID yang dipromote jadi speaker
	RaisedHands map[string]int64 // participantID -> timestamp raise hand
}

type Room struct {
	id         uint
	peers      map[string]*Peer
	tracks     []*TrackInfo
	lock       sync.RWMutex
	log        *logrus.Logger
	Conference *ConferenceState
}

func NewRoom(id uint, log *logrus.Logger) *Room {
	return &Room{
		id:     id,
		peers:  make(map[string]*Peer),
		tracks: make([]*TrackInfo, 0),
		log:    log,
		Conference: &ConferenceState{
			IsActive:    false,
			HostID:      "",
			Speakers:    make(map[string]bool),
			RaisedHands: make(map[string]int64),
		},
	}
}

func (r *Room) Join(participantID string, signalFunc func(interface{})) (*Peer, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Clean up existing if any
	if existing, ok := r.peers[participantID]; ok {
		existing.Close()
		// Also cleanup tracks from this participant
		r.cleanupTracksForPeer(participantID)
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
	for _, trackInfo := range r.tracks {
		if err := peer.AddTrack(trackInfo.Track); err != nil {
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

	// Cleanup tracks from this participant
	r.cleanupTracksForPeer(participantID)

	// Remove from raised hands if exists
	delete(r.Conference.RaisedHands, participantID)
	delete(r.Conference.Speakers, participantID)
}

// cleanupTracksForPeer removes all tracks from a specific peer (must be called with lock held)
func (r *Room) cleanupTracksForPeer(participantID string) {
	newTracks := make([]*TrackInfo, 0)
	for _, trackInfo := range r.tracks {
		if trackInfo.SourceID == participantID {
			// Signal the goroutine to stop
			close(trackInfo.Done)
			r.log.WithField("participant_id", participantID).Info("Track removed on leave")
		} else {
			newTracks = append(newTracks, trackInfo)
		}
	}
	r.tracks = newTracks
}

func (r *Room) GetPeer(participantID string) *Peer {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.peers[participantID]
}

func (r *Room) BroadcastTrack(sourceID string, remoteTrack *webrtc.TrackRemote) {
	// Create local track with descriptive stream ID for screen share detection
	streamID := remoteTrack.StreamID()
	trackID := remoteTrack.ID()

	r.log.WithFields(map[string]interface{}{
		"source_id": sourceID,
		"track_id":  trackID,
		"stream_id": streamID,
		"kind":      remoteTrack.Kind().String(),
	}).Info("Broadcasting new track")

	localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, trackID, streamID)
	if err != nil {
		r.log.Error("failed to create local track: ", err)
		return
	}

	// Create track info with done channel for cleanup
	trackInfo := &TrackInfo{
		Track:    localTrack,
		SourceID: sourceID,
		Done:     make(chan struct{}),
	}

	r.lock.Lock()
	r.tracks = append(r.tracks, trackInfo)
	peers := make([]*Peer, 0, len(r.peers))
	for pid, p := range r.peers {
		if pid != sourceID {
			peers = append(peers, p)
		}
	}
	r.lock.Unlock()

	// Add to other peers and trigger renegotiation
	for _, p := range peers {
		if err := p.AddTrack(localTrack); err != nil {
			r.log.Error("failed to add track to peer: ", err)
			continue
		}
		// Trigger renegotiation so audience receives the new track
		// This sends an offer to the client with the new track
		if err := p.Negotiate(); err != nil {
			r.log.WithField("error", err).Warn("failed to trigger renegotiation for new track")
		}
	}

	// Forward data with cleanup support
	go func() {
		buf := make([]byte, 1500)
		for {
			select {
			case <-trackInfo.Done:
				r.log.WithField("source_id", sourceID).Info("Track forwarding stopped")
				return
			default:
				i, _, err := remoteTrack.Read(buf)
				if err != nil {
					return
				}
				if _, err := localTrack.Write(buf[:i]); err != nil {
					return
				}
			}
		}
	}()
}

// Conference methods

// StartConference starts the conference (only host can do this)
func (r *Room) StartConference(hostID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Conference.IsActive {
		return nil // Already active
	}

	r.Conference.IsActive = true
	r.Conference.HostID = hostID
	r.Conference.Speakers[hostID] = true // Host is always a speaker
	return nil
}

// StopConference stops the conference
func (r *Room) StopConference(participantID string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Only host can stop
	if r.Conference.HostID != participantID {
		return nil
	}

	r.Conference.IsActive = false
	r.Conference.Speakers = make(map[string]bool)
	r.Conference.RaisedHands = make(map[string]int64)
	return nil
}

// RaiseHand participant raises hand to request speaking
func (r *Room) RaiseHand(participantID string, timestamp int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Conference.RaisedHands[participantID] = timestamp
}

// LowerHand participant lowers their hand
func (r *Room) LowerHand(participantID string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.Conference.RaisedHands, participantID)
}

// PromoteSpeaker promotes a participant to speaker (only host can do this)
func (r *Room) PromoteSpeaker(hostID, participantID string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Conference.HostID != hostID {
		return false
	}

	r.Conference.Speakers[participantID] = true
	delete(r.Conference.RaisedHands, participantID)
	return true
}

// DemoteSpeaker demotes a speaker back to audience
func (r *Room) DemoteSpeaker(hostID, participantID string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.Conference.HostID != hostID {
		return false
	}

	delete(r.Conference.Speakers, participantID)
	return true
}

// IsSpeaker checks if participant is a speaker
func (r *Room) IsSpeaker(participantID string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.Conference.Speakers[participantID]
}

// IsHost checks if participant is the host
func (r *Room) IsHost(participantID string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.Conference.HostID == participantID
}

// GetConferenceState returns current conference state
func (r *Room) GetConferenceState() ConferenceState {
	r.lock.RLock()
	defer r.lock.RUnlock()

	// Deep copy maps
	speakers := make(map[string]bool)
	for k, v := range r.Conference.Speakers {
		speakers[k] = v
	}
	raisedHands := make(map[string]int64)
	for k, v := range r.Conference.RaisedHands {
		raisedHands[k] = v
	}

	return ConferenceState{
		IsActive:    r.Conference.IsActive,
		HostID:      r.Conference.HostID,
		Speakers:    speakers,
		RaisedHands: raisedHands,
	}
}
