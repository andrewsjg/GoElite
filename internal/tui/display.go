package tui

import (
	"fmt"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/lipgloss"
)

// Functions for displaying formatted game information

// Return game state as a string
func SprintState(g *eliteEngine.Game) string {

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	//style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	gameState := ""

	gal := g.Galaxy
	shipLocation := g.Player.Ship.Location
	planet := gal.Systems[shipLocation.CurrentPlanet].Name

	gameState = fmt.Sprintf("%s\n\n", headerStyle.Render("System Info"))
	gameState = gameState + SprintSystem(g, planet, false)

	/*gameState = gameState + fmt.Sprintf("\n%s %.1f\n", style.Render("Cash:"), float64(g.Player.Cash)/float64(10))
	gameState = gameState + fmt.Sprintf("%s %.1f\n", style.Render("Fuel:"), float64(g.Player.Ship.Fuel)/float64(10))
	gameState = gameState + fmt.Sprintf("%s %dt", style.Render("Hold Space:"), g.Player.Ship.Holdspace) */

	return gameState

}

func SprintMarket(g *eliteEngine.Game) string {
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	fieldNameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	colNameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	marketData := ""
	numCommodities := len(g.Commodities) - 1
	gal := g.Galaxy
	p := gal.Systems[g.Player.Ship.Location.CurrentPlanet]
	mkt := p.Market

	marketData = fmt.Sprintf("%s\n\n", headerStyle.Render("Local Market"))
	//marketData = marketData + fmt.Sprintf("%-*sPrice  Quantity\n", 30, colNameStyle.Render("Commodity"))
	marketData = marketData + fmt.Sprintf("%-*s %s  %s\n", 30, colNameStyle.Render("Commodity"), colNameStyle.Render("Price"), colNameStyle.Render("Quantity"))
	//marketData = marketData + fmt.Sprintln("------------------------------------")
	for i := 0; i <= numCommodities; i++ {
		marketData = marketData + fmt.Sprintf("%-*s", 30, fieldNameStyle.Render(g.Commodities[i].Name))
		marketData = marketData + fmt.Sprintf(" %-*.1f", 5, float64(mkt.Price[i])/float64(10))
		marketData = marketData + fmt.Sprintf("  %d", mkt.Quantity[i])
		marketData = marketData + fmt.Sprintf(mkt.UnitNames[g.Commodities[1].Units])
		marketData = marketData + fmt.Sprintln("")

	}

	//marketData = marketData + fmt.Sprintln("------------------------------------")

	return marketData
}

func SprintSystem(game *eliteEngine.Game, systemName string, compressed bool) string {

	systemData := ""
	sd := game.GetPlanetaryData(systemName)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	if compressed {
		systemData = fmt.Sprintf("%10s", sd.Name)
		systemData = systemData + fmt.Sprintf(" %s %2d ", style.Render(" TL:"), sd.Techlev)
		systemData = systemData + fmt.Sprintf("%12s", sd.EconomyName)
		systemData = systemData + fmt.Sprintf(" %15s", sd.GovName)
	} else {
		systemData = systemData + fmt.Sprintf("%s %s\n", style.Render("System:"), sd.Name)
		systemData = systemData + fmt.Sprintf("%s (%d,", style.Render("Position:"), sd.X)
		systemData = systemData + fmt.Sprintf("%d)\n", sd.Y)
		systemData = systemData + fmt.Sprintf("%s (%d) ", style.Render("Economy:"), sd.Economy)
		systemData = systemData + fmt.Sprintf("%s\n", sd.EconomyName)
		systemData = systemData + fmt.Sprintf("%s (%d) ", style.Render("Government"), sd.Govtype)
		systemData = systemData + fmt.Sprintf("%s\n", sd.GovName)
		systemData = systemData + fmt.Sprintf("%s %2d\n", style.Render("Tech Level:"), sd.Techlev)
		systemData = systemData + fmt.Sprintf("%s %d\n", style.Render("Turnover:"), (sd.Productivity))
		systemData = systemData + fmt.Sprintf("%s %d\n", style.Render("Radius:"), sd.Radius)
		systemData = systemData + fmt.Sprintf("%s %d Billion\n", style.Render("Population:"), (sd.Population)>>3)

		systemData = systemData + sd.Description + "\n"

	}

	return systemData
}

func SprintLocal(game *eliteEngine.Game) string {
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	localSystems := headerStyle.Render("Local Systems\n\n")

	reachable := game.ReachableSystems()

	for _, navinfo := range reachable {
		if navinfo.ReachableWithCurrentFuel {
			localSystems = localSystems + "\n *"
		} else {
			localSystems = localSystems + "\n -"
		}

		localSystems = localSystems + SprintSystem(game, navinfo.System.Name, true)
		localSystems = localSystems + fmt.Sprintf(" (%.1f LY)", float64(navinfo.Distance)/float64(10))
	}

	return localSystems
}

func SprintShipData(game *eliteEngine.Game) string {
	shipData := "Ship Info\n"
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	holdSpace := fmt.Sprintf("\n%s %dt\n", nameStyle.Render("Hold Space:"), game.Player.Ship.Holdspace)
	cash := fmt.Sprintf("\n%s %.1f", nameStyle.Render("Cash:"), float64(game.Player.Cash)/float64(10))

	shipData = headerStyle.Render(shipData) + cash + holdSpace

	return shipData

}
