package ingress

import (
	"io"
	"log"
	"sync/atomic"
	"time"
	"unsafe"

	v2 "github.com/cloudfoundry-incubator/scalable-syslog/internal/api/loggregator/v2"
)

type ConnectionBuilder interface {
	Connect() (io.Closer, v2.EgressClient, error)
}

type connection struct {
	closer io.Closer
	client v2.EgressClient
}

// ClientManager manages loggregator egress clients and connections.
type ClientManager struct {
	connector     ConnectionBuilder
	connectionTTL time.Duration
	connections   []unsafe.Pointer
	nextIdx       uint64
	retryWait     time.Duration
}

type ClientManagerOpts func(*ClientManager)

func WithRetryWait(d time.Duration) func(*ClientManager) {
	return func(c *ClientManager) {
		c.retryWait = d
	}
}

// NewClientManager returns a ClientManager after opening the specified number
// of connections.
func NewClientManager(connector ConnectionBuilder, connCount int, ttl time.Duration, opts ...ClientManagerOpts) *ClientManager {
	c := &ClientManager{
		connector:     connector,
		connectionTTL: ttl,
		connections:   make([]unsafe.Pointer, connCount),
		retryWait:     2 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	for i := 0; i < connCount; i++ {
		c.openNewConnection(i)
	}

	go c.monitorConnectionsForRolling()

	return c
}

// Next returns the next available loggregator egress client. Next will block
// until a healthy client is available.
func (c *ClientManager) Next() v2.EgressClient {
	for {
		idx := atomic.AddUint64(&c.nextIdx, 1)
		conn := (*connection)(atomic.LoadPointer(&c.connections[idx%uint64(len(c.connections))]))
		if conn.client != nil {
			return conn.client
		}

		time.Sleep(c.retryWait)
	}
}

func (c *ClientManager) monitorConnectionsForRolling() {
	for range time.Tick(c.connectionTTL) {
		for i := 0; i < len(c.connections); i++ {
			conn := (*connection)(atomic.LoadPointer(&c.connections[i]))

			if conn.closer != nil {
				conn.closer.Close()
			}

			c.openNewConnection(i)
		}
	}
}

func (c *ClientManager) openNewConnection(idx int) {
	closer, client, err := c.connector.Connect()
	if err != nil {
		log.Printf("Failed to connect to loggregator API: %s", err)

		var nilConn *connection
		atomic.SwapPointer(&c.connections[idx], unsafe.Pointer(&nilConn))

		return
	}

	atomic.SwapPointer(&c.connections[idx], unsafe.Pointer(&connection{
		closer: closer,
		client: client,
	}))
}
