package tui

import (
	"log"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
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
	cmdPrompt.Focus()

	//TODO: Fix these sizes
	sysvp := viewport.New(100, 30)
	mktvp := viewport.New(100, 30)

	sysvp.SetContent(game.SprintState())
	mktvp.SetContent(game.Galaxy.Systems[game.Player.Ship.Location.CurrentPlanet].SprintMarket(game.Commodities))

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

	sb.SetContent(game.PlayerCurrentPlanetName(), "", "", "INFO: OK")
	//TODO: Fix these sizes
	sb.SetSize(144)

	return &Tui{
		game:           game,
		cmdInput:       cmdPrompt,
		systemViewport: sysvp,
		marketViewport: mktvp,
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
