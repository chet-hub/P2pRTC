package P2pRTC

import (
	"encoding/json"
	"github.com/pion/webrtc/v3"
	"log"
	"time"
)

type P2pSocket struct {
	StartTime          time.Time
	isServer           bool
	isTrickleICE       bool
	webRTConfiguration webrtc.Configuration
	dataChanelConfig   webrtc.DataChannelInit

	peerConnection *webrtc.PeerConnection

	onConnectionStateChange    func(webrtc.PeerConnectionState)
	onICEConnectionStateChange func(webrtc.ICEConnectionState)
	onICECandidate             func(*webrtc.ICECandidate)
	onICEGatheringStateChange  func(webrtc.ICEGathererState)
	onSignalingStateChange     func(webrtc.SignalingState)
	onNegotiationNeeded        func()
	onCreateDataChannel        func(*webrtc.DataChannel)

	dataChannel *webrtc.DataChannel

	OnSignal  func(string)
	OnOpen    func(*P2pSocket)
	OnMessage func(*P2pSocket, []byte)
	OnClose   func(*P2pSocket)
	OnError   func(*P2pSocket, []byte)
}

func newP2pSocket(isServer bool, isTrickleICE bool, webRtcConfig string, dataChannelConfig string) (*P2pSocket, error) {

	WebRTCConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				//http://olegh.ftp.sh/public-stun.txt
				URLs: []string{
					"stun:s1.taraba.net",
					"stun:s2.taraba.net",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:stun3.l.google.com:19302",
					"stun:stun4.l.google.com:19302",
				},
			},
		},
	}
	if webRtcConfig != "" {
		if e := json.Unmarshal([]byte(webRtcConfig), &WebRTCConfig); e != nil {
			return nil, e
		}
	}
	DataChanelConfig := webrtc.DataChannelInit{}
	if dataChannelConfig != "" {
		if e := json.Unmarshal([]byte(dataChannelConfig), &DataChanelConfig); e != nil {
			return nil, e
		}
	}
	configBytes, _ := json.Marshal(WebRTCConfig)
	log.Printf("webrtc.Configuration:  '%s'\n", string(configBytes))
	configBytes, _ = json.Marshal(WebRTCConfig)
	log.Printf("webrtc.DataChannelInit :  '%s'\n", string(configBytes))

	var Con = P2pSocket{}
	Con.isServer = isServer
	Con.isTrickleICE = isTrickleICE
	Con.webRTConfiguration = WebRTCConfig
	Con.dataChanelConfig = DataChanelConfig

	Con.OnSignal = func(message string) {
		log.Printf("Fire signalMessage:  '%s'\n", message)
	}

	Con.OnOpen = func(dataChannel *P2pSocket) {
		log.Printf("DataChannel[%d] OnOpen '%s'\n", Con.dataChannel.ID(), Con.dataChannel.Label())
	}
	Con.OnMessage = func(dataChannel *P2pSocket, msg []byte) {
		log.Printf("DataChannel[%d] Message '%s': '%s'\n", Con.dataChannel.ID(), Con.dataChannel.Label(), string(msg))
	}
	Con.OnClose = func(dataChannel *P2pSocket) {
		log.Printf("DataChannel[%d] OnClose '%s'\n", Con.dataChannel.ID(), Con.dataChannel.Label())
	}
	Con.OnError = func(dataChannel *P2pSocket, err []byte) {
		log.Printf("DataChannel[%d] OnError '%s' '%s' \n", Con.dataChannel.ID(), Con.dataChannel.Label(), string(err))
	}

	Con.onCreateDataChannel = func(d *webrtc.DataChannel) {
		//if Con.isServer {
		//	log.Printf("\tServer DataChannel New '%s' \n", d.Label())
		//} else {
		//	log.Printf("Client DataChannel New '%s' \n", d.Label())
		//}

		Con.dataChannel = d

		d.OnOpen(func() {
			Con.OnOpen(&Con)
		})
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			Con.OnMessage(&Con, msg.Data)
		})
		d.OnClose(func() {
			Con.OnClose(&Con)
			Con.dataChannel.Close()
			Con.peerConnection.Close()
		})
		d.OnError(func(err error) {
			log.Printf(err.Error())
			Con.OnError(&Con, []byte(err.Error()))
			Con.dataChannel.Close()
		})
	}

	//This happens whenever the aggregate state of the connection changes. The aggregate state is a combination of the
	//states of all of the individual network transports being used by the connection.
	Con.onConnectionStateChange = func(s webrtc.PeerConnectionState) {
		if Con.isServer {
			log.Printf("\tServer PeerConnectionState changed'%s' \n", s.String())
		} else {
			log.Printf("Client PeerConnectionState changed'%s' \n", s.String())
		}
		//todo
		//if s== webrtc.PeerConnectionStateClosed {
		//	Con.OnClose()
		//	Con.peerConnection.Close()
		//}
	}

	//This happens when the state of the connection's ICE agent, as represented by the iceConnectionState property, changes.
	Con.onICEConnectionStateChange = func(state webrtc.ICEConnectionState) {
		if Con.isServer {
			log.Printf("\tServer onICEConnectionStateChange '%s' \n", state.String())
		} else {
			log.Printf("Client onICEConnectionStateChange '%s' \n", state.String())
		}
	}

	//happens whenever the local ICE agent needs to deliver a message to the other peer through the signaling server.
	//This lets the ICE agent perform negotiation with the remote peer without the browser itself needing to know any specifics
	//about the technology being used for signaling; implement this method to use whatever messaging technology you
	//choose to send the ICE candidate to the remote peer.
	Con.onICECandidate = func(candidate *webrtc.ICECandidate) {
		if candidate != nil && Con.isTrickleICE {
			var m, e = SignalMessage{Type: SignalCandidate, Candidate: *candidate}.String()
			if e != nil {
				Con.OnError(&Con, []byte(e.Error()))
			} else {
				Con.OnSignal(m)
			}
		}
	}

	//This happens when the ICE gathering state—that is, whether or not the ICE agent is actively gathering candidates—changes.
	//You don't need to watch for this event unless you have specific reasons to want to closely monitor the state of ICE gathering.
	Con.onICEGatheringStateChange = func(state webrtc.ICEGathererState) {
		if Con.isServer {
			log.Printf("\tServer onICEGatheringStateChange '%s' \n", state.String())
		} else {
			log.Printf("Client onICEGatheringStateChange '%s' \n", state.String())
		}
	}

	//This event is fired when a change has occurred which requires session negotiation.
	//This negotiation should be carried out as the offerer, because some session changes cannot be negotiated as the answerer.
	//Most commonly, the negotiationneeded event is fired after a send track is added to the RTCPeerConnection.
	//If the session is modified in a manner that requires negotiation while a negotiation is already in progress,
	//no negotiationneeded event will fire until negotiation completes, and only then if negotiation is still needed.
	Con.onNegotiationNeeded = func() {
		if Con.isServer {
			log.Printf("\tServer onNegotiationNeeded '%d' \n", Con.dataChannel.ID())
		} else {
			log.Printf("Client onNegotiationNeeded '%d' \n", Con.dataChannel.ID())
		}
	}

	peerConnection, err := webrtc.NewPeerConnection(Con.webRTConfiguration)
	if err != nil {
		return nil, err
	}
	Con.peerConnection = peerConnection
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		Con.onConnectionStateChange(s)
	})
	peerConnection.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		Con.onICEConnectionStateChange(s)
	})
	peerConnection.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		Con.onICECandidate(iceCandidate)
	})
	peerConnection.OnICEGatheringStateChange(func(iceGathererState webrtc.ICEGathererState) {
		Con.onICEGatheringStateChange(iceGathererState)
	})
	peerConnection.OnSignalingStateChange(func(signalState webrtc.SignalingState) {
		Con.onSignalingStateChange(signalState)
	})
	peerConnection.OnNegotiationNeeded(func() {
		Con.onNegotiationNeeded()
	})

	//happen either because of a call to setLocalDescription() or to setRemoteDescription().
	Con.onSignalingStateChange = func(state webrtc.SignalingState) {

		if Con.isServer {
			log.Printf("\tServer onSignalingStateChange '%s' \n", state.String())
		} else {
			log.Printf("Client onSignalingStateChange '%s' \n", state.String())
		}

		if !Con.isTrickleICE {
			return
		}

		if Con.isServer {

			switch state {
			// SignalingStateStable indicates there is no offer/answer exchange in
			// progress. This is also the initial state, in which case the local and
			// remote descriptions are nil.
			case webrtc.SignalingStateStable:

			// SignalingStateHaveLocalOffer indicates that a local description, of
			// type "offer", has been successfully applied.
			case webrtc.SignalingStateHaveLocalOffer:

			// SignalingStateHaveRemoteOffer indicates that a remote description, of
			// type "offer", has been successfully applied.
			case webrtc.SignalingStateHaveRemoteOffer:

			// SignalingStateHaveLocalPranswer indicates that a remote description
			// of type "offer" has been successfully applied and a local description
			// of type "pranswer" has been successfully applied.
			case webrtc.SignalingStateHaveLocalPranswer:

			// SignalingStateHaveRemotePranswer indicates that a local description
			// of type "offer" has been successfully applied and a remote description
			// of type "pranswer" has been successfully applied.
			case webrtc.SignalingStateHaveRemotePranswer:

			// SignalingStateClosed indicates The peerConnection has been closed.
			case webrtc.SignalingStateClosed:

			}
		} else {
			switch state {
			// SignalingStateStable indicates there is no offer/answer exchange in
			// progress. This is also the initial state, in which case the local and
			// remote descriptions are nil.
			case webrtc.SignalingStateStable:

			// SignalingStateHaveLocalOffer indicates that a local description, of
			// type "offer", has been successfully applied.
			case webrtc.SignalingStateHaveLocalOffer:

			// SignalingStateHaveRemoteOffer indicates that a remote description, of
			// type "offer", has been successfully applied.
			case webrtc.SignalingStateHaveRemoteOffer:

			// SignalingStateHaveLocalPranswer indicates that a remote description
			// of type "offer" has been successfully applied and a local description
			// of type "pranswer" has been successfully applied.
			case webrtc.SignalingStateHaveLocalPranswer:

			// SignalingStateHaveRemotePranswer indicates that a local description
			// of type "offer" has been successfully applied and a remote description
			// of type "pranswer" has been successfully applied.
			case webrtc.SignalingStateHaveRemotePranswer:
				//client set remote answer

			// SignalingStateClosed indicates The peerConnection has been closed.
			case webrtc.SignalingStateClosed:

			}
		}

	}
	Con.StartTime = time.Now()
	return &Con, nil
}

func (Con *P2pSocket) _apply(isTrickleICE bool) (webrtc.SessionDescription, error) {
	empty := webrtc.SessionDescription{}
	if Con.isServer {
		return empty, nil
	}
	channel, err := Con.peerConnection.CreateDataChannel("Client", &Con.dataChanelConfig)
	if err != nil {
		return empty, err
	}
	Con.dataChannel = channel
	Con.onCreateDataChannel(channel)
	offer, e := Con.peerConnection.CreateOffer(nil)
	if e != nil {
		return empty, e
	}
	gatherComplete := webrtc.GatheringCompletePromise(Con.peerConnection)
	if e = Con.peerConnection.SetLocalDescription(offer); e != nil {
		return empty, e
	}
	var SignalType = SignalSdp
	if !isTrickleICE {
		<-gatherComplete
		SignalType = SignalApply
	}
	if Con.peerConnection.LocalDescription() != nil {
		var m, e = SignalMessage{Type: SignalType, Description: *Con.peerConnection.LocalDescription()}.String()
		if e != nil {
			Con.OnError(Con, []byte(e.Error()))
		} else {
			Con.OnSignal(m)
		}
	}

	//var applyString = utils.Encode(Con.peerConnection.LocalDescription())
	//log.Printf("applyString : '%s'\n", applyString)
	//return applyString, nil
	return *Con.peerConnection.LocalDescription(), nil
}

func (Con *P2pSocket) _connect(offer webrtc.SessionDescription) error {
	//offer := webrtc.SessionDescription{}
	//utils.Decode(approveString, &offer)
	if sdpErr := Con.peerConnection.SetRemoteDescription(offer); sdpErr != nil {
		return sdpErr
	}
	return nil
}

func (Con *P2pSocket) _approve(isTrickleICE bool, offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if !Con.isServer {
		return nil, nil
	}
	Con.peerConnection.OnDataChannel(func(channel *webrtc.DataChannel) {
		Con.dataChannel = channel
		Con.onCreateDataChannel(channel)
	})
	//offer := webrtc.SessionDescription{}
	//utils.Decode(applyString, &offer)
	err := Con.peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return nil, err
	}
	answer, err := Con.peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}
	gatherComplete := webrtc.GatheringCompletePromise(Con.peerConnection)
	err = Con.peerConnection.SetLocalDescription(answer)
	if err != nil {
		return nil, err
	}
	var SignalType = SignalSdp
	if !isTrickleICE {
		<-gatherComplete
		SignalType = SignalApprove
	}
	if Con.peerConnection.LocalDescription() != nil {
		var m, e = SignalMessage{Type: SignalType, Description: *Con.peerConnection.LocalDescription()}.String()
		if e != nil {
			Con.OnError(Con, []byte(e.Error()))
		} else {
			Con.OnSignal(m)
		}
	}
	//var approveString = utils.Encode(Con.peerConnection.LocalDescription())
	//return approveString, nil
	return Con.peerConnection.LocalDescription(), nil
}

func NewClient(isTrickleICE bool, webRtcConfig string, dataChannelConfig string) (*P2pSocket, error) {
	return newP2pSocket(false, isTrickleICE, webRtcConfig, dataChannelConfig)
}

func NewServer(isTrickleICE bool, webRtcConfig string, dataChannelConfig string) (*P2pSocket, error) {
	return newP2pSocket(true, isTrickleICE, webRtcConfig, dataChannelConfig)
}

func (Con *P2pSocket) Connect() {
	var _, err = Con._apply(Con.isTrickleICE)
	if err != nil {
		Con.OnError(Con, []byte(err.Error()))
	}
}

func (Con *P2pSocket) Signal(signalMsg string) {
	var message = SignalMessage{}
	if e := json.Unmarshal([]byte(signalMsg), &message); e != nil {
		Con.OnError(Con, []byte(e.Error()))
	}

	if Con.isTrickleICE {
		//Trickle ICE event, fast
		if message.Type == SignalSdp {
			if Con.isServer {
				if message.Description.Type == webrtc.SDPTypeOffer {
					var _, err = Con._approve(true, message.Description)
					if err != nil {
						Con.OnError(Con, []byte(err.Error()))
					}
				}
			} else {
				var err = Con._connect(message.Description)
				if err != nil {
					Con.OnError(Con, []byte(err.Error()))
				}
			}
		}
		if message.Type == SignalCandidate {
			var err = Con.peerConnection.AddICECandidate(message.Candidate.ToJSON())
			if err != nil {
				Con.OnError(Con, []byte(err.Error()))
			}
		}
	} else {
		//No Trickle ICE event, slow
		if message.Type == SignalApply { // Get apply as server, then approve the apply
			var _, err = Con._approve(false, message.Description)
			if err != nil {
				Con.OnError(Con, []byte(err.Error()))
			}
		}
		if message.Type == SignalApprove { //Get approve as client, then connect
			var err = Con._connect(message.Description)
			if err != nil {
				Con.OnError(Con, []byte(err.Error()))
			}
		}
	}
}


func (Con *P2pSocket) Send(data []byte) error {
	return Con.dataChannel.Send(data)
}
func (Con *P2pSocket) Close() error {
	return Con.dataChannel.Close()
}
func (Con *P2pSocket) ID() uint16 {
	return *Con.dataChannel.ID()
}
func (Con *P2pSocket) Label() string {
	return Con.dataChannel.Label()
}
func (Con *P2pSocket) Ordered() bool {
	return Con.dataChannel.Ordered()
}
func (Con *P2pSocket) Closing() bool {
	return Con.dataChannel.ReadyState() == webrtc.DataChannelStateClosing
}

func (Con *P2pSocket) Connecting() bool {
	return Con.dataChannel.ReadyState() == webrtc.DataChannelStateConnecting
}

func (Con *P2pSocket) Opened() bool {
	return Con.dataChannel.ReadyState() == webrtc.DataChannelStateOpen
}

func (Con *P2pSocket) Closed() bool {
	return Con.dataChannel.ReadyState() == webrtc.DataChannelStateClosed
}
