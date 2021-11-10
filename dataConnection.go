package P2pRTC

import (
	"github.com/pion/webrtc/v3"
)

type DataConnection struct {
	dataChannel *webrtc.DataChannel
}
func (DataConnection *DataConnection) SendText(text string) error {
	return DataConnection.dataChannel.SendText(text)
}
func (DataConnection *DataConnection) Send(data []byte) error {
	return DataConnection.dataChannel.Send(data)
}
func (DataConnection *DataConnection) Close() error {
	return DataConnection.dataChannel.Close()
}
func (DataConnection *DataConnection) ID() uint16 {
	return *DataConnection.dataChannel.ID()
}
func (DataConnection *DataConnection) Label() string {
	return DataConnection.dataChannel.Label()
}
func (DataConnection *DataConnection) Ordered() bool {
	return DataConnection.dataChannel.Ordered()
}
func (DataConnection *DataConnection) Closing() bool {
	return DataConnection.dataChannel.ReadyState() == webrtc.DataChannelStateClosing
}

func (DataConnection *DataConnection) Connecting() bool {
	return DataConnection.dataChannel.ReadyState() == webrtc.DataChannelStateConnecting
}

func (DataConnection *DataConnection) Opened() bool {
	return DataConnection.dataChannel.ReadyState() == webrtc.DataChannelStateOpen
}

func (DataConnection *DataConnection) Closed() bool {
	return DataConnection.dataChannel.ReadyState() == webrtc.DataChannelStateClosed
}