package frames

import (
	"bytes"
	"fmt"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/utils"
)

// An AckFrame in QUIC
type AckFrame struct {
	Entropy         byte
	LargestObserved protocol.PacketNumber
	DelayTime       uint16 // Todo: properly interpret this value as described in the specification
}

// Write writes an ACK frame.
func (f *AckFrame) Write(b *bytes.Buffer) error {
	typeByte := uint8(0x48)
	b.WriteByte(typeByte)
	b.WriteByte(f.Entropy)
	utils.WriteUint32(b, uint32(f.LargestObserved)) // TODO: send the correct length
	utils.WriteUint16(b, 1)                         // TODO: Ack delay time
	b.WriteByte(0x01)                               // Just one timestamp
	b.WriteByte(0x00)                               // Largest observed
	utils.WriteUint32(b, 0)                         // First timestamp
	return nil
}

// ParseAckFrame reads an ACK frame
func ParseAckFrame(r *bytes.Reader) (*AckFrame, error) {
	frame := &AckFrame{}

	typeByte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	hasNACK := false
	if typeByte&0x20 == 0x20 {
		hasNACK = true
	}
	if typeByte&0x10 == 0x10 {
		panic("truncated ACKs not yet implemented.")
	}

	largestObservedLen := 2 * ((typeByte & 0x0C) >> 2)
	if largestObservedLen == 0 {
		largestObservedLen = 1
	}

	missingSequenceNumberDeltaLen := 2 * (typeByte & 0x03)
	if missingSequenceNumberDeltaLen == 0 {
		missingSequenceNumberDeltaLen = 1
	}
	_ = missingSequenceNumberDeltaLen

	frame.Entropy, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	largestObserved, err := utils.ReadUintN(r, largestObservedLen)
	if err != nil {
		return nil, err
	}
	frame.LargestObserved = protocol.PacketNumber(largestObserved)

	frame.DelayTime, err = utils.ReadUint16(r)
	if err != nil {
		return nil, err
	}

	numTimestampByte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	numTimestamp := uint8(numTimestampByte)

	// Delta Largest observed
	_, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	// First Timestamp
	_, err = utils.ReadUint32(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(numTimestamp)-1; i++ {
		// Delta Largest observed
		_, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		// Time Since Previous Timestamp
		_, err = utils.ReadUint16(r)
		if err != nil {
			return nil, err
		}
	}

	if hasNACK {
		fmt.Println("NACK not implemented yet!")
		var numRanges uint8
		numRanges, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		p := make([]byte, largestObservedLen+1)
		for i := uint8(0); i < numRanges; i++ {
			_, err := r.Read(p)
			if err != nil {
				return nil, err
			}
		}
	}

	return frame, nil
}
