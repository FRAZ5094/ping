package pinger

import (
	"math/rand"
	"time"

	"github.com/FRAZ5094/ping/config"
)

type PingResult struct {
	Success bool
	Latency time.Duration
}

type Pinger interface {
	Ping(host config.Host) PingResult
}

type mockPinger struct{}

func (p *mockPinger) Ping(host config.Host) PingResult {
	latency := time.Millisecond * time.Duration(rand.Intn(100))
	success := rand.Intn(2) == 0
	time.Sleep(latency)
	return PingResult{
		Success: success,
		Latency: latency,
	}
}

func New() Pinger {
	return &mockPinger{}
}
