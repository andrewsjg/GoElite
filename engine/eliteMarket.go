package eliteEngine

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
	commodities = append(commodities, TradeGood{0x2D, -0x01, 0xFA, 0x0F, 2, "Gem-Strones"})
	commodities = append(commodities, TradeGood{0x35, +0x0F, 0xC0, 0x07, 0, "Alien Items"})

	return commodities
}
