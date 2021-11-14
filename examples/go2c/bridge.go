package main

/*
#include <stddef.h>
#include <stdint.h>

typedef uint64_t GoP2pSocket;
typedef void (*GoOnSignalCb)(GoP2pSocket p2pSocket,char *message);
typedef void (*GoOnOpenCb)(GoP2pSocket p2pSocket);
typedef void (*GoOnCloseCb)(GoP2pSocket p2pSocket);
typedef void (*GoOnMessageCb)(GoP2pSocket p2pSocket, char *message, size_t lenght);
typedef void (*GoOnErrorCb)(GoP2pSocket p2pSocket, char *message);

// Calling C function pointers is currently not supported, however you can
// declare Go variables which hold C function pointers and pass them back
// and forth between Go and C. C code may call function pointers received from Go.
// Reference: https://golang.org/cmd/cgo/#hdr-Go_references_to_C
inline void bridge_on_signal(GoOnSignalCb cb, GoP2pSocket p2pSocket,char *message)
{
    cb(p2pSocket,message);
}
inline void bridge_on_open(GoOnOpenCb cb, GoP2pSocket p2pSocket)
{
    cb(p2pSocket);
}
inline void bridge_on_close(GoOnCloseCb cb, GoP2pSocket p2pSocket)
{
    cb(p2pSocket);
}
inline void bridge_on_message(GoOnMessageCb cb, GoP2pSocket p2pSocket, char *message, size_t length)
{
    cb(p2pSocket,message,length);
}
inline void bridge_on_error(GoOnErrorCb cb, GoP2pSocket p2pSocket,char *message)
{
    cb(p2pSocket,message);
}
*/
import "C"

import (
	"github.com/chet-hub/P2pRTC"
	"sync/atomic"
	"unsafe"
)

var store = map[uint64]*P2pRTC.P2pSocket{}
var Id uint64 = 0

func getId() uint64 {
	atomic.AddUint64(&Id, 1)
	return Id
}

//export newClient
func newClient(cIsTrickleIce C.int, cWebrtcconfig *C.char, cDatachannelconfig *C.char) C.GoP2pSocket {
	var isTrickleIce = false
	if cIsTrickleIce != 0 {
		isTrickleIce = true
	}
	webRtcConfig := C.GoString(cWebrtcconfig)
	dataChannelConfig := C.GoString(cDatachannelconfig)
	connection, err := P2pRTC.NewClient(isTrickleIce, webRtcConfig, dataChannelConfig)
	if err != nil {
		return 0
	} else {
		id := getId()
		store[id] = connection
		return C.ulonglong(id)
	}
}

//export newServer
func newServer(cIsTrickleIce C.int, cWebrtcconfig *C.char, cDatachannelconfig *C.char) C.GoP2pSocket {
	var isTrickleIce = false
	if cIsTrickleIce != 0 {
		isTrickleIce = true
	}
	webRtcConfig := C.GoString(cWebrtcconfig)
	dataChannelConfig := C.GoString(cDatachannelconfig)
	connection, err := P2pRTC.NewServer(isTrickleIce, webRtcConfig, dataChannelConfig)
	if err != nil {
		return 0
	} else {
		id := getId()
		store[id] = connection
		return C.ulonglong(id)
	}
}
//export connect
func connect(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.Connect()
		return 1
	} else {
		return 0
	}
}

//export signal
func signal(p2pSocket C.GoP2pSocket, message *C.char) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.Signal(C.GoString(message))
		return 1
	} else {
		return 0
	}
}

//export closed
func closed(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		if connection.Closed() == true {
			return 1
		} else {
			return -1
		}
	} else {
		return 0
	}
}

//export opened
func opened(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		if connection.Opened() == true {
			return 1
		} else {
			return -1
		}
	} else {
		return 0
	}
}

//export connecting
func connecting(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		if connection.Connecting() == true {
			return 1
		} else {
			return -1
		}
	} else {
		return 0
	}
}

//export ordered
func ordered(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		if connection.Ordered() == true {
			return 1
		} else {
			return -1
		}
	} else {
		return 0
	}
}


//export close
func close(p2pSocket C.GoP2pSocket) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		if connection.Close() == nil {
			return 1
		} else {
			return -1
		}
	} else {
		return 0
	}
}

//export send
func send(p2pSocket C.GoP2pSocket, p unsafe.Pointer, length C.int) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.Send(C.GoBytes(p, length))
		return 1
	} else {
		return 0
	}
}


//export listenOnSignal
func listenOnSignal(p2pSocket C.GoP2pSocket, cb C.GoOnSignalCb) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.OnSignal = func(s string) {
			C.bridge_on_signal(cb, p2pSocket, C.CString(s))
		}
		return 1
	} else {
		return 0
	}
}

//export listenOnError
func listenOnError(p2pSocket C.GoP2pSocket, cb C.GoOnSignalCb) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.OnError = func(dataChannel *P2pRTC.P2pSocket, err []byte) {
			C.bridge_on_error(cb, p2pSocket, C.CString(string(err)))
		}
		return 1
	} else {
		return 0
	}
}

//export listenOnMessage
func listenOnMessage(p2pSocket C.GoP2pSocket, cb C.GoOnMessageCb) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.OnMessage = func(dataChannel *P2pRTC.P2pSocket, message []byte) {
			C.bridge_on_message(cb, p2pSocket, (*C.char)(C.CBytes(message)), C.size_t(len(message)))
		}
		return 1
	} else {
		return 0
	}
}

//export listenOnOpen
func listenOnOpen(p2pSocket C.GoP2pSocket, cb C.GoOnOpenCb) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.OnOpen = func(dataChannel *P2pRTC.P2pSocket) {
			C.bridge_on_open(cb, p2pSocket)
		}
		return 1
	} else {
		return 0
	}
}

//export listenOnClose
func listenOnClose(p2pSocket C.GoP2pSocket, cb C.GoOnOpenCb) C.int {
	if connection, has := store[uint64(p2pSocket)]; has {
		connection.OnClose = func(dataChannel *P2pRTC.P2pSocket) {
			C.bridge_on_close(cb, p2pSocket)
		}
		return 1
	} else {
		return 0
	}
}


func main() {}