package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Tui) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	//TODO: Look at how to make these more dynamic and cater for light/dark terminals
	sysBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(80).Height(35)

	mktBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(60).Height(35)

	cmdBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(142)

	leftViewPort := sysBorder.Render(m.systemViewport.View())
	rightViewPort := mktBorder.Render(m.marketViewport.View())
	statusBar := m.statusBar.View()

	commandInput := cmdBorder.Render(fmt.Sprintf("Command %s\n\n", m.cmdInput.View()))

	viewPorts := lipgloss.JoinHorizontal(lipgloss.Top, leftViewPort, rightViewPort)
	composedView := lipgloss.JoinVertical(lipgloss.Top, viewPorts, commandInput, statusBar)

	return fmt.Sprintf(
		style.Render("--== Elite v1.5 ==--")+"\n\n%s",
		composedView,
	) + "\n\n"
}
