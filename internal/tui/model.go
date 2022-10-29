package tui

import (
	eliteEngine "github.com/andrewsjg/GoElite/engine"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/knipferrc/teacup/statusbar"
)

type Tui struct {
	systemViewport viewport.Model
	marketViewport viewport.Model
	cmdrViewport   viewport.Model
	fuelBar        progress.Model
	holdSpaceBar   progress.Model
	cmdInput       textinput.Model
	statusBar      statusbar.Bubble
	holdTable      table.Model
	game           eliteEngine.Game
	gameCmd        string
}
