package main

import (
	"fmt"

	engine "github.com/andrewsjg/GoElite/engine"
)

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
