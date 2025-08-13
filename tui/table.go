package tui

import (
	"ping/config"
	"ping/pinger"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	gray         = lipgloss.Color("#c0caf5")
	errorColor   = lipgloss.Color("#f7768e")
	warningColor = lipgloss.Color("#ff9e64")
	successColor = lipgloss.Color("#9ece6a")
	PrimaryColor = lipgloss.Color("#7dcfff")

	BorderStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	HeaderStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Align(lipgloss.Center)

	rowStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(gray).AlignHorizontal(lipgloss.Left)

	UpStatus   = lipgloss.NewStyle().Foreground(successColor)
	DownStatus = lipgloss.NewStyle().Foreground(errorColor)

	LatencyStyleGood    = lipgloss.NewStyle().Foreground(successColor)
	LatencyStyleWarning = lipgloss.NewStyle().Foreground(warningColor)
	LatencyStyleBad     = lipgloss.NewStyle().Foreground(errorColor)
)

func getLatencyStyle(latency time.Duration) lipgloss.Style {
	if latency < 50*time.Millisecond {
		return LatencyStyleGood
	}
	if latency < 100*time.Millisecond {
		return LatencyStyleWarning
	}
	return LatencyStyleBad
}

func RenderRow(host string, addr string, status string, latency string) []string {
	return []string{host, addr, status, latency}
}

func CreateRowFromResult(result *pinger.PingResult, host config.Host, spinnerModel spinner.Model) []string {

	// If no result then its loading
	if result == nil {
		spinnerModel.Spinner = spinner.Dot
		return RenderRow(host.Name, host.Addr, spinnerModel.View(), "N/A")
	}

	// If result was a success render an UP status row
	if result.Success && result.Latency != nil {
		latency := *result.Latency
		latencyStyle := getLatencyStyle(latency)
		return RenderRow(host.Name, host.Addr, UpStatus.Render("UP"), latencyStyle.Render(latency.String()))
	}

	// Then the host is down, so render a DOWN status row
	return RenderRow(host.Name, host.Addr, DownStatus.Render("DOWN"), LatencyStyleBad.Render("N/A"))
}

func RenderTable(results map[string]*pinger.PingResult, hosts []config.Host, spinnerModel spinner.Model) string {
	headers := []string{"HOST", "ADDRESS", "STATUS", "LATENCY"}

	data := [][]string{}
	for _, host := range hosts {
		result := results[host.Name]
		row := CreateRowFromResult(result, host, spinnerModel)
		data = append(data, row)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(BorderStyle).
		Headers(headers...).
		Width(80).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderStyle
			}

			return rowStyle
		})

	return t.Render()
}
