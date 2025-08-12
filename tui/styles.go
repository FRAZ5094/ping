package tui

import "github.com/charmbracelet/lipgloss"

var NameStyle = lipgloss.NewStyle().Bold(true).Blink(true)

var CheckMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
var CrossMark = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).SetString("✗")

var SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
