package tui

import (
	"time"

	"github.com/FRAZ5094/ping/config"
	"github.com/FRAZ5094/ping/pinger"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	gray        = lipgloss.Color("245")
	BorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7dcfff"))
	HeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7dcfff")).Bold(true).Align(lipgloss.Center)

	rowStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(gray)).AlignHorizontal(lipgloss.Left)

	UpStatus   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("UP")
	DownStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).SetString("DOWN")

	LatencyStyleGood = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	LatencyStyleBad  = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
)

func getLatencyStyle(latency time.Duration) lipgloss.Style {
	if latency < 50*time.Millisecond {
		return LatencyStyleGood
	}
	return LatencyStyleBad
}

func RenderTable(results map[string]*pinger.PingResult, hosts []config.Host, spinnerModel spinner.Model) string {
	headers := []string{"HOST", "STATUS", "LATENCY"}

	data := [][]string{}
	for _, host := range hosts {
		result := results[host.Name]
		if result == nil {
			spinnerModel.Spinner = spinner.Dot
			data = append(data, []string{host.Name, spinnerModel.View(), "N/A"})
		} else {
			if result.Success {
				latencyStyle := getLatencyStyle(result.Latency)
				data = append(data, []string{host.Name, UpStatus.Render(), latencyStyle.Render(result.Latency.String())})
			} else {
				data = append(data, []string{host.Name, DownStatus.Render(), LatencyStyleBad.Render("N/A")})
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
			// if row == table.HeaderRow {
			// 	return headerStyle
			// }

			// if data[row][1] == "Pikachu" {
			// 	return selectedStyle
			// }

			// even := row%2 == 0

			// switch col {
			// case 2, 3: // Type 1 + 2
			// 	c := typeColors
			// 	if even {
			// 		c = dimTypeColors
			// 	}

			// 	color := c[fmt.Sprint(data[row][col])]
			// 	return baseStyle.Foreground(color)
			// }

			// if even {
			// 	return baseStyle.Foreground(lipgloss.Color("245"))
			// }
			// return baseStyle.Foreground(lipgloss.Color("252"))
		})

	return t.Render()
}
