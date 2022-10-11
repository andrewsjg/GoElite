package eliteEngine

import (
	"fmt"
	"math/rand"
)

// Go implementation of txtelite. See: http://www.iancgbell.clara.net/elite/text/
// Core game function

type player struct {
	Ship ship
	Cash int16
}

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
	Galaxy Galaxy
	//PlayerShip    ship
	Player        player
	maxFuel       uint16
	fuelCost      uint16
	lastrand      uint
	useNativeRand bool
	AlienItems    uint16
	Commodities   []TradeGood
}

type NavInfo struct {
	System                   planetarySystem
	ReachableWithMaxFuel     bool
	ReachableWithCurrentFuel bool
	Distance                 int
}

func (g *Game) gameRand() uint {

	if g.useNativeRand {
		return uint(rand.Intn(65536))
	}

	gRand := (((((((((((g.lastrand << 3) - g.lastrand) << 3) + g.lastrand) << 1) + g.lastrand) << 4) - g.lastrand) << 1) - g.lastrand) + 0xe60) & 0x7fffffff
	g.lastrand = gRand - 1

	return gRand
}

func (g *Game) randByte() uint {
	return g.gameRand() & 0xFF
}

func InitGame(useNativeRand bool) Game {
	game := Game{}
	ship := ship{}
	player := player{}

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

	player.Ship = ship
	player.Cash = 1000

	// Seed the local market. This needs to be done on each jump
	game.Galaxy.Systems[ship.Location.CurrentPlanet].marketFluctuation = 0

	game.Player = player
	game.AlienItems = 16 // Number of commodities per market
	game.Commodities = initCommodities(true)
	game.lastrand = 0
	game.useNativeRand = useNativeRand

	return game
}

// Game functions

// TODO: Probably want this to return an error rather than printing things out
func (g *Game) Jump(planetName string) {
	dest := g.Galaxy.Matchsys(planetName)

	if dest == g.Player.Ship.Location.CurrentPlanet {
		fmt.Println("Bad Jump")
		return
	}

	dist := distance(g.Galaxy.Systems[dest], g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet])

	fmt.Printf("Jump Distance: %d\n", dist)
	fmt.Printf("Current Fuel: %d\n", g.Player.Ship.Fuel)
	if dist > int(g.Player.Ship.Fuel) {
		fmt.Println("To far to jump. Not enough fuel")
		return
	}

	g.Player.Ship.Fuel = g.Player.Ship.Fuel - uint16(dist)
	g.Player.Ship.Location.CurrentPlanet = dest
	g.Galaxy.Systems[dest].marketFluctuation = uint16(g.randByte())
}

// Jump to a new Galaxy
func (g *Game) HyperSpaceJump() {

	g.Galaxy.galaxyNum = g.Galaxy.galaxyNum + 1

	if g.Galaxy.galaxyNum == 9 {
		g.Galaxy.galaxyNum = 1
	}

	g.Galaxy.buildGalaxy(g.Galaxy.galaxyNum)
}

// Show local systems. Replicates the orginal functionality.
func (g *Game) ShowLocal() {
	fmt.Printf("Galaxy Number: %d", g.Galaxy.galaxyNum)

	reachable := g.ReachableSystems()

	for _, navinfo := range reachable {
		if navinfo.ReachableWithCurrentFuel {
			fmt.Printf("\n *")
		} else {
			fmt.Printf("\n - ")
		}
		g.Galaxy.PrintSystem(navinfo.System, true)
		fmt.Printf(" (%.1f LY)", float64(navinfo.Distance)/float64(10))
	}

	fmt.Println()
}

// Return and array of reachable systems
func (g *Game) ReachableSystems() []NavInfo {
	reachable := []NavInfo{}

	currentPlanent := g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet]

	for syscount := 0; syscount < g.Galaxy.Size; syscount++ {

		dist := distance(g.Galaxy.Systems[syscount], currentPlanent)

		if dist <= int(g.maxFuel) {
			nav := NavInfo{}

			if dist <= int(g.Player.Ship.Fuel) {
				nav.ReachableWithCurrentFuel = true
				nav.ReachableWithMaxFuel = true
			} else {
				nav.ReachableWithCurrentFuel = false
				nav.ReachableWithMaxFuel = true
			}

			nav.System = g.Galaxy.Systems[syscount]
			nav.Distance = dist

			reachable = append(reachable, nav)
		}
	}

	return reachable

}

// Print out the current state of the game. Mostly for debug
func (g *Game) PrintState() {
	gal := g.Galaxy
	shipLocation := g.Player.Ship.Location
	planet := gal.Systems[shipLocation.CurrentPlanet]

	fmt.Printf("Current System is: %s", planet.Name)
	g.Galaxy.PrintSystem(planet, false)
	fmt.Println()
	planet.PrintMarket(g.Commodities)
	fmt.Println()
	fmt.Printf("Cash: \t\t%.1f\n", float64(g.Player.Cash)/float64(10))
	fmt.Printf("Fuel: \t\t%.1f\n", float64(g.Player.Ship.Fuel)/float64(10))
	fmt.Printf("Hold Space: \t%dt\n\n", g.Player.Ship.Holdspace)

}
