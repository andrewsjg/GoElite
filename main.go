package main

import (
	"fmt"

	engine "github.com/andrewsjg/GoElite/engine"
)

// Currently this does nothing except initialise the game as per the defaults
// and run some simply commands for debug

// TODO:  Add CLI and TUI to make game playable

func main() {

	game := engine.InitGame(false)
	debugTests(game)

}

func debugTests(game engine.Game) {

	// Initial State
	game.PrintState()

	// Jump to DISO
	err := game.Jump("DISO")

	if err != nil {
		fmt.Println("Jump failed: " + err.Error())
	}

	game.PrintState()

	game.PrintLocal()

	fmt.Printf("\nDoing Hyperspace Jump\n\n")
	game.HyperSpaceJump()

	game.PrintLocal()
}
