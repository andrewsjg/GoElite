package eliteEngine

import (
	"testing"
)

func TestGetCommodityIdx(t *testing.T) {

	game := InitGame(false)

	commodityName := "alloys"

	alloyIdx := getCommodityIdx(commodityName, game.Commodities)

	// Make sure we can find Alloys
	if alloyIdx == -1 {
		t.Error("Alloys dont exist in the commidities list")
	}

	if alloyIdx != 9 {
		t.Error("Alloys at the wrong index. Should be 9. Got: ", alloyIdx)
	}

	commodityName = "aloys"
	alloyIdx = getCommodityIdx(commodityName, game.Commodities)

	// Make sure getCommodityIdx works as expected for non-existent commodities
	if alloyIdx != -1 {
		t.Error("Found non-existent item in the commodity list")
	}

}

func TestLaveMarketState(t *testing.T) {
	game := InitGame(false)
	market := game.Galaxy.Systems[game.Player.Ship.Location.CurrentPlanet].Market

	commodityName := "alloys"
	alloyIdx := getCommodityIdx(commodityName, game.Commodities)

	alloyQuant := market.Quantity[alloyIdx]
	alloyPrice := market.Price[alloyIdx]

	if alloyQuant != 12 {
		t.Error("Expected 12t of alloys. Got:", alloyQuant)
	}

	if alloyPrice != 332 {
		t.Error("Expected alloys of 332. Got:", alloyPrice)
	}

}

func TestMarketBuySell(t *testing.T) {
	game := InitGame(false)
	market := game.Galaxy.Systems[game.Player.Ship.Location.CurrentPlanet].Market

	commodityName := "alloys"
	alloyIdx := getCommodityIdx(commodityName, game.Commodities)

	if market.Quantity[alloyIdx] != 12 {
		t.Error("Expected LAVE to have 12t of Alloys. Got:", market.Quantity[alloyIdx])
	}

	// Check the  starting values are correct
	if game.Player.Cash != 1000 {
		t.Error("Expected Player to have 1000 cash. Got:", game.Player.Cash)
	}

	if game.Player.Ship.Holdspace != 20 {
		t.Error("Expected ship to have 20 holdspace. Got:", game.Player.Ship.Holdspace)
	}

	// Buy 2 tonnes of alloys
	game.BuyCommodity(commodityName, 2)

	if market.Quantity[alloyIdx] != 10 {
		t.Error("Expected LAVE to have 10t of Alloys. Got:", market.Quantity[alloyIdx])
	}

	if game.Player.Cash != 336 {
		t.Error("Expected Player to have 336 cash. Got:", game.Player.Cash)
	}

	if game.Player.Ship.Hold[9] != 2 {
		t.Error("Expected ships hold to have 2t of Alloys. Got:", game.Player.Ship.Hold[9])
	}

	if game.Player.Ship.Holdspace != 18 {
		t.Error("Expected ship to have 18 holdspace. Got:", game.Player.Ship.Holdspace)
	}

	// Sell 2 tonnes of alloys

	game.SellCommodity(commodityName, 2)

	if market.Quantity[alloyIdx] != 12 {
		t.Error("Expected LAVE to have 12t of Alloys. Got:", market.Quantity[alloyIdx])
	}

	if game.Player.Cash != 1000 {
		t.Error("Expected Player to have 1000 cash. Got:", game.Player.Cash)
	}

	if game.Player.Ship.Hold[9] != 0 {
		t.Error("Expected ships hold to have 0t of Alloys. Got:", game.Player.Ship.Hold[9])
	}

	if game.Player.Ship.Holdspace != 20 {
		t.Error("Expected ship to have 20 holdspace. Got:", game.Player.Ship.Holdspace)
	}

	// Try to sell something we dont have
	_, err := game.SellCommodity(commodityName, 2)
	if err == nil {
		t.Error("Attempted to sell a commodity we dont have. No error raised")
	}

}

/*
LAVE Market

Food             3.6      16t
Textiles         6.0      15t
Radioactives     20.0     17t
Robot Slaves     6.0      0t
Beverages        23.2     20t
Luxuries         94.4     14t
Rare Species     49.6     55t
Computers        89.6     0t
Machinery        58.8     10t
Alloys           33.2     12t
Firearms         75.6     0t
Furs             52.4     9t
Minerals         10.8     58t
Gold             36.8     7t
Platinum         64.4     1t
Gem-Stones       16.0     0t
Alien Items      51.2     0t
*/
