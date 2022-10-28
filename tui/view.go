package tui

import (
	"pixel/tui/constants"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	focusedColor := lipgloss.Color("201")
	borderColor := lipgloss.Color("69")
	errString := ""
	m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(borderColor).Width(m.viewport.Width)

	// channel pane
	listPane := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(borderColor).Height(m.viewport.Height).Width(m.viewport.Width / 7).Padding(1)
	inputPane := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(1).Width(m.viewport.Width).BorderForeground(borderColor).Padding(1)
	switch m.Focus {
	case Input:
		inputPane.BorderForeground(focusedColor)
	case Feed:
		m.viewport.Style.BorderForeground(focusedColor)
	default:
		// list is default
		listPane.BorderForeground(focusedColor)

	}

	// chat window and input
	rightPane := lipgloss.JoinVertical(lipgloss.Center, m.viewport.View(), inputPane.Render(m.textarea.View()))
	if m.err != nil {
		errString = m.err.Error()
	}
	formatted := lipgloss.JoinHorizontal(lipgloss.Left, listPane.Render(m.list.View()), rightPane, errString)

	return constants.DocStyle.Render(formatted)
}
