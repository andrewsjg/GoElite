package tui

import (
	"log"
	"strconv"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

// Bubble Tea UI

func (m Tui) Init() tea.Cmd {
	return textinput.Blink
}

func New(game eliteEngine.Game) *Tui {

	cmdPrompt := textinput.New()

	//TODO: Fix these sizes
	sysvp := viewport.New(100, 24)
	mktvp := viewport.New(100, 24)
	cmdrvp := viewport.New(140, 5)

	fuelGuage := progress.New()
	fuelGuage.Width = 50

	holdSpace := progress.New()
	holdSpace.Width = 50

	holdTable := HoldTable(&game)

	sysvp.SetContent(SprintState(&game))
	mktvp.SetContent(SprintMarket(&game))
	cmdrvp.SetContent(SprintCmdrData(&game))

	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)

	system := game.GetPlanetaryData(game.PlayerCurrentPlanetName())
	pos := "(" + strconv.Itoa(int(system.X)) + "," + strconv.Itoa(int(system.Y)) + ")"

	sb.SetContent(game.PlayerCurrentPlanetName(), "", "", "POS: "+pos)
	//TODO: Fix these sizes
	sb.SetSize(144)

	cmdPrompt.Focus()

	return &Tui{
		game:           game,
		cmdInput:       cmdPrompt,
		systemViewport: sysvp,
		marketViewport: mktvp,
		cmdrViewport:   cmdrvp,
		fuelBar:        fuelGuage,
		holdSpaceBar:   holdSpace,
		holdTable:      holdTable,
		statusBar:      sb,
	}
}

func Start() error {

	game := eliteEngine.InitGame(false)

	p := tea.NewProgram(New(game))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	return nil
}
