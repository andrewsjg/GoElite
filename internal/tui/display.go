package tui

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/table"
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

	return gameState

}

// / Returns a formatted string showing the local market
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

// / Returns a formatted string showing system information
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

// Returns a formatted string showing reachable planets
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

// Returns a formatted string of the commander data
func SprintCmdrData(game *eliteEngine.Game) string {
	shipData := "Commander Info\n"
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	//holdSpace := fmt.Sprintf("\n%s %dt\n", nameStyle.Render("Hold Space:"), game.Player.Ship.Holdspace)
	cmdrInfo := fmt.Sprintf("\n%s TODO\n%s TODO", nameStyle.Render("Commander Name:"), nameStyle.Render("Rank:"))
	cash := fmt.Sprintf("\n%s %.1f", nameStyle.Render("Commander Cash:"), float64(game.Player.Cash)/float64(10))

	shipData = headerStyle.Render(shipData) + cash + cmdrInfo

	return shipData

}

// Returns a formatted help string
func SprintHelp() string {
	helpText := ""

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	argStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	buf := new(bytes.Buffer)
	w := new(tabwriter.Writer)
	w.Init(buf, 0, 0, 1, ' ', tabwriter.AlignRight)

	fmt.Fprintln(w, titleStyle.Render("Text Elite Help\n"))
	fmt.Fprintln(w, "Commands are:")
	fmt.Fprintln(w, nameStyle.Render(" Buy")+"\t "+argStyle.Render("<commodity>")+"\t "+argStyle.Render("<ammount>")+"\t - Buy <ammount> of a commodity")
	fmt.Fprintln(w, nameStyle.Render(" Buy")+"\t "+argStyle.Render("fuel")+"\t "+argStyle.Render("<ammount>"))
	fmt.Fprintln(w, nameStyle.Render(" Sell")+"\t "+argStyle.Render("<commodity>")+"\t "+argStyle.Render("<ammount>"))

	w.Flush()

	helpText = buf.String()

	return helpText
}

func HoldTable(game *eliteEngine.Game) table.Model {

	columns := []table.Column{}
	//rows := []table.Row{}
	quants := []string{}

	for _, commodity := range game.Commodities {

		col := table.Column{Title: commodity.Abbrievation, Width: len(commodity.Abbrievation)}
		columns = append(columns, col)
	}

	for _, quantity := range game.Player.Ship.Hold {

		quant := strconv.Itoa(int(quantity))
		quants = append(quants, quant)
	}

	rows := []table.Row{quants}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		//table.WithFocused(true),
		table.WithHeight(2),
	)

	return t
}
