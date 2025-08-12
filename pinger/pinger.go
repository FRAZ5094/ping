package pinger

import (
	"math/rand"
	"time"

	"github.com/FRAZ5094/ping/config"
)

type PingResult struct {
	Success  bool
	Duration time.Duration
}

type Pinger interface {
	Ping(host config.Host) PingResult
}

type mockPinger struct{}

func (p *mockPinger) Ping(host config.Host) PingResult {
	return PingResult{
		Success:  rand.Intn(2) == 0,
		Duration: 100 * time.Millisecond * time.Duration(rand.Intn(10)),
	}
}

func New() Pinger {
	return &mockPinger{}
}
