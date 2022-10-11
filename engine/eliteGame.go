package eliteEngine

import (
	"errors"
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

	game.maxFuel = 70    // 7 Light Year tank
	game.fuelCost = 2    // 0.2 CR/Light year
	game.AlienItems = 16 // Number of commodities per market

	ship.Hold = make([]uint16, game.AlienItems+1)
	ship.Holdspace = 20 //uint16(len(ship.Hold))
	ship.Fuel = game.maxFuel

	game.Galaxy = initGalaxy(1)

	// Start in Galaxy 1 at Lave
	// Not sure I need this bit?
	ship.Location.CurrentGalaxy = game.Galaxy.galaxyNum
	ship.Location.CurrentPlanet = game.Galaxy.CurrentPlanet

	player.Ship = ship
	player.Cash = 1000

	game.Commodities = initCommodities(true)

	// Generate the local market. This needs to be done on each jump
	game.Galaxy.Systems[ship.Location.CurrentPlanet].generateMarket(game.Commodities, 0)

	game.Player = player
	game.lastrand = 0
	game.useNativeRand = useNativeRand

	return game
}

// Game functions

func (g *Game) Jump(planetName string) error {
	dest := g.Galaxy.Matchsys(planetName)

	if dest == g.Player.Ship.Location.CurrentPlanet {
		return errors.New("bad jump")
	}

	dist := distance(g.Galaxy.Systems[dest], g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet])

	if dist > int(g.Player.Ship.Fuel) {
		return errors.New("to far to jump. Not enough fuel")
	}

	g.Player.Ship.Fuel = g.Player.Ship.Fuel - uint16(dist)
	g.Player.Ship.Location.CurrentPlanet = dest

	// Generate the local market. This is a bit ugly
	g.Galaxy.Systems[dest].generateMarket(g.Commodities, uint16(g.randByte()))

	return nil
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
func (g *Game) PrintLocal() {
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

// Function that returns a planetary system record. Can be used for print system info
func (g *Game) GetSystemData(systemName string) planetarySystem {
	return g.Galaxy.Systems[g.Galaxy.Matchsys(systemName)]
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
