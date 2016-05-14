package handshake

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/utils"
)

// ConnectionParametersManager stores the connection parameters
// Warning: Writes may only be done from the crypto stream, see the comment
// in GetSHLOMap().
type ConnectionParametersManager struct {
	params map[Tag][]byte
	mutex  sync.RWMutex

	idleConnectionStateLifetime        time.Duration
	sendStreamFlowControlWindow        protocol.ByteCount
	sendConnectionFlowControlWindow    protocol.ByteCount
	receiveStreamFlowControlWindow     protocol.ByteCount
	receiveConnectionFlowControlWindow protocol.ByteCount
}

// ErrTagNotInConnectionParameterMap is returned when a tag is not present in the connection parameters
var ErrTagNotInConnectionParameterMap = errors.New("Tag not found in ConnectionsParameter map")

// NewConnectionParamatersManager creates a new connection parameters manager
func NewConnectionParamatersManager() *ConnectionParametersManager {
	return &ConnectionParametersManager{
		params: map[Tag][]byte{
			TagMSPC: {0x64, 0x00, 0x00, 0x00}, // Max streams per connection = 100
		},
		idleConnectionStateLifetime:        protocol.InitialIdleConnectionStateLifetime,
		sendStreamFlowControlWindow:        protocol.InitialStreamFlowControlWindow,     // can only be changed by the client
		sendConnectionFlowControlWindow:    protocol.InitialConnectionFlowControlWindow, // can only be changed by the client
		receiveStreamFlowControlWindow:     protocol.ReceiveStreamFlowControlWindow,
		receiveConnectionFlowControlWindow: protocol.ReceiveConnectionFlowControlWindow,
	}
}

// SetFromMap reads all params
func (h *ConnectionParametersManager) SetFromMap(params map[Tag][]byte) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for key, value := range params {
		switch key {
		case TagMSPC, TagTCID:
			h.params[key] = value
		case TagICSL:
			clientValue, err := utils.ReadUint32(bytes.NewBuffer(value))
			if err != nil {
				return err
			}
			h.idleConnectionStateLifetime = h.negotiateIdleConnectionStateLifetime(time.Duration(clientValue) * time.Second)
		case TagSFCW:
			sendStreamFlowControlWindow, err := utils.ReadUint32(bytes.NewBuffer(value))
			if err != nil {
				return err
			}
			h.sendStreamFlowControlWindow = protocol.ByteCount(sendStreamFlowControlWindow)
		case TagCFCW:
			sendConnectionFlowControlWindow, err := utils.ReadUint32(bytes.NewBuffer(value))
			if err != nil {
				return err
			}
			h.sendConnectionFlowControlWindow = protocol.ByteCount(sendConnectionFlowControlWindow)
		}
	}

	return nil
}

func (h *ConnectionParametersManager) negotiateIdleConnectionStateLifetime(clientValue time.Duration) time.Duration {
	// TODO: what happens if the clients sets 0 seconds?
	return utils.MinDuration(clientValue, protocol.MaxIdleConnectionStateLifetime)
}

// getRawValue gets the byte-slice for a tag
func (h *ConnectionParametersManager) getRawValue(tag Tag) ([]byte, error) {
	h.mutex.RLock()
	rawValue, ok := h.params[tag]
	h.mutex.RUnlock()

	if !ok {
		return nil, ErrTagNotInConnectionParameterMap
	}
	return rawValue, nil
}

// GetSHLOMap gets all values (except crypto values) needed for the SHLO
func (h *ConnectionParametersManager) GetSHLOMap() map[Tag][]byte {
	sfcw := bytes.NewBuffer([]byte{})
	utils.WriteUint32(sfcw, uint32(h.GetReceiveStreamFlowControlWindow()))
	cfcw := bytes.NewBuffer([]byte{})
	utils.WriteUint32(cfcw, uint32(h.GetReceiveConnectionFlowControlWindow()))
	icsl := bytes.NewBuffer([]byte{})
	utils.Debugf("ICSL: %#v\n", h.GetIdleConnectionStateLifetime())
	utils.WriteUint32(icsl, uint32(h.GetIdleConnectionStateLifetime()/time.Second))

	return map[Tag][]byte{
		TagICSL: icsl.Bytes(),
		TagMSPC: []byte{0x64, 0x00, 0x00, 0x00}, //100
		TagCFCW: cfcw.Bytes(),
		TagSFCW: sfcw.Bytes(),
	}
}

// GetSendStreamFlowControlWindow gets the size of the stream-level flow control window for sending data
func (h *ConnectionParametersManager) GetSendStreamFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.sendStreamFlowControlWindow
}

// GetSendConnectionFlowControlWindow gets the size of the stream-level flow control window for sending data
func (h *ConnectionParametersManager) GetSendConnectionFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.sendConnectionFlowControlWindow
}

// GetReceiveStreamFlowControlWindow gets the size of the stream-level flow control window for receiving data
func (h *ConnectionParametersManager) GetReceiveStreamFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.receiveStreamFlowControlWindow
}

// GetReceiveConnectionFlowControlWindow gets the size of the stream-level flow control window for receiving data
func (h *ConnectionParametersManager) GetReceiveConnectionFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.receiveConnectionFlowControlWindow
}

// GetIdleConnectionStateLifetime gets the idle timeout
func (h *ConnectionParametersManager) GetIdleConnectionStateLifetime() time.Duration {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.idleConnectionStateLifetime
}

// TruncateConnectionID determines if the client requests truncated ConnectionIDs
func (h *ConnectionParametersManager) TruncateConnectionID() bool {
	rawValue, err := h.getRawValue(TagTCID)
	if err != nil {
		return false
	}

	var value uint32
	buf := bytes.NewBuffer(rawValue)
	err = binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		return false
	}

	if value == 0 {
		return true
	}
	return false
}
