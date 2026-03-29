package pinger

import (
	"math/rand"
	"runtime"
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

//nolint:unused
type mockPinger struct{}

//nolint:unused
func (p mockPinger) Ping(addr string) PingResult {
	latency := time.Millisecond * time.Duration(rand.Intn(100))
	success := rand.Intn(2) == 0
	time.Sleep(latency)
	return newPingResult(success, &latency)
}

//nolint:unused
type livePinger struct{}

//nolint:unused
func (p livePinger) Ping(addr string) PingResult {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		return newPingResult(
			false,
			nil,
		)
	}
	pinger.Count = 1
	pinger.SetPrivileged(runtime.GOOS == "windows")
	pinger.Timeout = 1 * time.Second
	err = pinger.Run()
	if err != nil {
		return newPingResult(false, nil)
	}
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return newPingResult(false, nil)
	}

	return newPingResult(true, &stats.AvgRtt)
}

func New() Pinger {
	return livePinger{}
	// return mockPinger{}
}
