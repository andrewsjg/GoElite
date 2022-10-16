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

				status := ""
				output := ""

				if strings.ToUpper(m.gameCmd) == "INFO" {
					output = SprintState(m.game)

				} else if strings.ToUpper(m.gameCmd) == "LOCAL" {
					output = SprintLocal(m.game)

				} else {
					status, output = m.executeTUICommand(m.gameCmd)
					//status, output = m.game.ExecuteCommand(m.gameCmd)
				}

				if output != "" {
					m.systemViewport.SetContent(output)
				} else {
					m.systemViewport.SetContent(SprintState(m.game))
				}

				m.marketViewport.SetContent(SprintMarket(m.game))
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

// filter commands to use only the commands required by the TUI.
// TODO: Might be a better way to do this?
func (m *Tui) executeTUICommand(cmd string) (status string, output string) {

	// Need to filter the commands that are valid in the tui
	tuiCmds := "jump,buy,sell,local,info"
	cmdParts := strings.Fields(cmd)

	if len(cmdParts) >= 1 {
		if strings.Contains(tuiCmds, cmdParts[0]) {
			// Call the command function
			status, output = m.game.ExecuteCommand(cmd)

		}
	}
	return status, output
}
