package mqtt

import (
	"net"
	"testing"
	"time"

	"github.com/256dpi/gomqtt/packet"
	"github.com/256dpi/gomqtt/transport"
	"github.com/256dpi/gomqtt/transport/flow"
	"github.com/creasty/defaults"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	t                      *testing.T
	expectedError          string
	expectedProcessPublish func(*packet.Publish) error
	expectedProcessPuback  func(*packet.Puback) error
}

func (h *mockHandler) ProcessPublish(pkt *packet.Publish) error {
	if h.expectedProcessPublish != nil {
		return h.expectedProcessPublish(pkt)
	}
	return nil
}

func (h *mockHandler) ProcessPuback(pkt *packet.Puback) error {
	if h.expectedProcessPuback != nil {
		return h.expectedProcessPuback(pkt)
	}
	return nil
}

func (h *mockHandler) ProcessError(err error) {
	if h.expectedError == "" {
		assert.NoError(h.t, err)
	} else {
		assert.EqualError(h.t, err, h.expectedError)
	}
}

func safeReceive(ch chan struct{}) {
	select {
	case <-time.After(1 * time.Minute):
		panic("nothing received")
	case <-ch:
	}
}

func newConfig(port string) (c ClientConfig) {
	c.CleanSession = true
	c.Address = "tcp://localhost:" + port
	defaults.Set(&c)
	return
}

func fakeBroker(t *testing.T, testFlows ...*flow.Flow) (chan struct{}, string) {
	done := make(chan struct{})

	server, err := transport.Launch("tcp://localhost:0")
	assert.NoError(t, err)

	go func() {
		for _, flow := range testFlows {
			conn, err := server.Accept()
			assert.NoError(t, err)

			err = flow.Test(conn)
			assert.NoError(t, err)
		}

		err = server.Close()
		assert.NoError(t, err)

		close(done)
	}()

	_, port, _ := net.SplitHostPort(server.Addr().String())
	return done, port
}

func connectPacket() *packet.Connect {
	pkt := packet.NewConnect()
	pkt.CleanSession = true
	return pkt
}

func connackPacket() *packet.Connack {
	pkt := packet.NewConnack()
	pkt.ReturnCode = packet.ConnectionAccepted
	pkt.SessionPresent = false
	return pkt
}

func disconnectPacket() *packet.Disconnect {
	return packet.NewDisconnect()
}
