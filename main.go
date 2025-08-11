package main

import (
	"fmt"
	"ping/config"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func ping(host config.Host, settings config.Settings, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := host.Addr
	fmt.Println("Pinging:", addr)

	pinger, err := probing.NewPinger(host.Addr)
	if err != nil {
		fmt.Println(addr, ": error resolving addr ", err)
		return
	}
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Duration(settings.TimeoutSeconds) * time.Second
	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		fmt.Println(addr, ": error pinging addr", addr, ": ", err)
		return
	}
	stats := pinger.Statistics()

	fmt.Println(addr, ":", stats.AvgRtt.Milliseconds(), "ms")
}

func main() {
	config := config.Parse("config.yaml")
	wg := new(sync.WaitGroup)

	for _, host := range config.Hosts {
		wg.Add(1)
		go ping(host, config.Settings, wg)
	}
	wg.Wait()
	fmt.Scanf("h")
}
