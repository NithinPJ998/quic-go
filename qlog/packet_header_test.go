package qlog

import (
	"bytes"
	"encoding/json"

	"github.com/francoispqt/gojay"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packet Header", func() {
	It("determines the packet type from the encryption level", func() {
		Expect(getPacketTypeFromEncryptionLevel(protocol.EncryptionInitial)).To(BeEquivalentTo(logging.PacketTypeInitial))
		Expect(getPacketTypeFromEncryptionLevel(protocol.EncryptionHandshake)).To(BeEquivalentTo(logging.PacketTypeHandshake))
		Expect(getPacketTypeFromEncryptionLevel(protocol.Encryption0RTT)).To(BeEquivalentTo(logging.PacketType0RTT))
		Expect(getPacketTypeFromEncryptionLevel(protocol.Encryption1RTT)).To(BeEquivalentTo(logging.PacketType1RTT))
	})

	Context("marshalling", func() {
		check := func(hdr *wire.ExtendedHeader, expected map[string]interface{}) {
			buf := &bytes.Buffer{}
			enc := gojay.NewEncoder(buf)
			ExpectWithOffset(1, enc.Encode(transformExtendedHeader(hdr))).To(Succeed())
			data := buf.Bytes()
			ExpectWithOffset(1, json.Valid(data)).To(BeTrue())
			checkEncoding(data, expected)
		}

		It("marshals a header for a 1-RTT packet", func() {
			check(
				&wire.ExtendedHeader{PacketNumber: 42},
				map[string]interface{}{
					"packet_number": 42,
					"dcil":          0,
					"key_phase_bit": "0",
				},
			)
		})

		It("marshals a header with a payload length", func() {
			check(
				&wire.ExtendedHeader{
					PacketNumber: 42,
					Header: wire.Header{
						IsLongHeader: true,
						Type:         protocol.PacketTypeInitial,
						Length:       123,
						Version:      protocol.VersionNumber(0xdecafbad),
					},
				},
				map[string]interface{}{
					"packet_number":  42,
					"payload_length": 123,
					"dcil":           0,
					"scil":           0,
					"version":        "decafbad",
				},
			)
		})

		It("marshals a packet with packet number 0", func() {
			check(
				&wire.ExtendedHeader{
					PacketNumber: 0,
					Header: wire.Header{
						IsLongHeader: true,
						Type:         protocol.PacketTypeHandshake,
						Version:      protocol.VersionNumber(0xdecafbad),
					},
				},
				map[string]interface{}{
					"packet_number": 0,
					"dcil":          0,
					"scil":          0,
					"version":       "decafbad",
				},
			)
		})

		It("marshals a header with a source connection ID", func() {
			check(
				&wire.ExtendedHeader{
					PacketNumber: 42,
					Header: wire.Header{
						IsLongHeader:    true,
						Type:            protocol.PacketTypeHandshake,
						SrcConnectionID: protocol.ConnectionID{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
						Version:         protocol.VersionNumber(0xdecafbad),
					},
				},
				map[string]interface{}{
					"packet_number": 42,
					"dcil":          0,
					"scil":          16,
					"scid":          "00112233445566778899aabbccddeeff",
					"version":       "decafbad",
				},
			)
		})

		It("marshals a 1-RTT header with a destination connection ID", func() {
			check(
				&wire.ExtendedHeader{
					PacketNumber: 42,
					Header:       wire.Header{DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef}},
					KeyPhase:     protocol.KeyPhaseOne,
				},
				map[string]interface{}{
					"packet_number": 42,
					"dcil":          4,
					"dcid":          "deadbeef",
					"key_phase_bit": "1",
				},
			)
		})
	})
})
