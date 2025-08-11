package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/FRAZ5094/ping/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type host struct {
	name string
	addr string
}

type model struct {
	hosts   []host
	results []pingResult
	spinner spinner.Model
}

func (h host) ping(c chan<- pingResult, ordering int) {
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	c <- NewPingResult(h, rand.Intn(2) == 0, delay, ordering)
}

func (m model) runPings() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		var results []pingResult
		c := make(chan pingResult, len(m.hosts))
		for i, h := range m.hosts {
			go h.ping(c, i)
		}
		for range m.hosts {
			results = append(results, <-c)
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].ordering < results[j].ordering
		})
		return pingMsg{results: results}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.runPings(),
		m.spinner.Tick,
	)
}

func (m model) View() string {
	s := "\n"
	m.spinner.Spinner = spinner.Line
	s += m.spinner.View() + " Running pings...\n\n"
	if len(m.results) == 0 {
		m.spinner.Spinner = spinner.Dot
		for _, host := range m.hosts {
			s += fmt.Sprintf("%s%s: ...\n", m.spinner.View(), styles.NameStyle.Render(host.name))
		}
	} else {
		for _, r := range m.results {
			if r.success {
				s += fmt.Sprintf("%s %s: %s\n", styles.CheckMark, styles.NameStyle.Render(r.host.name), r.duration)
			} else {
				s += fmt.Sprintf("%s %s: %s\n", styles.CrossMark, styles.NameStyle.Render(r.host.name), r.duration)
			}
		}
	}
	s += "\n\n"
	s += "Press any key to exit \n"
	return s
}

type pingResult struct {
	host     host
	success  bool
	duration time.Duration
	ordering int
}

func NewPingResult(host host, success bool, duration time.Duration, ordering int) pingResult {
	return pingResult{
		host:     host,
		success:  success,
		duration: duration,
		ordering: ordering,
	}
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
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func main() {
	hosts := []host{
		{
			name: "Google",
			addr: "www.google.com",
		},
		{
			name: "Cloudflare",
			addr: "1.1.1.1",
		},
	}
	var spinner = spinner.New()
	spinner.Style = styles.SpinnerStyle
	model := model{hosts: hosts, spinner: spinner}

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
