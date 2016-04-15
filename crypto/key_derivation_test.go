package crypto

import (
	"github.com/lucas-clemente/quic-go/protocol"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyDerivation", func() {
	It("derives proper keys", func() {
		aead, err := DeriveKeysChacha20(
			[]byte("0123456789012345678901"),
			[]byte("nonce"),
			protocol.ConnectionID(42),
			[]byte("chlo"),
			[]byte("scfg"),
			[]byte("cert"),
		)
		Expect(err).ToNot(HaveOccurred())
		chacha := aead.(*aeadChacha20Poly1305)
		// If the IVs match, the keys will match too, since the keys are read earlier
		Expect(chacha.myIV).To(Equal([]byte{0xf0, 0xf5, 0x4c, 0xa8}))
		Expect(chacha.otherIV).To(Equal([]byte{0x75, 0xd8, 0xa2, 0x8d}))
	})
})
