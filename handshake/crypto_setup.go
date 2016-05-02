package handshake

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"sync"

	"github.com/lucas-clemente/quic-go/crypto"
	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/utils"
)

// KeyDerivationFunction is used for key derivation
type KeyDerivationFunction func(forwardSecure bool, sharedSecret, nonces []byte, connID protocol.ConnectionID, chlo []byte, scfg []byte, cert []byte) (crypto.AEAD, error)

// The CryptoSetup handles all things crypto for the Session
type CryptoSetup struct {
	connID  protocol.ConnectionID
	version protocol.VersionNumber
	scfg    *ServerConfig
	nonce   []byte

	secureAEAD                  crypto.AEAD
	forwardSecureAEAD           crypto.AEAD
	receivedForwardSecurePacket bool
	receivedSecurePacket        bool

	keyDerivation KeyDerivationFunction

	cryptoStream utils.Stream

	mutex sync.RWMutex
}

var _ crypto.AEAD = &CryptoSetup{}

// NewCryptoSetup creates a new CryptoSetup instance
func NewCryptoSetup(connID protocol.ConnectionID, version protocol.VersionNumber, scfg *ServerConfig, cryptoStream utils.Stream) *CryptoSetup {
	nonce := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return &CryptoSetup{
		connID:        connID,
		version:       version,
		scfg:          scfg,
		nonce:         nonce,
		keyDerivation: crypto.DeriveKeysChacha20,
		cryptoStream:  cryptoStream,
	}
}

// HandleCryptoStream reads and writes messages on the crypto stream
func (h *CryptoSetup) HandleCryptoStream() error {
	for {
		cachingReader := utils.NewCachingReader(h.cryptoStream)
		messageTag, cryptoData, err := ParseHandshakeMessage(cachingReader)
		if err != nil {
			return err
		}
		if messageTag != TagCHLO {
			return errors.New("CryptoSetup: expected CHLO")
		}
		chloData := cachingReader.Get()

		var reply []byte
		if !h.isInchoateCHLO(cryptoData) {
			// We have a CHLO with a proper server config ID, do a 0-RTT handshake
			reply, err = h.handleCHLO(chloData, cryptoData)
			if err != nil {
				return err
			}
			_, err = h.cryptoStream.Write(reply)
			if err != nil {
				return err
			}
			return nil
		}

		// We have an inchoate or non-matching CHLO, we now send a rejection
		reply, err = h.handleInchoateCHLO(chloData)
		if err != nil {
			return err
		}
		_, err = h.cryptoStream.Write(reply)
		if err != nil {
			return err
		}
	}
}

// Open a message
func (h *CryptoSetup) Open(packetNumber protocol.PacketNumber, associatedData []byte, ciphertext []byte) ([]byte, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if h.forwardSecureAEAD != nil {
		res, err := h.forwardSecureAEAD.Open(packetNumber, associatedData, ciphertext)
		if err == nil {
			h.receivedForwardSecurePacket = true
			return res, nil
		}
		if h.receivedForwardSecurePacket {
			return nil, err
		}
	}
	if h.secureAEAD != nil {
		res, err := h.secureAEAD.Open(packetNumber, associatedData, ciphertext)
		if err == nil {
			h.receivedSecurePacket = true
			return res, nil
		}
		if h.receivedSecurePacket {
			return nil, err
		}
	}
	return (&crypto.NullAEAD{}).Open(packetNumber, associatedData, ciphertext)
}

// Seal a messageTag
func (h *CryptoSetup) Seal(packetNumber protocol.PacketNumber, associatedData []byte, plaintext []byte) []byte {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if h.receivedForwardSecurePacket {
		return h.forwardSecureAEAD.Seal(packetNumber, associatedData, plaintext)
	} else if h.secureAEAD != nil {
		return h.secureAEAD.Seal(packetNumber, associatedData, plaintext)
	} else {
		return (&crypto.NullAEAD{}).Seal(packetNumber, associatedData, plaintext)
	}
}

func (h *CryptoSetup) isInchoateCHLO(cryptoData map[Tag][]byte) bool {
	scid, ok := cryptoData[TagSCID]
	if !ok || !bytes.Equal(h.scfg.ID, scid) {
		return true
	}
	sno, ok := cryptoData[TagSNO]
	if !ok || !bytes.Equal(h.nonce, sno) {
		return true
	}
	return false
}

func (h *CryptoSetup) handleInchoateCHLO(data []byte) ([]byte, error) {
	proof, err := h.scfg.Sign(data)
	if err != nil {
		return nil, err
	}

	var serverReply bytes.Buffer
	WriteHandshakeMessage(&serverReply, TagREJ, map[Tag][]byte{
		TagSCFG: h.scfg.Get(),
		TagCERT: h.scfg.GetCertCompressed(),
		TagSNO:  h.nonce,
		TagPROF: proof,
	})
	return serverReply.Bytes(), nil
}

func (h *CryptoSetup) handleCHLO(data []byte, cryptoData map[Tag][]byte) ([]byte, error) {
	// We have a CHLO matching our server config, we can continue with the 0-RTT handshake
	sharedSecret, err := h.scfg.kex.CalculateSharedKey(cryptoData[TagPUBS])
	if err != nil {
		return nil, err
	}
	var nonce bytes.Buffer
	nonce.Write(cryptoData[TagNONC])
	nonce.Write(h.nonce)

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.secureAEAD, err = h.keyDerivation(false, sharedSecret, nonce.Bytes(), h.connID, data, h.scfg.Get(), h.scfg.signer.GetCertUncompressed())
	if err != nil {
		return nil, err
	}
	// TODO: Use new curve
	h.forwardSecureAEAD, err = h.keyDerivation(true, sharedSecret, nonce.Bytes(), h.connID, data, h.scfg.Get(), h.scfg.signer.GetCertUncompressed())
	if err != nil {
		return nil, err
	}

	var reply bytes.Buffer
	WriteHandshakeMessage(&reply, TagSHLO, map[Tag][]byte{
		TagPUBS: h.scfg.kex.PublicKey(),
		TagSNO:  h.nonce,
		TagVER:  protocol.SupportedVersionsAsTags,
		TagICSL: []byte{0x1e, 0x00, 0x00, 0x00}, //30
		TagMSPC: []byte{0x64, 0x00, 0x00, 0x00}, //100
	})
	return reply.Bytes(), nil
}
