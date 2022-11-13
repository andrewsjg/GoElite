package main

import (
	"fmt"
	"os"
	"strings"

	engine "github.com/andrewsjg/GoElite/engine"
	"github.com/andrewsjg/GoElite/internal/tui"
)

func main() {

	//game := engine.InitGame(false)
	//debugTests(game)

	// TODO: Probably a better way to encode version required
	const VER = "0.1.0"

	args := os.Args[1:]

	if len(args) >= 1 {
		arg0 := strings.ToUpper(args[0])

		if arg0 == "--VERSION" || arg0 == "--VER" {
			fmt.Println(VER)
			return
		}
	}

	tui.Start()

}

func debugTests(game engine.Game) {

	// Initial State
	//game.PrintState()

	plan := game.Galaxy.Systems[game.Player.Ship.Location.CurrentPlanet]

	fmt.Println(plan.SprintMarket(game.Commodities))

	/* Jump to DISO
	err := game.Jump("DISO")

	if err != nil {
		fmt.Println("Jump failed: " + err.Error())
	}

	game.PrintState()

	game.PrintLocal()

	fmt.Printf("\nDoing Hyperspace Jump\n\n")
	game.HyperSpaceJump()

	game.PrintLocal() */

}
