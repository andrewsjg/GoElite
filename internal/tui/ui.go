package tui

import (
	"log"
	"strings"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

// Bubble Tea UI

// TODO: Rename the model. This doesnt make sense any more
type CommandModel struct {
	systemViewport viewport.Model
	marketViewport viewport.Model
	cmdInput       textinput.Model
	statusBar      statusbar.Bubble
	game           eliteEngine.Game
	gameCmd        string
}

func (m CommandModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		tiCmd    tea.Cmd
		sysVpCmd tea.Cmd
		mktVpCmd tea.Cmd
	)

	m.cmdInput, tiCmd = m.cmdInput.Update(msg)
	m.systemViewport, sysVpCmd = m.systemViewport.Update(msg)
	m.marketViewport, mktVpCmd = m.marketViewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			//style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

			m.gameCmd = m.cmdInput.Value()

			if strings.ToUpper(m.gameCmd) == "EXIT" {
				return m, tea.Quit
			}

			if strings.ToUpper(m.gameCmd) == "Q" {
				return m, tea.Quit
			}

			if len(m.gameCmd) > 0 {

				status, output, err := m.executeCommand()

				if err != nil {
					output = "There was an error with the command: " + err.Error()
				}

				// Placeholder until the status bar is implemented
				output = status + "\n\n" + output
				m.systemViewport.SetContent(output)

			}

			m.cmdInput.Reset()
			m.cmdInput.Blink()
			m.systemViewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, sysVpCmd, mktVpCmd)
}

var _ tea.Model = &CommandModel{}

// Execute a game command
// Looks up the game command in the map of commands, if it finds a match it calls the
// function stored in the map value
func (m *CommandModel) executeCommand() (status string, output string, err error) {
	cmdOutput := "Command not found"
	cmdParts := strings.Fields(m.gameCmd)
	cmds := m.game.GameCommands

	if cmds[cmdParts[0]] != nil {

		// Pull the command function out of the commands map
		cmdFunc := cmds[cmdParts[0]]

		// Call the command function
		status, cmdOutput = cmdFunc(&m.game, cmdParts)
	}

	return status, cmdOutput, nil
}

func NewCommand(game eliteEngine.Game) *CommandModel {

	cmdPrompt := textinput.New()
	cmdPrompt.Focus()

	//TODO: Fix these sizes
	sysvp := viewport.New(100, 25)
	mktvp := viewport.New(100, 25)

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
			Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)
	sb.SetContent("test.txt", "~/.config/nvim", "1/23", "SB")
	//TODO: Fix these sizes
	sb.SetSize(112)

	return &CommandModel{
		game:           game,
		cmdInput:       cmdPrompt,
		systemViewport: sysvp,
		marketViewport: mktvp,
		statusBar:      sb,
	}
}

func Start() error {

	game := eliteEngine.InitGame(false)

	p := tea.NewProgram(NewCommand(game))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	return nil
}
