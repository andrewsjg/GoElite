package main

import (
	engine "github.com/andrewsjg/GoElite/engine"
)

func main() {

	game := engine.InitGame(false)
	debugTests(game)

}

func debugTests(game engine.Game) {

	/*gal := game.Galaxy
	shipLocation := game.Player.Ship.Location
	currentPlanet := gal.Systems[shipLocation.CurrentPlanet]

	fmt.Printf("Ship is currently at: %s in galaxy %d", gal.Systems[shipLocation.CurrentPlanet].Name, shipLocation.CurrentGalaxy)

	// Test current planet (Lave at start)
	gal.PrintSystem(gal.Systems[shipLocation.CurrentPlanet], false)
	fmt.Println()

	// Print the local Market
	fmt.Println()
	currentPlanet.PrintMarket(game.Commodities) */

	// Initial State
	game.PrintState()

	// Jump to DISO
	game.Jump("DISO")
	game.PrintState()
}
