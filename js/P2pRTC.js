class DataConnection {
    dataChannel

    constructor(d) {
        this.dataChannel = d
    }

    Send(data) {
        return this.dataChannel.send(data)
    }

    Close() {
        return this.dataChannel.close()
    }

    ID() {
        return this.dataChannel.id
    }

    Label() {
        return this.dataChannel.label
    }

    Ordered() {
        return this.dataChannel.ordered
    }

    Connecting() {
        return this.dataChannel.readyState == "connecting"
    }

    Opened() {
        return this.dataChannel.readyState == "open"
    }

    Closed() {
        return this.dataChannel.readyState == "closed"
    }
}

class SignalType {
    static SignalSdp = 1
    static SignalCandidate = 2
    static SignalApply = 3
    static SignalApprove = 4
}

class SignalMessage {
    Type
    Description
    Candidate

    constructor(type, value) {
        if (SignalType.SignalCandidate === type) {
            this.Type = type
            this.Candidate = value
        } else if (SignalType.SignalSdp === type) {
            this.Type = type
            this.Description = value
        } else if (SignalType.SignalApply === type) {
            this.Type = type
            this.Description = value
        } else if (SignalType.SignalApprove === type) {
            this.Type = type
            this.Description = value
        } else {
            console.error("unknown SignalMessage")
        }
    }

    toSignalString() {
        return JSON.stringify(this)
    }

    toString() {
        switch (this.Type) {
            case SignalType.SignalSdp:
                return "SignalSdp"
            case SignalType.SignalCandidate:
                return "SignalCandidate"
            case SignalType.SignalApply:
                return "SignalApply"
            case SignalType.SignalApprove:
                return "SignalApprove"
            default:
                console.error("unknown SignalMessage")
                return null
        }
    }
}

class P2pRTC {

    StartTime
    isServer
    isTrickleICE
    webRTConfiguration
    dataChanelConfig
    peerConnection
    onConnectionStateChange
    onICEConnectionStateChange
    onICECandidate
    onICEGatheringStateChange
    onSignalingStateChange
    onNegotiationNeeded
    onCreateDataChannel
    dataChannel
    dataConnection

    constructor(isServer, isTrickleICE, webRtcConfig, dataChannelConfig) {
        const Con = this
        Con.isServer = isServer
        Con.isTrickleICE = isTrickleICE
        if (!webRtcConfig) {
            webRtcConfig = {"iceServers": [{"urls": ["stun:s1.taraba.net", "stun:s2.taraba.net", "stun:stun1.l.google.com:19302", "stun:stun2.l.google.com:19302", "stun:stun3.l.google.com:19302", "stun:stun4.l.google.com:19302"]}]}
        }
        Con.webRTConfiguration = webRtcConfig
        Con.dataChanelConfig = dataChannelConfig
        Con.onCreateDataChannel = (d) => {
            if (Con.isServer) {
                console.debug("\tServer DataChannel created")
            } else {
                console.debug("Client DataChannel created")
            }

            Con.dataChannel = d
            Con.dataConnection = new DataConnection(d)

            d.onopen = function () {
                Con.OnOpen(Con.dataConnection)
            }
            d.onmessage = function (msg) {
                Con.OnMessage(Con.dataConnection, msg.data)
            }
            d.onclose = function () {
                Con.OnClose(Con.dataConnection)
                Con.dataChannel.Close()
                Con.peerConnection.Close()
            }
            d.onerror = function (err) {
                Con.OnError(Con.dataConnection,err)
            }
        }

        //This happens whenever the aggregate state of the connection changes. The aggregate state is a combination of the
        //states of all of the individual network transports being used by the connection.
        Con.onConnectionStateChange = function (peerConnectionState) {
            //"closed" | "connected" | "connecting" | "disconnected" | "failed" | "new"
            const state = peerConnectionState.currentTarget.connectionState
            if (Con.isServer) {
                console.debug("\tServer PeerConnectionState changed '%s' \n", state)
            } else {
                console.debug("Client PeerConnectionState changed '%s' \n", state)
            }
        }

        //This happens when the state of the connection's ICE agent, as represented by the iceConnectionState property, changes.
        Con.onICEConnectionStateChange = function (iCEConnectionState) {
            const state = iCEConnectionState.currentTarget.iceConnectionState
            if (Con.isServer) {
                console.debug("\tServer onICEConnectionStateChange '%s' \n", state)
            } else {
                console.debug("Client onICEConnectionStateChange '%s' \n", state)
            }
        }

        //happens whenever the local ICE agent needs to deliver a message to the other peer through the signaling server.
        //This lets the ICE agent perform negotiation with the remote peer without the browser itself needing to know any specifics
        //about the technology being used for signaling; implement this method to use whatever messaging technology you
        //choose to send the ICE candidate to the remote peer.
        Con.onICECandidate = function (event) {
            if (event && event.candidate && Con.isTrickleICE) {
                Con.OnSignal(new SignalMessage(SignalType.SignalCandidate, event.candidate).toSignalString())
            }
        }

        //This happens when the ICE gathering state—that is, whether or not the ICE agent is actively gathering candidates—changes.
        //You don't need to watch for this event unless you have specific reasons to want to closely monitor the state of ICE gathering.
        Con.onICEGatheringStateChange = function (state) {
            if (Con.isServer) {
                console.debug("\tServer onICEGatheringStateChange '%s' \n", state)
            } else {
                console.debug("Client onICEGatheringStateChange '%s' \n", state)
            }
        }

        //This event is fired when a change has occurred which requires session negotiation.
        //This negotiation should be carried out as the offerer, because some session changes cannot be negotiated as the answerer.
        //Most commonly, the negotiationneeded event is fired after a send track is added to the RTCPeerConnection.
        //If the session is modified in a manner that requires negotiation while a negotiation is already in progress,
        //no negotiationneeded event will fire until negotiation completes, and only then if negotiation is still needed.
        Con.onNegotiationNeeded = function () {
            if (Con.isServer) {
                console.debug("\tServer onNegotiationNeeded '%d' \n", Con.dataChannel.id)
            } else {
                console.debug("Client onNegotiationNeeded '%d' \n", Con.dataChannel.id)
            }
        }

        Con.peerConnection = new RTCPeerConnection(Con.webRTConfiguration);

        Con.peerConnection.onconnectionstatechange = function (s) {
            Con.onConnectionStateChange(s)
        }
        Con.peerConnection.oniceconnectionstatechange = function (s) {
            Con.onICEConnectionStateChange(s)
        }
        Con.peerConnection.onicecandidate = function (iceCandidate) {
            Con.onICECandidate(iceCandidate)
        }
        Con.peerConnection.onicegatheringstatechange = function (event) {
            const state = event.currentTarget.iceGatheringState //iceConnectionState, connectionState,signalingState
            Con.onICEGatheringStateChange(state)
            if (state === "complete" && !Con.isTrickleICE) {
                if (Con.isServer) {
                    Con.OnSignal(new SignalMessage(SignalType.SignalApprove, Con.peerConnection.localDescription).toSignalString())
                } else {
                    Con.OnSignal(new SignalMessage(SignalType.SignalApply, Con.peerConnection.localDescription).toSignalString())
                }
            }
        }
        Con.peerConnection.onsignalingstatechange = function (signalState) {
            Con.onSignalingStateChange(signalState)
        }
        Con.peerConnection.onnegotiationneeded = function () {
            Con.onNegotiationNeeded()
        }

        //happen either because of a call to setLocalDescription() or to setRemoteDescription().
        Con.onSignalingStateChange = function (event) {
            const state = event.currentTarget.signalingState
            if (Con.isServer) {
                console.debug("\tServer onSignalingStateChange '%s' \n", state)
            } else {
                console.debug("Client onSignalingStateChange '%s' \n", state)
            }

            if (!Con.isTrickleICE) {
                return
            }

            if (Con.isServer) {

                //"closed" | "have-local-offer" | "have-local-pranswer" | "have-remote-offer" | "have-remote-pranswer" | "stable";
                switch (state) {
                    // SignalingStateStable indicates there is no offer/answer exchange in
                    // progress. This is also the initial state, in which case the local and
                    // remote descriptions are nil.
                    case "stable":
                        // if (Con.peerConnection.localDescription != null) {
                        //     Con.OnSignal(new SignalMessage(SignalType.SignalSdp, Con.peerConnection.localDescription))
                        // }
                        return
                    // SignalingStateHaveLocalOffer indicates that a local description, of
                    // type "offer", has been successfully applied.
                    case "have-local-offer" :
                        // if (Con.peerConnection.localDescription != null) {
                        //     Con.OnSignal(new SignalMessage(SignalType.SignalSdp, Con.peerConnection.localDescription))
                        // }
                        return
                    // SignalingStateHaveRemoteOffer indicates that a remote description, of
                    // type "offer", has been successfully applied.
                    case "have-remote-offer":
                        return

                    // SignalingStateHaveLocalPranswer indicates that a remote description
                    // of type "offer" has been successfully applied and a local description
                    // of type "pranswer" has been successfully applied.
                    case "have-local-pranswer":
                        return
                    // SignalingStateHaveRemotePranswer indicates that a local description
                    // of type "offer" has been successfully applied and a remote description
                    // of type "pranswer" has been successfully applied.
                    case "have-remote-pranswer":
                        return
                    // SignalingStateClosed indicates The peerConnection has been closed.
                    case "closed":
                        return

                }
            } else {
                switch (state) {

                    // SignalingStateStable indicates there is no offer/answer exchange in
                    // progress. This is also the initial state, in which case the local and
                    // remote descriptions are nil.
                    case "stable":
                        return
                    // SignalingStateHaveLocalOffer indicates that a local description, of
                    // type "offer", has been successfully applied.
                    case "have-local-offer" :
                        //client create offer and set local offer
                        // if (Con.peerConnection.localDescription != null) {
                        //     Con.OnSignal(new SignalMessage(SignalType.SignalSdp, Con.peerConnection.localDescription))
                        // }
                        return
                    // SignalingStateHaveRemoteOffer indicates that a remote description, of
                    // type "offer", has been successfully applied.
                    case "have-remote-offer":
                        return
                    // SignalingStateHaveLocalPranswer indicates that a remote description
                    // of type "offer" has been successfully applied and a local description
                    // of type "pranswer" has been successfully applied.
                    case "have-local-pranswer":
                        return
                    // SignalingStateHaveRemotePranswer indicates that a local description
                    // of type "offer" has been successfully applied and a remote description
                    // of type "pranswer" has been successfully applied.
                    case "have-remote-pranswer":
                        return
                    // SignalingStateClosed indicates The peerConnection has been closed.
                    case "closed":
                        return
                }
            }

        }

        Con.StartTime = new Date().getTime()
        if (!Con.isServer) {
            Con.apply(Con.isTrickleICE)
        }
    }

    static NewClient(isTrickleICE, webRtcConfig, dataChannelConfig) {
        return new P2pRTC(false, isTrickleICE, webRtcConfig, dataChannelConfig)
    }

    static NewServer(isTrickleICE, webRtcConfig, dataChannelConfig) {
        return new P2pRTC(true, isTrickleICE | false, webRtcConfig, dataChannelConfig)
    }

    Signal(message) {
        message = JSON.parse(message)
        const Con = this
        if (Con.isServer) {
            console.debug("\tServer signal got '%s'\n", message.toString())
        } else {
            console.debug("Client signal got '%s'\n", message.toString())
        }

        //Trickle ICE event, fast
        if (Con.isTrickleICE) {
            if (message.Type === SignalType.SignalSdp) {
                if (Con.isServer) {
                    Con.approve(true, message.Description)
                } else {
                    Con.connect(message.Description)
                }
            }
            if (message.Type === SignalType.SignalCandidate) {
                Con.peerConnection.addIceCandidate(message.Candidate).catch(Con.OnError)
            }
        } else {//No Trickle ICE event, slow
            if (message.Type === SignalType.SignalApply) { // Get apply as server, then approve the apply
                Con.approve(false, message.Description)
            }

            if (message.Type === SignalType.SignalApprove) { //Get approve as client, then connect
                Con.connect(message.Description)
            }
        }
    }

    OnSignal = function (message) {
        console.debug("Fire signalMessage:  '%s'\n", message)
    }

    OnOpen = function (dataChannel) {
        console.debug("DataChannel[%d] OnOpen '%s'\n", dataChannel.id, dataChannel.label)
    }

    OnMessage = function (dataChannel, msg) {
        console.debug("DataChannel[%d] Message '%s': '%s'\n", Con.dataChannel.id, Con.dataChannel.label, msg)
    }

    OnClose = function (dataChannel) {
        console.debug("DataChannel[%d] OnClose '%s'\n", Con.dataChannel.id, Con.dataChannel.label)
    }

    OnError = function (dataChannel,err) {
        console.debug("DataChannel[%d] OnError '%s' '%s' \n", Con.dataChannel.id, Con.dataChannel.label, err)
    }

    apply(isTrickleICE) {
        const Con = this
        Con.dataChannel = Con.peerConnection.createDataChannel("Client", Con.dataChanelConfig)
        Con.onCreateDataChannel(Con.dataChannel)
        if (isTrickleICE) {
            return Con.peerConnection.createOffer().then((localOffer) => {
                return Con.peerConnection.setLocalDescription(localOffer);
            }).then(() => {
                if (Con.peerConnection.localDescription != null) {
                    Con.OnSignal(new SignalMessage(SignalType.SignalSdp, Con.peerConnection.localDescription).toSignalString())
                }
            }).catch(Con.OnError);
        } else {
            return Con.peerConnection.createOffer().then((offer) => {
                return Con.peerConnection.setLocalDescription(offer)
            }).then(() => {
                return Promise.resolve(Con.peerConnection.localDescription)
            }).catch(Con.OnError)
        }
    }

    connect(offer) {
        return this.peerConnection.setRemoteDescription(offer)
    }

    approve(isTrickleICE, offer) {
        const Con = this
        Con.peerConnection.ondatachannel = function (event) {
            Con.dataChannel = event.channel
            Con.onCreateDataChannel(Con.dataChannel)
        }
        if (isTrickleICE) {
            Con.peerConnection.setRemoteDescription(offer).catch(Con.OnError)
            return Con.peerConnection.createAnswer().then((remoteAnswer) => {
                return Con.peerConnection.setLocalDescription(remoteAnswer);
            }).then(() => {
                if (Con.peerConnection.localDescription != null) {
                    Con.OnSignal(new SignalMessage(SignalType.SignalSdp, Con.peerConnection.localDescription).toSignalString())
                }
            }).catch(Con.OnError);
        } else {
            return Con.peerConnection.setRemoteDescription(offer).then(() => {
                return Con.peerConnection.createAnswer()
            }).then((answer) => {
                return Con.peerConnection.setLocalDescription(answer)
            }).then(() => {
                return Promise.resolve(Con.peerConnection.localDescription)
            }).catch(Con.OnError)
        }
    }

}












