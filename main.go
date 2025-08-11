package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type host struct {
	name string
	addr string
}

type model struct {
	hosts   []host
	results []pingResult
}

func (h host) ping(c chan<- pingResult) {
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	c <- pingResult{
		host:     h,
		success:  rand.Intn(2) == 0,
		duration: delay,
	}
}

func (m model) runPings() tea.Cmd {
	d := time.Duration(1 * time.Second)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		var results []pingResult
		c := make(chan pingResult, len(m.hosts))
		for _, h := range m.hosts {
			go h.ping(c)
		}
		for range m.hosts {
			results = append(results, <-c)
		}
		return pingMsg{results: results}
	})
}

func (m model) Init() tea.Cmd {
	return m.runPings()
}

func (m model) View() string {
	s := ""
	for _, r := range m.results {
		s += fmt.Sprintf("%s: %s\n", r.host.name, r.duration)
	}
	s += "\n\n"
	s += fmt.Sprintf("Press q to quit \n")
	return s
}

type pingResult struct {
	host     host
	success  bool
	duration time.Duration
}

type pingMsg struct {
	results []pingResult
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// is called every time something happens including when a key is pressed

	switch msg := msg.(type) {
	case pingMsg:
		m.results = msg.results
		return m, m.runPings()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func main() {
	hosts := []host{
		{
			name: "google",
			addr: "www.google.com",
		},
		{
			name: "cloudflare",
			addr: "1.1.1.1",
		},
	}
	model := model{hosts: hosts}

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
