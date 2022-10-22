package eliteEngine

import (
	"strings"
	"testing"
)

func TestJump(t *testing.T) {
	game := InitGame(false)

	// Jump to diso
	game.Jump("DISO") // Remaining Fuel 34. Cost 36
	reachable := game.ReachableSystems()

	// Should be 8 reachable systems
	if len(reachable) != 8 {
		t.Error("expected 8 reachable systems. Got:", len(reachable))
	}

	// Because of a problem with the way the names are generated using byte arrays there
	// seems to be a stray 0 at the end of the byte array for the system name (golang rune thing?).
	// Need to strip that stray 0 byte before the string compare works

	bName := []byte(reachable[2].System.Name)
	name := string(bName[:len(bName)-1])

	// Second reachable system is RIEDQUAT. Tests planet name generation.
	if !strings.EqualFold(string(name), "RIEDQUAT") {
		t.Error("expected the second reachable system to be RIEDQUAT. Got:", reachable[2].System.Name)
	}

}

func TestFuel(t *testing.T) {
	game := InitGame(false)

	// Jump to diso
	game.Jump("DISO") // Remaining Fuel 34. Cost 36

	// Check remaining fuel
	if game.Player.Ship.Fuel != 34 {
		t.Error("expected to have 34 fuel left. Got:", game.Player.Ship.Fuel)
	}

	// Attempt to by 40 fuel. Should fill the tank to a max of 70
	err := game.BuyFuel(40)

	if err != nil {
		t.Error("Buy fuel cause and error:", err.Error())
	}

	if game.Player.Cash != 928 {
		t.Error("Fuel bought but didnt cost what was expected. Cash:", game.Player.Cash)
	}

	if game.Player.Ship.Fuel != 70 {
		t.Error("expected 70 fuel. Got:", game.Player.Ship.Fuel)
	}

	// Jump to around a bit testing fuel along the way
	game.Jump("LAVE")

	// This should error
	err = game.Jump("DISO")
	if err == nil {
		t.Error("Able to jump without enough fuel")
	}

	err = game.BuyFuel(36)
	if err != nil {
		t.Error("Error buying fuel:", err.Error())
	}

	if game.Player.Ship.Fuel != 70 {
		t.Error("expected 70 fuel. Got:", game.Player.Ship.Fuel)
	}

	game.Jump("LAVE")

	err = game.BuyFuel(36)
	if err != nil {
		t.Error("Error buying fuel:", err.Error())
	}

	if game.Player.Ship.Fuel != 70 {
		t.Error("expected 70 fuel. Got:", game.Player.Ship.Fuel)
	}

	game.Jump("RIEDQUAT")
	if game.Player.Ship.Fuel != 1 {
		t.Error("expected 1 fuel. Got:", game.Player.Ship.Fuel)
	}

	// Drain cash
	game.Player.Cash = 0

	// Attempt to buy fuel with no money. This will error
	err = game.BuyFuel(70)
	if err == nil {
		t.Error("Should not have been able to buy fuel", err.Error())
	}

}
