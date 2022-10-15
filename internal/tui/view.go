package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m CommandModel) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	border := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	leftViewPort := border.Render(m.systemViewport.View())
	rightViewPort := border.Render(m.marketViewport.View())
	statusBar := m.statusBar.View()

	commandInput := fmt.Sprintf("Command %s\n\n", m.cmdInput.View())

	viewPorts := lipgloss.JoinHorizontal(lipgloss.Top, leftViewPort, rightViewPort)
	composedView := lipgloss.JoinVertical(lipgloss.Top, viewPorts, statusBar, commandInput)

	return fmt.Sprintf(
		style.Render("--== Elite v1.5 ==--")+"\n\n%s",
		composedView,
	) + "\n\n"
}
