package pinger

import (
	"math/rand"
	"ping/config"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingResult struct {
	Success bool
	Latency *time.Duration
}

func newPingResult(success bool, latency *time.Duration) PingResult {
	return PingResult{
		Success: success,
		Latency: latency,
	}
}

type Pinger interface {
	Ping(addr string) PingResult
}

type mockPinger struct{}

func (p *mockPinger) Ping(host config.Host) PingResult {
	latency := time.Millisecond * time.Duration(rand.Intn(100))
	success := rand.Intn(2) == 0
	time.Sleep(latency)
	return newPingResult(success, &latency)
}

type livePinger struct{}

func (p *livePinger) Ping(addr string) PingResult {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		return newPingResult(
			false,
			nil,
		)
	}
	pinger.Count = 1
	pinger.SetPrivileged(true)
	pinger.Timeout = 1 * time.Second
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return newPingResult(false, nil)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	return newPingResult(true, &stats.AvgRtt)
}

func New() Pinger {
	return &livePinger{}
}
