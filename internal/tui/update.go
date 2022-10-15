package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
					status = "There was an error with the command: " + err.Error()
					output = ""
				}

				if output != "" {
					m.systemViewport.SetContent(output)
				} else {
					m.systemViewport.SetContent(m.game.SprintState())
				}

				m.marketViewport.SetContent(m.game.Galaxy.Systems[m.game.Player.Ship.Location.CurrentPlanet].SprintMarket(m.game.Commodities))

				// TODO: Think of something to add for Info
				m.statusBar.SetContent(m.game.PlayerCurrentPlanetName(), "  "+cases.Title(language.English).String(status), "", "INFO: OK")
			}

			m.cmdInput.Reset()
			m.cmdInput.Blink()
			m.systemViewport.GotoBottom()

		}
	}

	return m, tea.Batch(tiCmd, sysVpCmd, mktVpCmd)
}

var _ tea.Model = &Tui{}

// Execute a game command
// Looks up the game command in the map of commands, if it finds a match it calls the
// function stored in the map value
func (m *Tui) executeCommand() (status string, output string, err error) {
	cmdOutput := ""
	status = "Command not found"
	cmdParts := strings.Fields(m.gameCmd)
	cmds := m.game.GameCommands

	// Need to filter the commands that are valid in the tui
	tuiCmds := "jump,buy,local,info"

	if cmds[cmdParts[0]] != nil {

		// Pull the command function out of the commands map
		cmdFunc := cmds[cmdParts[0]]

		if strings.Contains(tuiCmds, cmdParts[0]) {
			// Call the command function
			status, cmdOutput = cmdFunc(&m.game, cmdParts)
		}
	}

	return status, cmdOutput, nil
}
