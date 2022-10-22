package eliteEngine

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// Go implementation of txtelite. See: http://www.iancgbell.clara.net/elite/text/

// Market functions

type TradeGood struct { /* In 6502 version these were: */
	Baseprice uint16 /* one byte */
	Gradient  int16  /* five bits plus sign */
	Basequant uint16 /* one byte */
	Maskbyte  uint16 /* one byte */
	Units     uint   /* two bits */
	Name      string /* longest="Radioactives" */
}

// TODO: Revisit this. Storing the unit names as part of the market struct
// is a departure from the effiecent storage method
// used by the rest of the game and it doesnt really make
// sense to do it like this here.
type Market struct {
	Price     []uint16
	Quantity  []uint16
	UnitNames []string
}

func initCommodities(politicallCorrect bool) []TradeGood {
	commodities := []TradeGood{}

	// Base price,  Gradient, Base Quant, Mask,Unit,Name

	commodities = append(commodities, TradeGood{0x13, -0x02, 0x06, 0x01, 0, "Food"})
	commodities = append(commodities, TradeGood{0x14, -0x01, 0x0A, 0x03, 0, "Textiles"})
	commodities = append(commodities, TradeGood{0x41, -0x03, 0x02, 0x07, 0, "Radioactives"})

	if politicallCorrect {
		commodities = append(commodities, TradeGood{0x28, -0x05, 0xE2, 0x1F, 0, "Robot Slaves"})
		commodities = append(commodities, TradeGood{0x53, -0x05, 0xFB, 0x0F, 0, "Beverages"})

	} else {
		commodities = append(commodities, TradeGood{0x28, -0x05, 0xE2, 0x1F, 0, "Slaves"})
		commodities = append(commodities, TradeGood{0x53, -0x05, 0xFB, 0x0F, 0, "Liquor/Wines"})

	}

	commodities = append(commodities, TradeGood{0xC4, +0x08, 0x36, 0x03, 0, "Luxuries"})

	if politicallCorrect {
		commodities = append(commodities, TradeGood{0xEB, +0x1D, 0x08, 0x78, 0, "Rare Species"})
	} else {
		commodities = append(commodities, TradeGood{0xEB, +0x1D, 0x08, 0x78, 0, "Narcotics"})
	}

	commodities = append(commodities, TradeGood{0x9A, +0x0E, 0x38, 0x03, 0, "Computers"})
	commodities = append(commodities, TradeGood{0x75, +0x06, 0x28, 0x07, 0, "Machinery"})
	commodities = append(commodities, TradeGood{0x4E, +0x01, 0x11, 0x1F, 0, "Alloys"})
	commodities = append(commodities, TradeGood{0x7C, +0x0d, 0x1D, 0x07, 0, "Firearms"})
	commodities = append(commodities, TradeGood{0xB0, -0x09, 0xDC, 0x3F, 0, "Furs"})
	commodities = append(commodities, TradeGood{0x20, -0x01, 0x35, 0x03, 0, "Minerals"})
	commodities = append(commodities, TradeGood{0x61, -0x01, 0x42, 0x07, 1, "Gold"})
	commodities = append(commodities, TradeGood{0xAB, -0x02, 0x37, 0x1F, 1, "Platinum"})
	commodities = append(commodities, TradeGood{0x2D, -0x01, 0xFA, 0x0F, 2, "Gem-Stones"})
	commodities = append(commodities, TradeGood{0x35, +0x0F, 0xC0, 0x07, 0, "Alien Items"})

	return commodities
}

/* Prices and availabilities are influenced by the planet's economy type
   (0-7) and a random "fluctuation" byte that was kept within the saved
   commander position to keep the market prices constant over gamesaves.
   Availabilities must be saved with the game since the player alters them
   by buying (and selling(?))

   Almost all operations are one byte only and overflow "errors" are
   extremely frequent and exploited.

   Trade Item prices are held internally in a single byte=true value/4.
   The decimal point in prices is introduced only when printing them.
   Internally, all prices are integers.
   The player's cash is held in four bytes.
*/

func (p *planetarySystem) generateMarket(commodities []TradeGood, marketFluctuation uint16) {
	mkt := Market{}

	mkt.Quantity = make([]uint16, len(commodities))
	mkt.Price = make([]uint16, len(commodities))

	mkt.UnitNames = []string{"t", "kg", "g"}

	numCommodities := len(commodities) - 1
	AlienItemsIdx := 16

	for i := 0; i <= numCommodities; i++ {
		product := int16((p.Economy)) * (commodities[i].Gradient)
		changing := int16(marketFluctuation & (commodities[i].Maskbyte))
		q := int16((commodities[i].Basequant)) + changing - product
		q = q & 0xFF

		// Clip to positive 8-bit
		// NOTE: Not sure about this. Keep screwing up the bit-wise oprations
		if q&0x80 == 128 {
			q = 0
		}

		mkt.Quantity[i] = uint16(q & 0x3F) // Mask to 6-bits

		q = int16((commodities[i].Baseprice)) + changing + product
		q = q & 0xFF

		mkt.Price[i] = uint16(q * 4)
	}

	mkt.Quantity[AlienItemsIdx] = 0 // NOTE: Why? Override to force nonavailability.

	p.Market = mkt
}

// Buy commodity from the local market
// This isnt the same implementation as the C version, but should be functionally the same
// TODO: Check if we need to deal with units properly

func (g *Game) BuyCommodity(commodity string, amount int) (bought int, err error) {
	bought = 0
	err = nil

	// All these dots feel wrong?
	market := g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet].Market
	commodityIdx := getCommodityIdx(commodity, g.Commodities)

	// Didnt find the commodity is the game commodities list
	if commodityIdx == -1 {
		return 0, errors.New("no such commodity")
	}

	// Not enough to buy
	if market.Quantity[commodityIdx] < uint16(amount) {
		return 0, errors.New("not enough of " + commodity + " in the market to buy")
	}

	// Not enough hold space
	if int(g.Player.Ship.Holdspace) < amount {
		return 0, errors.New("not enough of hold space available to buy " + strconv.Itoa(amount) + " of " + commodity)
	}

	// Not enough cash
	if g.Player.Cash < (int16(market.Price[commodityIdx]) * int16(amount)) {
		return 0, errors.New("not cash to buy " + strconv.Itoa(amount) + " of " + commodity)
	}

	// Everything is in order. Do the trade
	// Add the amount of the commodity to the ships hold
	g.Player.Ship.Hold[commodityIdx] = g.Player.Ship.Hold[commodityIdx] + uint16(amount)

	// Deduct the space from the ships hold space
	g.Player.Ship.Holdspace = g.Player.Ship.Holdspace - uint16(amount)

	// Deduct the amount from the market
	market.Quantity[commodityIdx] = market.Quantity[commodityIdx] - uint16(amount)

	// Deduct the amount spent from players cash
	g.Player.Cash = g.Player.Cash - (int16(market.Price[commodityIdx]) * int16(amount))

	bought = amount
	err = nil

	return bought, err
}

func (g *Game) SellCommodity(commodity string, amount int) (sold int, err error) {
	sold = 0
	err = nil

	// All these dots feel wrong?
	market := g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet].Market
	commodityIdx := getCommodityIdx(commodity, g.Commodities)

	// Didnt find the commodity is the game commodities list
	if commodityIdx == -1 {
		return 0, errors.New("no such commodity")
	}

	// Not enough to buy
	if g.Player.Ship.Hold[commodityIdx] < uint16(amount) {
		return 0, errors.New("not enough of " + commodity + " in the ships hold to sell")
	}

	// Everything is in order. Do the trade
	// Remove the amount of the commodity from the ships hold
	g.Player.Ship.Hold[commodityIdx] = g.Player.Ship.Hold[commodityIdx] - uint16(amount)

	// Add the space from the ships hold space
	g.Player.Ship.Holdspace = g.Player.Ship.Holdspace + uint16(amount)

	// Add the amount to the market
	market.Quantity[commodityIdx] = market.Quantity[commodityIdx] + uint16(amount)

	// add the amount earnt to the players cash
	g.Player.Cash = g.Player.Cash + (int16(market.Price[commodityIdx]) * int16(amount))

	sold = amount
	err = nil

	return sold, err
}

// Helper function to get the index of a commodity in the market
func getCommodityIdx(commodity string, commodities []TradeGood) int {
	for idx, tradegood := range commodities {
		if strings.EqualFold(tradegood.Name, commodity) {
			return idx
		}
	}

	return -1
}

func (g *Game) BuyFuel(amount int16) error {

	// TODO: This can never return an error. Maybe it could, but if not change this function
	// so that it doesnt return anything

	amountToBuy := amount * 10

	currentCapacity := g.maxFuel - g.Player.Ship.Fuel
	//fmt.Printf("Current Capacity: %d amount requested: %d\n", currentCapacity, amountToBuy)

	if amountToBuy > int16(currentCapacity) {
		amountToBuy = int16(currentCapacity)
	}

	// Check if affordable. If not only buy the amount we can afford
	if g.Player.Cash < (int16(g.fuelCost))*int16(amountToBuy) {
		amountToBuy = int16((g.Player.Cash) / int16(g.fuelCost))

	}

	//fmt.Printf("Current Capacity: %d amount requested: %d\n", currentCapacity, amountToBuy)

	// Add amount of fuel to the ship
	g.Player.Ship.Fuel = g.Player.Ship.Fuel + uint16(amountToBuy)

	// Deduct cost from player cash.
	g.Player.Cash = g.Player.Cash - (int16(g.fuelCost) * (amountToBuy))

	return nil
}

// Market Display Functions

func (p *planetarySystem) PrintMarket(commodities []TradeGood) {
	numCommodities := len(commodities) - 1
	mkt := p.Market
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 5, ' ', 0)

	fmt.Fprintln(w, "Local Market")
	fmt.Fprintf(w, "Commodity\tPrice\tQuantity\n")
	fmt.Fprintln(w, "------------------------------------")
	for i := 0; i <= numCommodities; i++ {

		fmt.Fprintf(w, commodities[i].Name)
		fmt.Fprintf(w, "\t%.1f", float64(mkt.Price[i])/float64(10))
		fmt.Fprintf(w, "\t%d", mkt.Quantity[i])
		fmt.Fprintf(w, mkt.UnitNames[commodities[1].Units])
		fmt.Fprintln(w, "")

	}
	w.Flush()
	fmt.Println("------------------------------------")

}

// Returns market data as a string rather than printing to the screen
func (p *planetarySystem) SprintMarket(commodities []TradeGood) string {

	marketData := ""
	numCommodities := len(commodities) - 1
	mkt := p.Market

	marketData = fmt.Sprintf("%s\n\n", "Local Market")
	marketData = marketData + fmt.Sprintf("%-*sPrice  Quantity\n", 21, "Commodity")
	marketData = marketData + fmt.Sprintln("------------------------------------")
	for i := 0; i <= numCommodities; i++ {
		marketData = marketData + fmt.Sprintf("%-*s", 21, commodities[i].Name)
		marketData = marketData + fmt.Sprintf(" %-*.1f", 5, float64(mkt.Price[i])/float64(10))
		marketData = marketData + fmt.Sprintf(" %d", mkt.Quantity[i])
		marketData = marketData + fmt.Sprintf(mkt.UnitNames[commodities[1].Units])
		marketData = marketData + fmt.Sprintln("")

	}
	//marketData = marketData + fmt.Sprintln("------------------------------------")

	return marketData
}
