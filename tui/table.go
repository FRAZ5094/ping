package tui

import (
	"os"

	"github.com/FRAZ5094/ping/config"
	"github.com/FRAZ5094/ping/pinger"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func RenderTable(results map[string]*pinger.PingResult, hosts []config.Host) string {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	// headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	headers := []string{"Host", "Status", "Duration"}

	data := [][]string{}
	for _, host := range hosts {
		result := results[host.Name]
		if result == nil {
			data = append(data, []string{host.Name, "Loading...", "..."})
		} else {
			if result.Success {
				data = append(data, []string{host.Name, "OK", result.Duration.String()})
			} else {
				data = append(data, []string{host.Name, "Error", result.Duration.String()})
			}
		}
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(80).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return baseStyle
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
