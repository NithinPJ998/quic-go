package crypto

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"strings"

	"github.com/lucas-clemente/quic-go/utils"
)

// rsaSigner stores a key and a certificate for the server proof
type rsaSigner struct {
	key    *rsa.PrivateKey
	cert   *x509.Certificate
	config *tls.Config
}

// NewRSASigner loads the key and cert from files
func NewRSASigner(tlsConfig *tls.Config) (Signer, error) {
	if len(tlsConfig.Certificates) == 0 {
		return nil, errors.New("Expected at least one certificate in TLS config")
	}
	cert := tlsConfig.Certificates[0]

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	rsaKey, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Only RSA private keys are supported for now")
	}

	return &rsaSigner{key: rsaKey, cert: x509Cert, config: tlsConfig}, nil
}

// SignServerProof signs CHLO and server config for use in the server proof
func (kd *rsaSigner) SignServerProof(sni string, chlo []byte, serverConfigData []byte) ([]byte, error) {
	hash := sha256.New()
	if len(chlo) > 0 {
		// Version >= 31
		hash.Write([]byte("QUIC CHLO and server config signature\x00"))
		chloHash := sha256.Sum256(chlo)
		hash.Write([]byte{32, 0, 0, 0})
		hash.Write(chloHash[:])
	} else {
		hash.Write([]byte("QUIC server config signature\x00"))
	}
	hash.Write(serverConfigData)
	return rsa.SignPSS(rand.Reader, kd.key, crypto.SHA256, hash.Sum(nil), &rsa.PSSOptions{SaltLength: 32})
}

// GetCertCompressed gets the certificate in the format described by the QUIC crypto doc
func (kd *rsaSigner) GetCertCompressed(sni string) ([]byte, error) {
	b := &bytes.Buffer{}
	b.WriteByte(1) // Entry type compressed
	b.WriteByte(0) // Entry type end_of_list
	utils.WriteUint32(b, uint32(len(kd.cert.Raw)+4))
	gz, err := zlib.NewWriterLevelDict(b, flate.BestCompression, certDictZlib)
	if err != nil {
		panic(err)
	}
	lenCert := len(kd.cert.Raw)
	gz.Write([]byte{
		byte(lenCert & 0xff),
		byte((lenCert >> 8) & 0xff),
		byte((lenCert >> 16) & 0xff),
		byte((lenCert >> 24) & 0xff),
	})
	gz.Write(kd.cert.Raw)
	gz.Close()
	return b.Bytes(), nil
}

// GetCertUncompressed gets the certificate in DER
func (kd *rsaSigner) GetCertUncompressed(sni string) ([]byte, error) {
	return kd.cert.Raw, nil
}

func (kd *rsaSigner) getCertForSNI(sni string) (*tls.Certificate, error) {
	if kd.config.GetCertificate != nil {
		cert, err := kd.config.GetCertificate(&tls.ClientHelloInfo{ServerName: sni})
		if err != nil {
			return nil, err
		}
		if cert != nil {
			return cert, nil
		}
	}
	if len(kd.config.NameToCertificate) != 0 {
		if cert, ok := kd.config.NameToCertificate[sni]; ok {
			return cert, nil
		}
		wildcardSNI := "*" + strings.TrimLeftFunc(sni, func(r rune) bool { return r != '.' })
		if cert, ok := kd.config.NameToCertificate[wildcardSNI]; ok {
			return cert, nil
		}
	}
	if len(kd.config.Certificates) != 0 {
		return &kd.config.Certificates[0], nil
	}
	return nil, errors.New("no matching certificate found")
}
