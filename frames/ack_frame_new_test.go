package frames

import (
	"bytes"
	"time"

	"github.com/lucas-clemente/quic-go/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AckFrame", func() {
	Context("when parsing", func() {
		It("accepts a sample frame", func() {
			b := bytes.NewReader([]byte{0x40, 0x1c, 0x8e, 0x0, 0x1c, 0x1, 0x1, 0x6b, 0x26, 0x3, 0x0})
			frame, err := ParseAckFrameNew(b, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x1c)))
			Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
			Expect(frame.DelayTime).To(Equal(142 * time.Microsecond))
			Expect(frame.HasMissingRanges()).To(BeFalse())
			Expect(b.Len()).To(BeZero())
		})

		It("parses a frame without a timestamp", func() {
			b := bytes.NewReader([]byte{0x40, 0x3, 0x50, 0x15, 0x3, 0x0})
			frame, err := ParseAckFrameNew(b, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(3)))
		})

		It("parses a frame with a 48 bit packet number", func() {
			b := bytes.NewReader([]byte{0x4c, 0x37, 0x13, 0xad, 0xfb, 0xca, 0xde, 0x0, 0x0, 0x0, 0x1, 0, 0, 0, 0, 0})
			frame, err := ParseAckFrameNew(b, protocol.Version34)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0xdecafbad1337)))
			Expect(frame.HasMissingRanges()).To(BeFalse())
			Expect(b.Len()).To(BeZero())
		})

		It("parses a frame, when packet 1 was lost", func() {
			b := bytes.NewReader([]byte{0x40, 0x9, 0x92, 0x7, 0x8, 0x3, 0x2, 0x69, 0xa3, 0x0, 0x0, 0x1, 0xc9, 0x2, 0x0, 0x46, 0x10})
			frame, err := ParseAckFrameNew(b, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(9)))
			Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(2)))
			Expect(b.Len()).To(BeZero())
		})

		It("parses a frame with multiple timestamps", func() {
			b := bytes.NewReader([]byte{0x40, 0x10, 0x0, 0x0, 0x10, 0x4, 0x1, 0x6b, 0x26, 0x4, 0x0, 0x3, 0, 0, 0x2, 0, 0, 0x1, 0, 0})
			_, err := ParseAckFrameNew(b, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(b.Len()).To(BeZero())
		})

		It("errors when the ACK range is too large", func() {
			// LargestAcked: 0x1c
			// Length: 0x1d => LowestAcked would be -1
			b := bytes.NewReader([]byte{0x40, 0x1c, 0x8e, 0x0, 0x1d, 0x1, 0x1, 0x6b, 0x26, 0x3, 0x0})
			_, err := ParseAckFrameNew(b, 0)
			Expect(err).To(MatchError(ErrInvalidAckRanges))
		})

		Context("ACK blocks", func() {
			It("parses a frame with one ACK block", func() {
				b := bytes.NewReader([]byte{0x60, 0x18, 0x94, 0x1, 0x1, 0x3, 0x2, 0x10, 0x2, 0x1, 0x5c, 0xd5, 0x0, 0x0, 0x0, 0x95, 0x0})
				frame, err := ParseAckFrameNew(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(24)))
				Expect(frame.HasMissingRanges()).To(BeTrue())
				Expect(frame.AckRanges).To(HaveLen(2))
				Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 22, LastPacketNumber: 24}))
				Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 4, LastPacketNumber: 19}))
				Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(4)))
				Expect(b.Len()).To(BeZero())
			})

			It("rejects a frame with invalid ACK ranges", func() {
				// like the test before, but increased the last ACK range, such that the FirstPacketNumber would be negative
				b := bytes.NewReader([]byte{0x60, 0x18, 0x94, 0x1, 0x1, 0x3, 0x2, 0x15, 0x2, 0x1, 0x5c, 0xd5, 0x0, 0x0, 0x0, 0x95, 0x0})
				_, err := ParseAckFrameNew(b, 0)
				Expect(err).To(MatchError(ErrInvalidAckRanges))
			})

			It("parses a frame with multiple single packets missing", func() {
				b := bytes.NewReader([]byte{0x60, 0x27, 0xda, 0x0, 0x6, 0x9, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x13, 0x2, 0x1, 0x71, 0x12, 0x3, 0x0, 0x0, 0x47, 0x2})
				frame, err := ParseAckFrameNew(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x27)))
				Expect(frame.HasMissingRanges()).To(BeTrue())
				Expect(frame.AckRanges).To(HaveLen(7))
				Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 31, LastPacketNumber: 0x27}))
				Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 29, LastPacketNumber: 29}))
				Expect(frame.AckRanges[2]).To(Equal(AckRange{FirstPacketNumber: 27, LastPacketNumber: 27}))
				Expect(frame.AckRanges[3]).To(Equal(AckRange{FirstPacketNumber: 25, LastPacketNumber: 25}))
				Expect(frame.AckRanges[4]).To(Equal(AckRange{FirstPacketNumber: 23, LastPacketNumber: 23}))
				Expect(frame.AckRanges[5]).To(Equal(AckRange{FirstPacketNumber: 21, LastPacketNumber: 21}))
				Expect(frame.AckRanges[6]).To(Equal(AckRange{FirstPacketNumber: 1, LastPacketNumber: 19}))
				Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
				Expect(b.Len()).To(BeZero())
			})

			It("parses a packet with packet 1 and one more packet lost", func() {
				b := bytes.NewReader([]byte{0x60, 0xc, 0x92, 0x0, 0x1, 0x1, 0x1, 0x9, 0x2, 0x2, 0x53, 0x43, 0x1, 0x0, 0x0, 0xa7, 0x0})
				frame, err := ParseAckFrameNew(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(12)))
				Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(2)))
				Expect(frame.AckRanges).To(HaveLen(2))
				Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 12, LastPacketNumber: 12}))
				Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 2, LastPacketNumber: 10}))
				Expect(b.Len()).To(BeZero())
			})

			It("parses a frame with multiple longer ACK blocks", func() {
				b := bytes.NewReader([]byte{0x60, 0x52, 0xd1, 0x0, 0x3, 0x17, 0xa, 0x10, 0x4, 0x8, 0x2, 0x12, 0x2, 0x1, 0x6c, 0xc8, 0x2, 0x0, 0x0, 0x7e, 0x1})
				frame, err := ParseAckFrameNew(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x52)))
				Expect(frame.HasMissingRanges()).To(BeTrue())
				Expect(frame.AckRanges).To(HaveLen(4))
				Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 60, LastPacketNumber: 0x52}))
				Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 34, LastPacketNumber: 49}))
				Expect(frame.AckRanges[2]).To(Equal(AckRange{FirstPacketNumber: 22, LastPacketNumber: 29}))
				Expect(frame.AckRanges[3]).To(Equal(AckRange{FirstPacketNumber: 2, LastPacketNumber: 19}))
				Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(2)))
				Expect(b.Len()).To(BeZero())
			})

			Context("more than 256 lost packets in a row", func() {
				It("parses a frame with one long range, spanning 2 blocks, of missing packets", func() { // 280 missing packets
					b := bytes.NewReader([]byte{0x64, 0x44, 0x1, 0xa7, 0x0, 0x2, 0x19, 0xff, 0x0, 0x19, 0x13, 0x2, 0x1, 0xb, 0x59, 0x2, 0x0, 0x0, 0xb6, 0x0})
					frame, err := ParseAckFrameNew(b, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x144)))
					Expect(frame.HasMissingRanges()).To(BeTrue())
					Expect(frame.AckRanges).To(HaveLen(2))
					Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 300, LastPacketNumber: 0x144}))
					Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 1, LastPacketNumber: 19}))
					Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
					Expect(b.Len()).To(BeZero())
				})

				It("parses a frame with one long range, spanning mulitple blocks, of missing packets", func() { // 2345 missing packets
					b := bytes.NewReader([]byte{0x64, 0x5b, 0x9, 0x66, 0x1, 0xa, 0x1f, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0xff, 0x0, 0x32, 0x13, 0x4, 0x3, 0xb4, 0xda, 0x1, 0x0, 0x2, 0xe0, 0x0, 0x1, 0x9a, 0x0, 0x0, 0x81, 0x0})
					frame, err := ParseAckFrameNew(b, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x95b)))
					Expect(frame.HasMissingRanges()).To(BeTrue())
					Expect(frame.AckRanges).To(HaveLen(2))
					Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 2365, LastPacketNumber: 0x95b}))
					Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 1, LastPacketNumber: 19}))
					Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
					Expect(b.Len()).To(BeZero())
				})

				It("parses a frame with multiple long ranges of missing packets", func() {
					b := bytes.NewReader([]byte{0x65, 0x66, 0x9, 0x23, 0x1, 0x7, 0x7, 0x0, 0xff, 0x0, 0x0, 0xf5, 0x8a, 0x2, 0xc8, 0xe6, 0x0, 0xff, 0x0, 0x0, 0xff, 0x0, 0x0, 0xff, 0x0, 0x0, 0x23, 0x13, 0x0, 0x2, 0x1, 0x13, 0xae, 0xb, 0x0, 0x0, 0x80, 0x5})
					frame, err := ParseAckFrameNew(b, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x966)))
					Expect(frame.HasMissingRanges()).To(BeTrue())
					Expect(frame.AckRanges).To(HaveLen(4))
					Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 2400, LastPacketNumber: 0x966}))
					Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 1250, LastPacketNumber: 1899}))
					Expect(frame.AckRanges[2]).To(Equal(AckRange{FirstPacketNumber: 820, LastPacketNumber: 1049}))
					Expect(frame.AckRanges[3]).To(Equal(AckRange{FirstPacketNumber: 1, LastPacketNumber: 19}))
					Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
					Expect(b.Len()).To(BeZero())
				})

				It("parses a frame with short ranges and one long range", func() {
					b := bytes.NewReader([]byte{0x64, 0x8f, 0x3, 0x65, 0x1, 0x5, 0x3d, 0x1, 0x32, 0xff, 0x0, 0xff, 0x0, 0xf0, 0x1c, 0x2, 0x13, 0x3, 0x2, 0x23, 0xaf, 0x2, 0x0, 0x1, 0x3, 0x1, 0x0, 0x8e, 0x0})
					frame, err := ParseAckFrameNew(b, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(protocol.PacketNumber(0x38f)))
					Expect(frame.HasMissingRanges()).To(BeTrue())
					Expect(frame.AckRanges).To(HaveLen(4))
					Expect(frame.AckRanges[0]).To(Equal(AckRange{FirstPacketNumber: 851, LastPacketNumber: 0x38f}))
					Expect(frame.AckRanges[1]).To(Equal(AckRange{FirstPacketNumber: 800, LastPacketNumber: 849}))
					Expect(frame.AckRanges[2]).To(Equal(AckRange{FirstPacketNumber: 22, LastPacketNumber: 49}))
					Expect(frame.AckRanges[3]).To(Equal(AckRange{FirstPacketNumber: 1, LastPacketNumber: 19}))
					Expect(frame.LowestAcked).To(Equal(protocol.PacketNumber(1)))
					Expect(b.Len()).To(BeZero())
				})
			})
		})
	})

	Context("when writing", func() {
		var b *bytes.Buffer

		BeforeEach(func() {
			b = &bytes.Buffer{}
		})

		Context("self-consistency", func() {
			It("writes a simple ACK frame", func() {
				frameOrig := &AckFrameNew{
					LargestAcked: 1,
					LowestAcked:  1,
				}
				err := frameOrig.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				r := bytes.NewReader(b.Bytes())
				frame, err := ParseAckFrameNew(r, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
				Expect(frame.HasMissingRanges()).To(BeFalse())
				Expect(r.Len()).To(BeZero())
			})

			It("writes the correct block length in a simple ACK frame", func() {
				frameOrig := &AckFrameNew{
					LargestAcked: 20,
					LowestAcked:  10,
				}
				err := frameOrig.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				r := bytes.NewReader(b.Bytes())
				frame, err := ParseAckFrameNew(r, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
				Expect(frame.LowestAcked).To(Equal(frameOrig.LowestAcked))
				Expect(frame.HasMissingRanges()).To(BeFalse())
				Expect(r.Len()).To(BeZero())
			})

			It("writes a simple ACK frame with a high packet number", func() {
				frameOrig := &AckFrameNew{
					LargestAcked: 0xDEADBEEFCAFE,
					LowestAcked:  0xDEADBEEFCAFE,
				}
				err := frameOrig.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				r := bytes.NewReader(b.Bytes())
				frame, err := ParseAckFrameNew(r, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
				Expect(frame.HasMissingRanges()).To(BeFalse())
				Expect(r.Len()).To(BeZero())
			})

			It("writes an ACK frame with one packet missing", func() {
				frameOrig := &AckFrameNew{
					LargestAcked: 40,
					LowestAcked:  1,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 25, LastPacketNumber: 40},
						AckRange{FirstPacketNumber: 1, LastPacketNumber: 23},
					},
				}
				err := frameOrig.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				r := bytes.NewReader(b.Bytes())
				frame, err := ParseAckFrameNew(r, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
				Expect(frame.LowestAcked).To(Equal(frameOrig.LowestAcked))
				Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				Expect(r.Len()).To(BeZero())
			})

			It("writes an ACK frame with multiple missing packets", func() {
				frameOrig := &AckFrameNew{
					LargestAcked: 25,
					LowestAcked:  1,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 22, LastPacketNumber: 25},
						AckRange{FirstPacketNumber: 15, LastPacketNumber: 18},
						AckRange{FirstPacketNumber: 13, LastPacketNumber: 13},
						AckRange{FirstPacketNumber: 1, LastPacketNumber: 10},
					},
				}
				err := frameOrig.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				r := bytes.NewReader(b.Bytes())
				frame, err := ParseAckFrameNew(r, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
				Expect(frame.LowestAcked).To(Equal(frameOrig.LowestAcked))
				Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				Expect(r.Len()).To(BeZero())
			})

			It("rejects a frame with incorrect LargestObserved value", func() {
				frame := &AckFrameNew{
					LargestAcked: 26,
					LowestAcked:  1,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 12, LastPacketNumber: 25},
						AckRange{FirstPacketNumber: 1, LastPacketNumber: 10},
					},
				}
				err := frame.Write(b, 0)
				Expect(err).To(MatchError(errInconsistentAckLargestAcked))
			})

			It("rejects a frame with incorrect LargestObserved value", func() {
				frame := &AckFrameNew{
					LargestAcked: 25,
					LowestAcked:  2,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 12, LastPacketNumber: 25},
						AckRange{FirstPacketNumber: 1, LastPacketNumber: 10},
					},
				}
				err := frame.Write(b, 0)
				Expect(err).To(MatchError(errInconsistentAckLowestAcked))
			})

			Context("longer ACK blocks", func() {
				It("only writes one block for 254 lost packets", func() {
					frameOrig := &AckFrameNew{
						LargestAcked: 300,
						LowestAcked:  1,
						AckRanges: []AckRange{
							AckRange{FirstPacketNumber: 20 + 254, LastPacketNumber: 300},
							AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
						},
					}
					Expect(frameOrig.numWrittenNackRanges()).To(Equal(uint64(2)))
					err := frameOrig.Write(b, 0)
					Expect(err).ToNot(HaveOccurred())
					r := bytes.NewReader(b.Bytes())
					frame, err := ParseAckFrameNew(r, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
					Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				})

				It("only writes one block for 255 lost packets", func() {
					frameOrig := &AckFrameNew{
						LargestAcked: 300,
						LowestAcked:  1,
						AckRanges: []AckRange{
							AckRange{FirstPacketNumber: 20 + 255, LastPacketNumber: 300},
							AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
						},
					}
					Expect(frameOrig.numWrittenNackRanges()).To(Equal(uint64(2)))
					err := frameOrig.Write(b, 0)
					Expect(err).ToNot(HaveOccurred())
					r := bytes.NewReader(b.Bytes())
					frame, err := ParseAckFrameNew(r, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
					Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				})

				It("writes two blocks for 256 lost packets", func() {
					frameOrig := &AckFrameNew{
						LargestAcked: 300,
						LowestAcked:  1,
						AckRanges: []AckRange{
							AckRange{FirstPacketNumber: 20 + 256, LastPacketNumber: 300},
							AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
						},
					}
					Expect(frameOrig.numWrittenNackRanges()).To(Equal(uint64(3)))
					err := frameOrig.Write(b, 0)
					Expect(err).ToNot(HaveOccurred())
					// Expect(b.Bytes()[13+0*(1+6) : 13+1*(1+6)]).To(Equal([]byte{0xFF, 0, 0, 0, 0, 0, 0}))
					// Expect(b.Bytes()[13+1*(1+6) : 13+2*(1+6)]).To(Equal([]byte{0x1, 0, 0, 0, 0, 0, 19}))
					r := bytes.NewReader(b.Bytes())
					frame, err := ParseAckFrameNew(r, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
					Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				})

				It("writes multiple blocks for a lot of lost packets", func() {
					frameOrig := &AckFrameNew{
						LargestAcked: 3000,
						LowestAcked:  1,
						AckRanges: []AckRange{
							AckRange{FirstPacketNumber: 2900, LastPacketNumber: 3000},
							AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
						},
					}
					err := frameOrig.Write(b, 0)
					Expect(err).ToNot(HaveOccurred())
					r := bytes.NewReader(b.Bytes())
					frame, err := ParseAckFrameNew(r, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
					Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				})

				It("writes multiple longer blocks for 256 lost packets", func() {
					frameOrig := &AckFrameNew{
						LargestAcked: 3600,
						LowestAcked:  1,
						AckRanges: []AckRange{
							AckRange{FirstPacketNumber: 2900, LastPacketNumber: 3600},
							AckRange{FirstPacketNumber: 1000, LastPacketNumber: 2500},
							AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
						},
					}
					err := frameOrig.Write(b, 0)
					Expect(err).ToNot(HaveOccurred())
					r := bytes.NewReader(b.Bytes())
					frame, err := ParseAckFrameNew(r, 0)
					Expect(err).ToNot(HaveOccurred())
					Expect(frame.LargestAcked).To(Equal(frameOrig.LargestAcked))
					Expect(frame.AckRanges).To(Equal(frameOrig.AckRanges))
				})
			})
		})

		Context("min length", func() {
			It("has proper min length", func() {
				f := &AckFrameNew{
					LargestAcked: 1,
				}
				err := f.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(f.MinLength(0)).To(Equal(protocol.ByteCount(b.Len())))
			})

			It("has proper min length with a large LargestObserved", func() {
				f := &AckFrameNew{
					LargestAcked: 0xDEADBEEFCAFE,
				}
				err := f.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(f.MinLength(0)).To(Equal(protocol.ByteCount(b.Len())))
			})

			It("has the proper min length for an ACK with missing packets", func() {
				f := &AckFrameNew{
					LargestAcked: 2000,
					LowestAcked:  10,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 1000, LastPacketNumber: 2000},
						AckRange{FirstPacketNumber: 50, LastPacketNumber: 900},
						AckRange{FirstPacketNumber: 10, LastPacketNumber: 23},
					},
				}
				err := f.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(f.MinLength(0)).To(Equal(protocol.ByteCount(b.Len())))
			})

			It("has the proper min length for an ACK with long gaps of missing packets", func() {
				f := &AckFrameNew{
					LargestAcked: 2000,
					LowestAcked:  1,
					AckRanges: []AckRange{
						AckRange{FirstPacketNumber: 1500, LastPacketNumber: 2000},
						AckRange{FirstPacketNumber: 290, LastPacketNumber: 295},
						AckRange{FirstPacketNumber: 1, LastPacketNumber: 19},
					},
				}
				err := f.Write(b, 0)
				Expect(err).ToNot(HaveOccurred())
				Expect(f.MinLength(0)).To(Equal(protocol.ByteCount(b.Len())))
			})
		})
	})

	Context("ACK range validator", func() {
		It("accepts an ACK without NACK Ranges", func() {
			ack := AckFrameNew{LargestAcked: 7}
			Expect(ack.validateAckRanges()).To(BeTrue())
		})

		It("rejects ACK ranges with a single range", func() {
			ack := AckFrameNew{
				LargestAcked: 10,
				AckRanges:    []AckRange{AckRange{FirstPacketNumber: 1, LastPacketNumber: 10}},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects ACK ranges with LastPacketNumber of the first range unequal to LargestObserved", func() {
			ack := AckFrameNew{
				LargestAcked: 10,
				AckRanges: []AckRange{
					AckRange{FirstPacketNumber: 8, LastPacketNumber: 9},
					AckRange{FirstPacketNumber: 2, LastPacketNumber: 3},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects ACK ranges with FirstPacketNumber greater than LastPacketNumber", func() {
			ack := AckFrameNew{
				LargestAcked: 10,
				AckRanges: []AckRange{
					AckRange{FirstPacketNumber: 8, LastPacketNumber: 10},
					AckRange{FirstPacketNumber: 4, LastPacketNumber: 3},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects ACK ranges with FirstPacketNumber greater than LargestObserved", func() {
			ack := AckFrameNew{
				LargestAcked: 5,
				AckRanges: []AckRange{
					AckRange{FirstPacketNumber: 4, LastPacketNumber: 10},
					AckRange{FirstPacketNumber: 1, LastPacketNumber: 2},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects ACK ranges in the wrong order", func() {
			ack := AckFrameNew{
				LargestAcked: 7,
				AckRanges: []AckRange{
					{FirstPacketNumber: 2, LastPacketNumber: 2},
					{FirstPacketNumber: 6, LastPacketNumber: 7},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects with overlapping ACK ranges", func() {
			ack := AckFrameNew{
				LargestAcked: 7,
				AckRanges: []AckRange{
					{FirstPacketNumber: 5, LastPacketNumber: 7},
					{FirstPacketNumber: 2, LastPacketNumber: 5},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects ACK ranges that are part of a larger ACK range", func() {
			ack := AckFrameNew{
				LargestAcked: 7,
				AckRanges: []AckRange{
					{FirstPacketNumber: 4, LastPacketNumber: 7},
					{FirstPacketNumber: 5, LastPacketNumber: 6},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("rejects with directly adjacent ACK ranges", func() {
			ack := AckFrameNew{
				LargestAcked: 7,
				AckRanges: []AckRange{
					{FirstPacketNumber: 5, LastPacketNumber: 7},
					{FirstPacketNumber: 2, LastPacketNumber: 4},
				},
			}
			Expect(ack.validateAckRanges()).To(BeFalse())
		})

		It("accepts an ACK with one lost packet", func() {
			ack := AckFrameNew{
				LargestAcked: 10,
				AckRanges: []AckRange{
					{FirstPacketNumber: 5, LastPacketNumber: 10},
					{FirstPacketNumber: 1, LastPacketNumber: 3},
				},
			}
			Expect(ack.validateAckRanges()).To(BeTrue())
		})

		It("accepts an ACK with multiple lost packets", func() {
			ack := AckFrameNew{
				LargestAcked: 20,
				AckRanges: []AckRange{
					{FirstPacketNumber: 15, LastPacketNumber: 20},
					{FirstPacketNumber: 10, LastPacketNumber: 12},
					{FirstPacketNumber: 1, LastPacketNumber: 3},
				},
			}
			Expect(ack.validateAckRanges()).To(BeTrue())
		})
	})
})
