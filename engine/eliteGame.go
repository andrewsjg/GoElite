package eliteEngine

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Go implementation of txtelite. See: http://www.iancgbell.clara.net/elite/text/

// Core game functions and objects

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
	Galaxy        Galaxy
	Player        player
	maxFuel       uint16
	fuelCost      uint16
	lastrand      uint
	useNativeRand bool
	AlienItems    uint16
	Commodities   []TradeGood
	GameCommands  GameCommands
}

// A gamecommand is a map of commands to the functions that execute the command
// TODO: Shoult func return a generic type here? For now returning strings works because all
// the result of a command is is some status output that goes on screen
type GameCommands map[string]func(game *Game, args ...[]string) (stauts string, output string)

type NavInfo struct {
	System                   planetarySystem
	ReachableWithMaxFuel     bool
	ReachableWithCurrentFuel bool
	Distance                 int
}

// Builds a map of game commands.
// Each command is a function type that can be called to execute the command
func buildGameCommands(game *Game) GameCommands {

	gc := GameCommands{}

	// Info Command

	infoCmd := func(game *Game, args ...[]string) (status string, output string) {
		return "", game.SprintState()
	}
	gc["info"] = infoCmd

	// Jump Command

	jumpCmd := func(game *Game, args ...[]string) (status string, output string) {
		// Because this is variadic and how I am passing the arguments as a split string,
		// the args are an array of string arrays. This makes this code a bit wierd. Might be a way around this?

		jumpResult := "Unknown Destination"

		//fmt.Println(args)
		if len(args[0]) >= 2 {

			dest := args[0][1]
			err := game.Jump(dest)

			if err != nil {
				jumpResult = err.Error()
			} else {
				jumpResult = "Jump Complete"
			}
		}

		status = jumpResult
		//output = game.SprintState()

		return status, output
	}
	gc["jump"] = jumpCmd

	mktCmd := func(game *Game, args ...[]string) (status string, output string) {

		system := game.Galaxy.Systems[game.Player.Ship.Location.CurrentPlanet]
		output = system.SprintMarket(game.Commodities)
		status = "OK"

		return status, output
	}

	gc["mkt"] = mktCmd

	buyCmd := func(game *Game, args ...[]string) (status string, output string) {

		if len(args[0]) >= 3 {

			commodity := args[0][1]
			amount, err := strconv.Atoi(args[0][2])

			if err != nil {
				amount = 0
				return "Buy Failed", output
			}

			if strings.ToUpper(commodity) == "FUEL" {
				err = game.BuyFuel(int16(amount))

				if err != nil {
					status = "Buy fuel failed. " + err.Error()
					return status, output
				}

				return "Bought Fuel", output
			}

			bought, err := game.BuyCommodity(commodity, amount)

			if err != nil {
				status = "Buy Failed. " + err.Error()
				return status, output
			}

			status = "Bought " + fmt.Sprint(bought) + " tonne of " + commodity
		}

		return status, output
	}

	gc["buy"] = buyCmd

	sellCmd := func(game *Game, args ...[]string) (status string, output string) {

		if len(args[0]) >= 3 {
			commodity := args[0][1]

			amount, err := strconv.Atoi(args[0][2])

			if err != nil {
				amount = 0
				return "Sell Failed", output
			}

			bought, err := game.SellCommodity(commodity, amount)

			if err != nil {
				status = "Sell Failed. " + err.Error()
				return status, output
			}

			status = "Sold " + fmt.Sprint(bought) + " tonne of " + commodity
		}

		return status, output
	}

	gc["sell"] = sellCmd

	hyperCmd := func(game *Game, args ...[]string) (status string, output string) {

		return status, output
	}

	gc["hyper"] = hyperCmd

	localCmd := func(game *Game, args ...[]string) (status string, output string) {

		return status, output
	}

	gc["local"] = localCmd

	helpCmd := func(game *Game, args ...[]string) (status string, output string) {

		return status, output
	}

	gc["help"] = helpCmd
	return gc
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
	ship.Location.CurrentPlanet = 7 //game.Galaxy.CurrentPlanet

	player.Ship = ship
	player.Cash = 1000

	game.Commodities = initCommodities(true)

	// Generate the local market. This needs to be done on each jump
	game.Galaxy.Systems[ship.Location.CurrentPlanet].generateMarket(game.Commodities, 0)

	game.Player = player
	game.lastrand = 0
	game.useNativeRand = useNativeRand

	game.GameCommands = buildGameCommands(&game)

	return game
}

// Game functions

func (g *Game) Jump(planetName string) error {
	dest := g.Matchsys(strings.ToUpper(planetName))

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
	return g.Galaxy.Systems[g.Matchsys(systemName)]
}

// Game Display Functions

func (g *Game) GetPlanetaryData(systemName string) PlanetarySystem {

	systemData := PlanetarySystem{}

	gal := g.Galaxy
	dataTables := gal.dataTables
	psys := gal.Systems[g.Matchsys(systemName)]

	gal.prng.rnd_seed = psys.goatsoupseed
	gs := gal.sgoatSoup("", "\x8F is \x97.", &psys)

	systemData.Description = gs
	systemData.Techlev = psys.Techlev + 1
	systemData.Economy = psys.Economy
	systemData.EconomyName = dataTables.econnames[psys.Economy]
	systemData.Govtype = psys.Govtype
	systemData.GovName = dataTables.govnames[psys.Govtype]
	systemData.Productivity = psys.Productivity
	systemData.Population = psys.Population >> 3
	systemData.Radius = psys.Radius
	systemData.X = psys.X
	systemData.Y = psys.Y

	return systemData
}

// Print out the current state of the game. Mostly for debug or simple CLI
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

// Return game state as a string
func (g *Game) SprintState() string {

	gameState := ""

	gal := g.Galaxy
	shipLocation := g.Player.Ship.Location
	planet := gal.Systems[shipLocation.CurrentPlanet]

	gameState = fmt.Sprintf("%s\n\n", "System Info")
	gameState = gameState + g.Galaxy.SprintSystem(planet, false)
	//gameState = gameState + planet.SprintMarket(g.Commodities) + "\n"

	gameState = gameState + fmt.Sprintf("\n%s %.1f\n", "Cash:", float64(g.Player.Cash)/float64(10))
	gameState = gameState + fmt.Sprintf("%s %.1f\n", "Fuel:", float64(g.Player.Ship.Fuel)/float64(10))
	gameState = gameState + fmt.Sprintf("%s %dt", "Hold Space:", g.Player.Ship.Holdspace)

	return gameState

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
