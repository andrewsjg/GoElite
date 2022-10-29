package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func (m Tui) View() string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).MarginBottom(0)

	//TODO: Look at how to make these more dynamic and cater for light/dark terminals
	sysBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(80).Height(26)

	mktBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(60).Height(26)

	cmdrBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(60).Height(5)

	guageBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).Width(80).Height(5)

	cmdBorder := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(142)

	leftViewPort := sysBorder.Render(m.systemViewport.View())
	rightViewPort := mktBorder.Render(m.marketViewport.View())

	statusBar := m.statusBar.View()

	fuelValue := ((float64(m.game.Player.Ship.Fuel) / float64(70)) * float64(100)) / 100
	fuelGauge := m.fuelBar.ViewAs(fuelValue)
	fuelTitle := nameStyle.Render("\nFuel") + fmt.Sprintf(" (%.1fLY):     ", float64(m.game.Player.Ship.Fuel)/float64(10))

	holdSpaceValue := ((float64(m.game.Player.Ship.Holdspace) / float64(20)) * float64(100)) / 100
	holdSpaceGuage := m.holdSpaceBar.ViewAs(holdSpaceValue)
	holdTitle := nameStyle.Render("\nHold Space") + fmt.Sprintf(" (%dt): ", m.game.Player.Ship.Holdspace)

	cmdrInfo := m.cmdrViewport.View()
	cmdrInfo = cmdrBorder.Render(cmdrInfo)

	guageViews := titleStyle.Render("Ship Info\n") + fuelTitle + fuelGauge + " full\n" + holdTitle + holdSpaceGuage + " available\n\n" + nameStyle.Render("Hold Contents::") + "\n" + m.holdTable.View()
	guageViews = guageBorder.Render(guageViews)

	// Ship + Cmdr views
	dataViews := lipgloss.JoinHorizontal(lipgloss.Top, guageViews, cmdrInfo)

	commandInput := cmdBorder.Render(fmt.Sprintf("Command %s\n\n", m.cmdInput.View()))

	viewPorts := lipgloss.JoinHorizontal(lipgloss.Top, leftViewPort, rightViewPort)
	composedView := lipgloss.JoinVertical(lipgloss.Top, viewPorts, dataViews, commandInput, statusBar)
	figTitle := figure.NewFigure("         --== Elite v1.5 ==--", "small", true)
	title := figTitle.String()

	return fmt.Sprintf(

		//titleStyle.Render("--== Elite v1.5 ==--")+"\n\n%s",
		titleStyle.Render(title)+"\n%s",
		composedView,
	) + "\n\n"
}
