package main

import (
	"fmt"

	engine "github.com/andrewsjg/GoElite/engine"
)

func main() {

	galaxy := engine.InitGalaxy(1)
	debugTests(galaxy)

}

func debugTests(gal engine.Galaxy) {

	fmt.Printf("The current system is: %s", gal.Systems[gal.CurrentPlanet].Name)
	// test current planet (Lave at start)
	gal.PrintSystem(gal.Systems[gal.CurrentPlanet], false)
	fmt.Println()
	// test matchsys
	fmt.Printf("DISO is system numner: %d", gal.Matchsys("DISO")) // 147

	diso := gal.Matchsys("DISO")
	gal.PrintSystem(gal.Systems[diso], false)

}
