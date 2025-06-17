package balance

import (
	"log/slog"
	"net"
	"strconv"
	"sync"
)

type Manager struct {
	logger *slog.Logger
	addr   string

	mu      sync.Mutex
	balance int
}

func NewManager(addr string, logger *slog.Logger) *Manager {
	m := &Manager{
		logger: logger,
		addr:   addr,
	}
	go m.keepUpdatingBalance()
	return m
}

func (m *Manager) Balance() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.balance
}

func (m *Manager) Decrease(value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.balance -= value
}

func (m *Manager) keepUpdatingBalance() {
	addr, err := net.ResolveUDPAddr("udp4", m.addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			m.logger.With("err", err.Error()).
				Error("failed to read udp in keepUpdatingBalance")
			return
		}
		payment, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
		if err != nil {
			m.logger.With("err", err.Error()).Error("invalid balance received")
			continue
		}
		m.mu.Lock()
		m.logger.Debug("received payment, updating balance")
		m.balance += int(payment)
		m.mu.Unlock()
	}
}
