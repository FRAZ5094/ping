package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/FRAZ5094/ping/config"
	"github.com/FRAZ5094/ping/pinger"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hosts       []config.Host
	results     map[string]*pinger.PingResult
	resultsChan chan pingResultMsg
	spinner     spinner.Model
}

func ping(host config.Host) pinger.PingResult {
	pinger := pinger.New()
	result := pinger.Ping(host)
	return result
}

type pingResultMsg struct {
	host   config.Host
	result pinger.PingResult
}

func runPings(hosts []config.Host, resultsChan chan pingResultMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			for _, host := range hosts {
				go func() {
					result := ping(host)
					resultMsg := pingResultMsg{host: host, result: result}
					resultsChan <- resultMsg
				}()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func waitForResult(resultsChan chan pingResultMsg) tea.Cmd {
	return func() tea.Msg {
		result := <-resultsChan
		return result
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		runPings(m.hosts, m.resultsChan),
		waitForResult(m.resultsChan),
		m.spinner.Tick,
	)
}

func (m model) View() string {
	s := "\n"
	m.spinner.Spinner = spinner.Line
	s += m.spinner.View() + " Running pings...\n\n"

	s += RenderTable(m.results, m.hosts)

	return s

	for _, host := range m.hosts {
		result := m.results[host.Name]
		if result == nil {
			s += fmt.Sprintf("%s %s: Loading...\n", m.spinner.View(), NameStyle.Render(host.Name))
		} else {
			if result.Success {
				s += fmt.Sprintf("%s %s: %s\n", CheckMark, NameStyle.Render(host.Name), result.Duration)
			} else {
				s += fmt.Sprintf("%s %s: %s\n", CrossMark, NameStyle.Render(host.Name), result.Duration)
			}
		}
	}
	// if len(results) == 0 {
	// 	m.spinner.Spinner = spinner.Dot
	// 	for _, host := range m.hosts {
	// 		s += fmt.Sprintf("%s%s: ...\n", m.spinner.View(), NameStyle.Render(host.name))
	// 	}
	// } else {
	// 	for _, r := range *m.results {
	// 		if r.success {
	// 			s += fmt.Sprintf("%s %s: %s\n", CheckMark, NameStyle.Render(r.host.name), r.duration)
	// 		} else {
	// 			s += fmt.Sprintf("%s %s: %s\n", CrossMark, NameStyle.Render(r.host.name), r.duration)
	// 		}
	// 	}
	// }
	s += "\n\n"
	s += "Press any key to exit \n"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pingResultMsg:
		result := msg.result
		host := msg.host
		m.results[host.Name] = &result
		return m, waitForResult(m.resultsChan)
	case tea.KeyMsg:
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func Start(hosts []config.Host) {

	var spinner = spinner.New()
	spinner.Style = SpinnerStyle

	resultsChan := make(chan pingResultMsg)
	results := make(map[string]*pinger.PingResult)

	for _, host := range hosts {
		results[host.Name] = nil
	}

	model := model{hosts: hosts, spinner: spinner, resultsChan: resultsChan, results: results}

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
