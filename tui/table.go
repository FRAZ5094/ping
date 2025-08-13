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

	UpStatus   = lipgloss.NewStyle().Foreground(successColor).SetString("UP")
	DownStatus = lipgloss.NewStyle().Foreground(errorColor).SetString("DOWN")

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

func RenderTable(results map[string]*pinger.PingResult, hosts []config.Host, spinnerModel spinner.Model) string {
	headers := []string{"HOST", "ADDRESS", "STATUS", "LATENCY"}

	data := [][]string{}
	for _, host := range hosts {
		result := results[host.Name]
		if result == nil {
			spinnerModel.Spinner = spinner.Dot
			data = append(data, []string{host.Name, host.Addr, spinnerModel.View(), "N/A"})
		} else {
			if result.Success && result.Latency != nil {
				latencyStyle := getLatencyStyle(*result.Latency)
				data = append(data, []string{host.Name, host.Addr, UpStatus.Render(), latencyStyle.Render(result.Latency.String())})
			} else {
				data = append(data, []string{host.Name, host.Addr, DownStatus.Render(), LatencyStyleBad.Render("N/A")})
			}
		}
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
