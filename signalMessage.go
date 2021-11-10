package P2pRTC

import (
	"encoding/json"
	"errors"
	"github.com/pion/webrtc/v3"
)

type SignalType int

const (
	SignalSdp       SignalType = iota + 1 // c0 == 1
	SignalCandidate                       // c1 == 2
	SignalApply                           // c2 == 3
	SignalApprove                         // c2 == 4
)

type SignalMessage struct {
	Type        SignalType
	Description webrtc.SessionDescription
	Candidate   webrtc.ICECandidate
}

func (m SignalMessage) String() (string,error) {
	 bytes,e := json.Marshal(m)
	 return string(bytes),e
}


func (t SignalType) String() string {
	switch t {
	case SignalSdp:
		return "SignalSdp"
	case SignalCandidate:
		return "SignalCandidate"
	case SignalApply:
		return "SignalApply"
	case SignalApprove:
		return "SignalApprove"
	default:
		return errors.New("unknown").Error()
	}
}
