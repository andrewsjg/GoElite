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

			// TODO: Revisit - Not sure any of this is any good. Basically a command should execute, return output
			// and status. Then the TUI should simply print the current game state. Feels like this could be simplified.

			if len(m.gameCmd) > 0 {

				var err error
				// TODO: Think of something to add for Info
				status := "INFO: OK"
				output := ""

				if strings.ToUpper(m.gameCmd) == "INFO" {
					output = SprintState(m.game)

				} else {
					status, output, err = m.executeCommand()

					if err != nil {
						status = "There was an error with the command: " + err.Error()
						output = ""
					}
				}

				if output != "" {
					m.systemViewport.SetContent(output)
				} else {
					m.systemViewport.SetContent(SprintState(m.game))
				}

				m.marketViewport.SetContent(SprintMarket(m.game))
				m.statusBar.SetContent(m.game.PlayerCurrentPlanetName(), "  "+cases.Title(language.English).String(status), "", status)
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
	tuiCmds := "jump,buy,sell,local,info"

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
