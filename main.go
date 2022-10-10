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

	gal := game.Galaxy
	shipLocation := game.PlayerShip.Location

	fmt.Printf("Ship is currently at: %s in galaxy %d", gal.Systems[shipLocation.CurrentPlanet].Name, shipLocation.CurrentGalaxy)

	// test current planet (Lave at start)
	gal.PrintSystem(gal.Systems[shipLocation.CurrentPlanet], false)
	fmt.Println()

	// test matchsys
	fmt.Printf("DISO is system numner: %d", gal.Matchsys("DISO")) // 147

	diso := gal.Matchsys("DISO")
	gal.PrintSystem(gal.Systems[diso], false)

}
