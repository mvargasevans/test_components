package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

// Envelope is the wire format for all IPC messages.
type Envelope struct {
	Interaction string          `json:"interaction"`
	OK          bool            `json:"ok,omitempty"`
	Error       string          `json:"error,omitempty"`
	Body        json.RawMessage `json:"body"`
}

// Conn wraps a net.Conn with line-delimited JSON encoding/decoding.
type Conn struct {
	conn    net.Conn
	enc     *json.Encoder
	scanner *bufio.Scanner
}

func newConn(c net.Conn) *Conn {
	return &Conn{
		conn:    c,
		enc:     json.NewEncoder(c),
		scanner: bufio.NewScanner(c),
	}
}

// Send encodes and sends an Envelope as a newline-delimited JSON message.
func (c *Conn) Send(env Envelope) error {
	return c.enc.Encode(env)
}

// Receive reads the next newline-delimited JSON message into an Envelope.
func (c *Conn) Receive() (Envelope, error) {
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return Envelope{}, err
		}
		return Envelope{}, fmt.Errorf("connection closed")
	}
	var env Envelope
	if err := json.Unmarshal(c.scanner.Bytes(), &env); err != nil {
		return Envelope{}, err
	}
	return env, nil
}

// Close closes the underlying connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Dial connects to a Unix domain socket at path and returns a Conn.
func Dial(path string) (*Conn, error) {
	c, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return newConn(c), nil
}

// Listen binds a Unix domain socket at path and returns a net.Listener.
func Listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

// Accept accepts one connection from l and returns a Conn.
func Accept(l net.Listener) (*Conn, error) {
	c, err := l.Accept()
	if err != nil {
		return nil, err
	}
	return newConn(c), nil
}

// TempSocket returns a unique temp socket path and registers cleanup with t.
func TempSocket(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "ipc-test-*")
	if err != nil {
		t.Fatalf("TempSocket: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return filepath.Join(dir, "test.sock")
}
