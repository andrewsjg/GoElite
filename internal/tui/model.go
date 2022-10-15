package tui

import (
	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/knipferrc/teacup/statusbar"
)

// TODO: Rename the model. This doesnt make sense any more
type Tui struct {
	systemViewport viewport.Model
	marketViewport viewport.Model
	cmdInput       textinput.Model
	statusBar      statusbar.Bubble
	game           eliteEngine.Game
	gameCmd        string
}
