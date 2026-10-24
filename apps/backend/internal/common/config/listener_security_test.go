package config

import (
	"net"
	"testing"
)

// @covers AC-PLATFORM-LOCAL-LISTENER-001.1
func TestDefaultServerHostIsLoopback(t *testing.T) {
	binds, err := (&ServerConfig{}).ResolvedBinds()
	if err != nil || len(binds) != 1 {
		t.Fatalf("default binds = %v, %v", binds, err)
	}
	listener, err := net.Listen("tcp4", net.JoinHostPort(binds[0], "0"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	if !listener.Addr().(*net.TCPAddr).IP.IsLoopback() {
		t.Fatalf("default listener accepts non-loopback traffic: %s", listener.Addr())
	}
}
