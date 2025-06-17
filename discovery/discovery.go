package discovery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	broadcastAddress      = "255.255.255.255:8829"
	localBroadcastAddress = ":8829"
)

type Discovery struct {
	name          []byte
	localAddr     string
	advertisement []byte
	logger        *slog.Logger
	mu            sync.Mutex
	peers         map[string]string
}

func New(name, localAddr string, logger *slog.Logger) *Discovery {
	d := &Discovery{
		name:          []byte(name),
		localAddr:     strings.Split(localAddr, ":")[0],
		advertisement: []byte(name + " " + localAddr),
		peers:         make(map[string]string),
		logger:        logger,
	}
	go d.broadcasting()
	go d.collectPeers()
	return d
}

func (c *Discovery) Write(p []byte) (n int, err error) {
	tokens := bytes.SplitN(p, []byte(" "), 2)
	name, amount := tokens[0], tokens[1]
	c.mu.Lock()
	addr, ok := c.peers[string(name)]
	c.mu.Unlock()
	if !ok {
		return 0, errors.New("peer not connected")
	}
	var peer net.Conn
	peer, err = net.Dial("udp4", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to dial, %w", err)
	}
	defer peer.Close()
	return peer.Write(amount)
}

func (c *Discovery) broadcasting() error {
	broadcast, err := net.Dial("udp", broadcastAddress)
	if err != nil {
		return fmt.Errorf("failed to allocate broadcast address, %w", err)
	}
	ticker := time.NewTicker(3 * time.Second)
	defer broadcast.Close()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if _, err = broadcast.Write(c.advertisement); err != nil {
				return fmt.Errorf("failed to write broadcast message, %w", err)
			}
		}
	}
}

func control(_, _ string, c syscall.RawConn) error {
	var operr error
	err := c.Control(func(fd uintptr) {
		operr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if operr != nil {
			return
		}
		operr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if operr != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	return operr
}

func (c *Discovery) collectPeers() error {
	conf := net.ListenConfig{Control: control}
	conn, err := conf.ListenPacket(context.Background(), "udp4", localBroadcastAddress)
	if err != nil {
		return fmt.Errorf("failed to listen broadcast packet, %w", err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		var n int
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			return fmt.Errorf("failed to read from broadcast to collect peers, %w", err)
		}
		if bytes.Contains(buffer[:n], c.name) {
			continue
		}
		tokens := bytes.SplitN(buffer[:n], []byte(" "), 2)
		name, localAddr := tokens[0], tokens[1]
		key := string(name)
		c.mu.Lock()
		_, ok := c.peers[key]
		if !ok {
			c.peers[key] = string(localAddr)
			c.logger.With("peer", string(buffer[:n])).Debug("New peer has been collected")
		}
		c.mu.Unlock()
	}
}
