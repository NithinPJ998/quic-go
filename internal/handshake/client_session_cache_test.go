package handshake

import (
	"bytes"
	"crypto/tls"
	"time"
	"unsafe"

	"github.com/lucas-clemente/quic-go/internal/utils"

	"github.com/marten-seemann/qtls"

	"github.com/lucas-clemente/quic-go/internal/congestion"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientSessionCache", func() {
	encodeIntoSessionTicket := func(data []byte) *tls.ClientSessionState {
		var session clientSessionState
		sessBytes := (*[unsafe.Sizeof(session)]byte)(unsafe.Pointer(&session))[:]
		session.nonce = data
		var tlsSession tls.ClientSessionState
		tlsSessBytes := (*[unsafe.Sizeof(tlsSession)]byte)(unsafe.Pointer(&tlsSession))[:]
		copy(tlsSessBytes, sessBytes)
		return &tlsSession
	}

	It("puts and gets", func() {
		get := make(chan []byte, 100)
		set := make(chan []byte, 100)

		csc := newClientSessionCache(
			tls.NewLRUClientSessionCache(100),
			congestion.NewRTTStats(),
			func() []byte { return <-get },
			func(b []byte) { set <- b },
		)

		get <- []byte("foobar")
		csc.Put("localhost", &qtls.ClientSessionState{})
		Expect(set).To(BeEmpty())
		state, ok := csc.Get("localhost")
		Expect(ok).To(BeTrue())
		Expect(state).ToNot(BeNil())
		Expect(set).To(Receive(Equal([]byte("foobar"))))
	})

	It("saves the RTT", func() {
		rttStatsOrig := congestion.NewRTTStats()
		rttStatsOrig.UpdateRTT(10*time.Second, 0, time.Now())
		Expect(rttStatsOrig.SmoothedRTT()).To(Equal(10 * time.Second))
		cache := tls.NewLRUClientSessionCache(100)
		csc1 := newClientSessionCache(
			cache,
			rttStatsOrig,
			func() []byte { return nil },
			func([]byte) {},
		)
		csc1.Put("localhost", &qtls.ClientSessionState{})

		rttStats := congestion.NewRTTStats()
		csc2 := newClientSessionCache(
			cache,
			rttStats,
			func() []byte { return nil },
			func([]byte) {},
		)
		Expect(rttStats.SmoothedRTT()).ToNot(Equal(10 * time.Second))
		_, ok := csc2.Get("localhost")
		Expect(ok).To(BeTrue())
		Expect(rttStats.SmoothedRTT()).To(Equal(10 * time.Second))
	})

	It("refuses a session state that is too short for the revision", func() {
		cache := tls.NewLRUClientSessionCache(1)
		cache.Put("localhost", encodeIntoSessionTicket([]byte{}))

		csc := newClientSessionCache(
			cache,
			congestion.NewRTTStats(),
			func() []byte { return nil },
			func([]byte) {},
		)
		_, ok := csc.Get("localhost")
		Expect(ok).To(BeFalse())
	})

	It("refuses a session state with the wrong revision", func() {
		cache := tls.NewLRUClientSessionCache(1)
		b := &bytes.Buffer{}
		utils.WriteVarInt(b, clientSessionStateRevision+1)
		cache.Put("localhost", encodeIntoSessionTicket(b.Bytes()))

		csc := newClientSessionCache(
			cache,
			congestion.NewRTTStats(),
			func() []byte { return nil },
			func([]byte) {},
		)
		_, ok := csc.Get("localhost")
		Expect(ok).To(BeFalse())
	})

	It("refuses a session state when unmarshalling fails", func() {
		cache := tls.NewLRUClientSessionCache(1)
		b := &bytes.Buffer{}
		utils.WriteVarInt(b, clientSessionStateRevision)
		b.Write([]byte("foobar"))
		cache.Put("localhost", encodeIntoSessionTicket(b.Bytes()))

		csc := newClientSessionCache(
			cache,
			congestion.NewRTTStats(),
			func() []byte { return nil },
			func([]byte) {},
		)
		_, ok := csc.Get("localhost")
		Expect(ok).To(BeFalse())
	})
})
