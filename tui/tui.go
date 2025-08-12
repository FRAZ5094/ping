package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/FRAZ5094/ping/config"
	"github.com/FRAZ5094/ping/pinger"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	hosts          []config.Host
	results        map[string]*pinger.PingResult
	resultsChan    chan pingResultMsg
	spinner        spinner.Model
	timer          timer.Model
	timerResetChan chan struct{}
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

var pingInterval = 5 * time.Second

func runPings(hosts []config.Host, resultsChan chan pingResultMsg, timerResetChan chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			for _, host := range hosts {
				go func() {
					result := ping(host)
					resultMsg := pingResultMsg{host: host, result: result}
					resultsChan <- resultMsg
				}()
			}
			time.Sleep(pingInterval)
			timerResetChan <- struct{}{}
		}
	}
}

func waitForResult(resultsChan chan pingResultMsg) tea.Cmd {
	return func() tea.Msg {
		result := <-resultsChan
		return result
	}
}

type timerResetMsg struct{}

func waitForTimerReset(timerResetChan chan struct{}) tea.Cmd {
	return func() tea.Msg {
		resetEvent := <-timerResetChan
		return timerResetMsg(resetEvent)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		runPings(m.hosts, m.resultsChan, m.timerResetChan),
		waitForResult(m.resultsChan),
		m.spinner.Tick,
		m.timer.Init(),
		waitForTimerReset(m.timerResetChan),
	)
}

func (m model) View() string {
	s := "\n"

	// s += m.spinner.View() + " Running pings...\n\n"

	s += RenderTable(m.results, m.hosts, m.spinner)
	s += "\n"
	m.spinner.Spinner = spinner.Dot
	s += fmt.Sprintf("Next ping in: %s\n", m.timer.View())
	s += "\n"
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
	case timerResetMsg:
		m.timer = timer.NewWithInterval(pingInterval, time.Millisecond)
		return m, tea.Batch(
			m.timer.Init(),
			waitForTimerReset(m.timerResetChan),
		)
	default:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		var timerCmd tea.Cmd
		m.timer, timerCmd = m.timer.Update(msg)
		return m, tea.Batch(spinnerCmd, timerCmd)
	}
}

func Start(hosts []config.Host) {

	var spinner = spinner.New()
	var SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	spinner.Style = SpinnerStyle

	resultsChan := make(chan pingResultMsg)
	results := make(map[string]*pinger.PingResult)

	for _, host := range hosts {
		results[host.Name] = nil
	}

	model := model{hosts: hosts, spinner: spinner, resultsChan: resultsChan, results: results, timer: timer.NewWithInterval(pingInterval, time.Millisecond), timerResetChan: make(chan struct{})}

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
