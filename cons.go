package P2pRTC

import (
	"sync"
	"time"
)

var (
	listLock   sync.RWMutex
	connections []*p2pSocket

	OnSignal = func(message string){}
	OnOpen  = func(dataChannel *DataConnection){}
	OnMessage = func (dataChannel *DataConnection, msg []byte){}
	OnClose = func (dataChannel *DataConnection){}
	OnError = func (dataChannel *DataConnection, err []byte){}
)

func add(socket *p2pSocket){
	defer func() {
		listLock.Lock()
	}()
	listLock.Unlock()
	connections = append(connections,socket)
}

func remove(socket *p2pSocket){
	defer func() {
		listLock.Lock()
	}()
	listLock.Unlock()
	for i, v := range connections {
		if v == socket{
			connections[i] = connections[len(connections)-1]
			connections = connections[:len(connections)-1]
			return
		}
	}
}

func run(connectTime int,isServer bool, isTrickleICE bool,webRtcConfig string, dataChannelConfig string) error {
	var connect, e = newP2pSocket(isServer,isTrickleICE,webRtcConfig,dataChannelConfig)
	if e != nil {
		return e
	}
	var timer = time.NewTimer(time.Duration(connectTime) * time.Second)

	connect.OnSignal = func(message string) {
		OnSignal(message)
	}
	connect.OnOpen = func(dataChannel *DataConnection) {
		OnOpen(dataChannel)
		timer.Stop()
	}
	connect.OnMessage = func(dataChannel *DataConnection,msg []byte) {
		OnMessage(dataChannel,msg)
	}
	connect.OnClose = func(dataChannel *DataConnection) {
		OnClose(connect.dataConnection)
		remove(connect)
		timer.Stop()
	}
	connect.OnError = func(dataChannel *DataConnection,err []byte) {
		OnError(connect.dataConnection,err)
	}

	connect.Connect()
	add(connect)

	<-timer.C
	if connect.dataConnection == nil || !connect.dataConnection.Opened(){
		remove(connect)
	}
	return nil
}

func Connect(connectionTime int, isTrickleICE bool,webRtcConfig string, dataChannelConfig string) error {
	return run(connectionTime,false,isTrickleICE,webRtcConfig,dataChannelConfig)
}

func Accept(connectionTime int, isTrickleICE bool,webRtcConfig string, dataChannelConfig string) error {
	return run(connectionTime,true,isTrickleICE,webRtcConfig,dataChannelConfig)
}






