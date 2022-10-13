package tui

import (
	"fmt"
	"log"

	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rivo/tview"
)

// tview UI experiments
func CreateUI() {
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

}

func StartTView() {
	game := eliteEngine.InitGame(false)

	app := tview.NewApplication()

	commands := []string{"Jump", "Buy Fuel", "Hyperspace Jump", "Buy Commodity", "Sell Commodity", "Show Hold"}
	commandPanel := tview.NewList().ShowSecondaryText(false)
	commandPanel.SetBorder(true).SetTitle("Commands")

	flex := tview.NewFlex().
		AddItem(commandPanel, 0, 1, false)

	for _, command := range commands {
		commandPanel.AddItem(command, "", 0, func() { game.Jump("DISO") })
	}

	if err := app.SetRoot(flex, true).SetFocus(commandPanel).Run(); err != nil {
		panic(err)
	}
}

// Bubble Tea UI experiements

type CommandModel struct {
	viewport viewport.Model
	cmdInput textinput.Model
	done     bool
	game     eliteEngine.Game
	gameCmd  string
}

func (m CommandModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.cmdInput, tiCmd = m.cmdInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

			m.gameCmd = m.cmdInput.Value()
			output, err := m.executeCommand()

			if err != nil {
				output = "There was an error with the command: " + err.Error()
			}

			m.viewport.SetContent(style.Render(output))
			m.cmdInput.Reset()
			m.cmdInput.Blink()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m CommandModel) View() string {

	return fmt.Sprintf(
		"--== Elite v1.5 ==--\n%s\n\nCommand %s",
		m.viewport.View(),
		m.cmdInput.View(),
	) + "\n\n"
}

var _ tea.Model = &CommandModel{}

// Execute a game command
func (m *CommandModel) executeCommand() (string, error) {

	cmdOutput := "Command: " + m.gameCmd + " was requested"

	return cmdOutput, nil
}

func NewCommand(game eliteEngine.Game) *CommandModel {

	cmdPrompt := textinput.New()

	//cmdPrompt.Placeholder = "info"
	cmdPrompt.Focus()

	vp := viewport.New(30, 5)

	return &CommandModel{
		game:     game,
		cmdInput: cmdPrompt,
		viewport: vp,
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
