package sfu

import (
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

type Peer struct {
	id                    string
	pc                    *webrtc.PeerConnection
	log                   *logrus.Logger
	signalFunc            func(interface{})
	OnTrack               func(*webrtc.TrackRemote, *webrtc.RTPReceiver)
	OnRenegotiationNeeded func()
	lock                  sync.Mutex
}

func NewPeer(id string, log *logrus.Logger, signalFunc func(interface{})) (*Peer, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// MediaEngine setup might be needed for v4 if codecs need registration
	// For now using default which typically includes VP8, Opus etc.

	// Create PeerConnection
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	p := &Peer{
		id:         id,
		pc:         pc,
		log:        log,
		signalFunc: signalFunc,
	}

	// Handlers
	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		// Convert to proper JSON struct
		candidateInit := c.ToJSON()

		payload := map[string]interface{}{
			"type":             "candidate",
			"candidate":        candidateInit.Candidate,
			"sdpMid":           candidateInit.SDPMid,
			"sdpMLineIndex":    candidateInit.SDPMLineIndex,
			"usernameFragment": candidateInit.UsernameFragment,
		}
		p.signalFunc(payload)
	})

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if p.OnTrack != nil {
			p.OnTrack(track, receiver)
		}
	})

	// Handle renegotiation
	pc.OnNegotiationNeeded(func() {
		if p.OnRenegotiationNeeded != nil {
			p.OnRenegotiationNeeded()
		}
	})

	return p, nil
}

func (p *Peer) HandleOffer(offer webrtc.SessionDescription) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if err := p.pc.SetRemoteDescription(offer); err != nil {
		return err
	}

	answer, err := p.pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	if err := p.pc.SetLocalDescription(answer); err != nil {
		return err
	}

	// Send answer
	payload := map[string]interface{}{
		"type": answer.Type.String(),
		"sdp":  answer.SDP,
	}
	p.signalFunc(payload)

	return nil
}

func (p *Peer) HandleAnswer(answer webrtc.SessionDescription) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.pc.SetRemoteDescription(answer)
}

func (p *Peer) HandleCandidate(candidate webrtc.ICECandidateInit) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.pc.AddICECandidate(candidate)
}

func (p *Peer) AddTrack(track webrtc.TrackLocal) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, err := p.pc.AddTrack(track); err != nil {
		return err
	}

	// We might need to renegotiate now
	// But defer it to the caller or OnNegotiationNeeded
	return nil
}

// Negotiate creates an offer and sends it to the other peer (client)
func (p *Peer) Negotiate() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	offer, err := p.pc.CreateOffer(nil)
	if err != nil {
		return err
	}

	if err := p.pc.SetLocalDescription(offer); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type": offer.Type.String(),
		"sdp":  offer.SDP,
	}
	p.signalFunc(payload)
	return nil
}

func (p *Peer) Close() {
	if p.pc != nil {
		p.pc.Close()
	}
}

// Helper structs for Signals if needed
type WebsocketMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
