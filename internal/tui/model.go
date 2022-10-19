package tui

import (
	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/knipferrc/teacup/statusbar"
)

type Tui struct {
	systemViewport viewport.Model
	marketViewport viewport.Model
	shipViewport   viewport.Model
	fuelBar        progress.Model
	cmdInput       textinput.Model
	statusBar      statusbar.Bubble
	game           eliteEngine.Game
	gameCmd        string
}
