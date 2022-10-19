package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Tui) View() string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).MarginBottom(0)

	//TODO: Look at how to make these more dynamic and cater for light/dark terminals
	sysBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(80).Height(24)

	mktBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(60).Height(24)

	shipBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(142).Height(5)

	cmdBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(142)

	leftViewPort := sysBorder.Render(m.systemViewport.View())
	rightViewPort := mktBorder.Render(m.marketViewport.View())

	statusBar := m.statusBar.View()

	fuelValue := ((float64(m.game.Player.Ship.Fuel) / float64(70)) * float64(100)) / 100
	fuelGauge := m.fuelBar.ViewAs(fuelValue)
	fuelTitle := nameStyle.Render("\nFuel") + fmt.Sprintf(" (%.1fLY):\n", float64(m.game.Player.Ship.Fuel)/10)

	bottomViews := m.shipViewport.View() + fuelTitle + fuelGauge
	bottomViews = shipBorder.Render(bottomViews)

	commandInput := cmdBorder.Render(fmt.Sprintf("Command %s\n\n", m.cmdInput.View()))

	viewPorts := lipgloss.JoinHorizontal(lipgloss.Top, leftViewPort, rightViewPort)
	composedView := lipgloss.JoinVertical(lipgloss.Top, viewPorts, bottomViews, commandInput, statusBar)

	return fmt.Sprintf(
		titleStyle.Render("--== Elite v1.5 ==--")+"\n\n%s",
		composedView,
	) + "\n\n"
}
