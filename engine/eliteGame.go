package eliteEngine

// Go implementation of txtelite. See: http://www.iancgbell.clara.net/elite/text/
// Core game function

type shipLocation struct {
	CurrentPlanet int
	CurrentGalaxy int
}

type ship struct {
	Hold      []uint16
	Holdspace uint16
	Fuel      uint16
	Location  shipLocation
}

type Game struct {
	Galaxy      Galaxy
	PlayerShip  ship
	maxFuel     uint16
	fuelCost    uint16
	AlienItems  uint16
	Commodities []TradeGood
}

func InitGame() Game {
	game := Game{}
	ship := ship{}

	game.maxFuel = 70 // 7 Light Year tank
	game.fuelCost = 2 // 0.2 CR/Light year

	// TODO: Fix this so its configurable
	ship.Hold = make([]uint16, 20)
	ship.Holdspace = uint16(len(ship.Hold))
	ship.Fuel = game.maxFuel

	game.Galaxy = initGalaxy(1)

	// Start in Galaxy 1 at Lave
	// Not sure I need this bit?
	ship.Location.CurrentGalaxy = game.Galaxy.galaxyNum
	ship.Location.CurrentPlanet = game.Galaxy.CurrentPlanet

	game.PlayerShip = ship
	game.AlienItems = 16 // Number of commodities per market
	game.Commodities = initCommodities(true)

	return game
}
