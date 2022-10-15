package eliteEngine

// Go implementation of txtelite. See: http://www.iancgbell.clara.net/elite/text/
// Galaxy functions

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// four byte random number used for planet description
type fastseed struct {
	a uint16
	b uint16
	c uint16
	d uint16
}

// six byte random number used as seed for planets
type seed struct {
	w0 uint16
	w1 uint16
	w2 uint16
}

type planetarySystem struct {
	X            uint16
	Y            uint16 /* One byte unsigned */
	Economy      uint16 /* These two are actually only 0-7  */
	Govtype      uint16
	Techlev      uint16 /* 0-16 i think */
	Population   uint16 /* One byte */
	Productivity uint16 /* Two byte */
	Radius       uint16 /* Two byte (not used by game at all) */
	goatsoupseed fastseed
	Name         string
	Market       Market
}

type Galaxy struct {
	galaxyNum int
	Size      int
	Systems   []planetarySystem

	prng       *galPRNG
	dataTables planetDataTables
}

type galPRNG struct {
	// PRNG For galaxy generation
	galaxySeed seed
	rnd_seed   fastseed
	base0      uint16
	base1      uint16
	base2      uint16

	lastrand int64
}

type planetDataTables struct {
	pairs0    []byte
	pairs     []byte
	govnames  []string
	econnames []string
}

func initGalaxy(galaxyNumber int) Galaxy {
	galaxy := Galaxy{}

	// Data Tables used to generate planet data
	dataTables := planetDataTables{}

	dataTables.econnames = []string{"Rich Ind", "Average Ind", "Poor Ind", "Mainly Ind",
		"Mainly Agri", "Rich Agri", "Average Agri", "Poor Agri"}

	dataTables.govnames = []string{"Anarchy", "Feudal", "Multi-gov", "Dictatorship",
		"Communist", "Confederacy", "Democracy", "Corporate State"}

	dataTables.pairs0 = []byte("ABOUSEITILETSTONLONUTHNOALLEXEGEZACEBISOUSESARMAINDIREA.ERATENBERALAVETIEDORQUANTEISRION")
	dataTables.pairs = []byte("..LEXEGEZACEBISOUSESARMAINDIREA.ERATENBERALAVETIEDORQUANTEISRION") /* Dots should be nullprint characters */
	galaxy.dataTables = dataTables

	// Galaxy PRNG
	galRNG := galPRNG{}

	galRNG.base0 = 0x5A4A
	galRNG.base1 = 0x0248
	galRNG.base2 = 0xB753
	galRNG.lastrand = 0
	// Seend the RNG
	galRNG.mysrand(12345)

	// Galaxy parameters
	galaxy.Size = 256 // Should pass this as a parameter?
	galaxy.galaxyNum = galaxyNumber
	galaxy.prng = &galRNG
	galaxy.Systems = make([]planetarySystem, galaxy.Size)

	// Populate the galaxy with planetary systems
	galaxy.buildGalaxy(galaxy.galaxyNum)

	return galaxy
}

func (g *galPRNG) gen_rnd_number() uint16 {
	var a, x uint16
	x = (g.rnd_seed.a * 2) & 0xFF
	a = x + g.rnd_seed.c

	if g.rnd_seed.a > 127 {
		a++
	}

	g.rnd_seed.a = a & 0xFF
	g.rnd_seed.c = x

	a = a / 256 /* a = any carry left from above */
	x = g.rnd_seed.b
	a = (a + x + g.rnd_seed.d) & 0xFF
	g.rnd_seed.b = a
	g.rnd_seed.d = x

	return a
}

func (g *galPRNG) mysrand(seed int64) {
	rand.Seed(seed)

	g.lastrand = seed - 1
}

func (g *galPRNG) tweakseed(s *seed) {

	temp := (s.w0) + (s).w1 + (s.w2) /* 2 byte aritmetic */
	(*s).w0 = (*s).w1
	(*s).w1 = (*s).w2
	(*s).w2 = temp
}

/* rotate 8 bit number leftwards */

func rotatel(x uint16) uint16 {
	return ((x << 1) & 0xfe) | ((x >> 7) & 0x01)
}

func twist(x uint16) uint16 {
	return (uint16)((256 * rotatel(x>>8)) + rotatel(x&255))
}

/*
Apply to base seed; once for galaxy 2

	twice for galaxy 3, etc.
	Eighth application gives galaxy 1 again
*/
func (g *Galaxy) nextgalaxy(s *seed) {
	s.w0 = twist(s.w0)
	s.w1 = twist(s.w1)
	s.w2 = twist(s.w2)
}

func (g *Galaxy) buildGalaxy(galaxyNum int) {
	var syscount, galcount int

	/* Initialise seed for galaxy 1 */
	g.prng.galaxySeed.w0 = g.prng.base0
	g.prng.galaxySeed.w1 = g.prng.base1
	g.prng.galaxySeed.w2 = g.prng.base2

	for galcount = 1; galcount < galaxyNum; galcount++ {
		g.nextgalaxy(&g.prng.galaxySeed)
	}

	/* Put galaxy data into array of structures */
	for syscount = 0; syscount < g.Size; syscount++ {
		g.Systems[syscount] = g.makeSystem(&g.prng.galaxySeed)

	}
}

func (g *Galaxy) makeSystem(s *seed) planetarySystem {

	var pair1, pair2, pair3, pair4 uint16
	var longnameflag uint16 = (s.w0) & 0x40

	planSys := planetarySystem{}

	planSys.X = ((s.w1) >> 8)
	planSys.Y = ((s.w0) >> 8)

	planSys.Govtype = (((s.w1) >> 3) & 7) /* bits 3,4 &5 of w1 */

	planSys.Economy = (((s.w0) >> 8) & 7) /* bits 8,9 &A of w0 */
	if planSys.Govtype <= 1 {
		planSys.Economy = ((planSys.Economy) | 2)
	}

	planSys.Techlev = (((s.w1) >> 8) & 3) + ((planSys.Economy) ^ 7)
	planSys.Techlev += ((planSys.Govtype) >> 1)

	if ((planSys.Govtype) & 1) == 1 {
		planSys.Techlev += 1
	}
	/* C simulation of 6502's LSR then ADC */

	planSys.Population = 4*(planSys.Techlev) + (planSys.Economy)
	planSys.Population += (planSys.Govtype) + 1

	planSys.Productivity = (((planSys.Economy) ^ 7) + 3) * ((planSys.Govtype) + 4)
	planSys.Productivity *= (planSys.Population) * 8

	planSys.Radius = 256*((((s.w2)>>8)&15)+11) + planSys.X

	planSys.goatsoupseed.a = (s.w1 & 0xFF)
	planSys.goatsoupseed.b = (s.w1 >> 8)
	planSys.goatsoupseed.c = (s.w2 & 0xFF)
	planSys.goatsoupseed.d = (s.w2 >> 8)

	pair1 = 2 * (((s.w2) >> 8) & 31)
	g.prng.tweakseed(s)

	pair2 = 2 * (((s.w2) >> 8) & 31)
	g.prng.tweakseed(s)

	pair3 = 2 * (((s.w2) >> 8) & 31)
	g.prng.tweakseed(s)

	pair4 = 2 * (((s.w2) >> 8) & 31)
	g.prng.tweakseed(s)
	/* Always four iterations of random number */

	name := make([]byte, 9)

	name[0] = g.dataTables.pairs[pair1]
	name[1] = g.dataTables.pairs[pair1+1]
	name[2] = g.dataTables.pairs[pair2]
	name[3] = g.dataTables.pairs[pair2+1]
	name[4] = g.dataTables.pairs[pair3]
	name[5] = g.dataTables.pairs[pair3+1]

	/* bit 6 of ORIGINAL w0 flags a four-pair name */
	if longnameflag == 64 {
		name[6] = g.dataTables.pairs[pair4]
		name[7] = g.dataTables.pairs[pair4+1]
		name[8] = 0

	} else {
		name[6] = 0
	}

	planSys.Name = strings.ReplaceAll(string(name), ".", "")

	return planSys
}

/* Return id of the planet whose name matches passed strinmg
   closest to currentplanet - if none return currentplanet

func (g *Galaxy) Matchsys(platnetName string) int {

	p := g.CurrentPlanet
	d := 9999

	for syscount := 0; syscount < g.Size; syscount++ {

		if strings.HasPrefix(g.Systems[syscount].Name, platnetName) {
			if distance(g.Systems[syscount], g.Systems[g.CurrentPlanet]) < d {

				d = distance(g.Systems[syscount], g.Systems[g.CurrentPlanet])
				p = syscount
			}
		}
	}

	return p
}  */

func (g *Game) Matchsys(platnetName string) int {

	p := g.Player.Ship.Location.CurrentPlanet
	currentPlanet := p
	gal := g.Galaxy
	d := 9999

	for syscount := 0; syscount < gal.Size; syscount++ {

		if strings.HasPrefix(gal.Systems[syscount].Name, platnetName) {
			if distance(gal.Systems[syscount], gal.Systems[currentPlanet]) < d {

				d = distance(gal.Systems[syscount], gal.Systems[currentPlanet])
				p = syscount
			}
		}
	}

	return p
}

// Seperation between two planets (4*sqrt(X*X+Y*Y/4))
func distance(planetA planetarySystem, planetB planetarySystem) int {

	val := (planetA.X-planetB.X)*(planetA.X-planetB.X) + (planetA.Y-planetB.Y)*(planetA.Y-planetB.Y)/4.0

	return int(4 * math.Sqrt(float64(val)))
}

func (g *Galaxy) goatSoup(source string, psy *planetarySystem) {

	desc := [][]string{
		{"fabled", "notable", "well known", "famous", "noted"},
		{"very", "mildly", "most", "reasonably", ""},
		{"ancient", "\x95", "great", "vast", "pink"},
		{"\x9E \x9D plantations", "mountains", "\x9C", "\x94 forests", "oceans"},
		{"shyness", "silliness", "mating traditions", "loathing of \x86", "love for \x86"},
		{"food blenders", "tourists", "poetry", "discos", "\x8E"},
		{"talking tree", "crab", "bat", "lobst", "\xB2"},
		{"beset", "plagued", "ravaged", "cursed", "scourged"},
		{"\x96 civil war", "\x9B \x98 \x99s", "a \x9B disease", "\x96 earthquakes", "\x96 solar activity"},
		{"its \x83 \x84", "the \xB1 \x98 \x99", "its inhabitants' \x9A \x85", "\xA1", "its \x8D \x8E"},
		{"juice", "brandy", "water", "brew", "gargle blasters"},
		{"\xB2", "\xB1 \x99", "\xB1 \xB2", "\xB1 \x9B", "\x9B \xB2"},
		{"fabulous", "exotic", "hoopy", "unusual", "exciting"},
		{"cuisine", "night life", "casinos", "sit coms", " \xA1 "},
		{"\xB0", "The planet \xB0", "The world \xB0", "This planet", "This world"},
		{"n unremarkable", " boring", " dull", " tedious", " revolting"},
		{"planet", "world", "place", "little planet", "dump"},
		{"wasp", "moth", "grub", "ant", "\xB2"},
		{"poet", "arts graduate", "yak", "snail", "slug"},
		{"tropical", "dense", "rain", "impenetrable", "exuberant"},
		{"funny", "wierd", "unusual", "strange", "peculiar"},
		{"frequent", "occasional", "unpredictable", "dreadful", "deadly"},
		{"\x82 \x81 for \x8A", "\x82 \x81 for \x8A and \x8A", "\x88 by \x89", "\x82 \x81 for \x8A but \x88 by \x89", "a\x90 \x91"},
		{"\x9B", "mountain", "edible", "tree", "spotted"},
		{"\x9F", "\xA0", "\x87oid", "\x93", "\x92"},
		{"ancient", "exceptional", "eccentric", "ingrained", "\x95"},
		{"killer", "deadly", "evil", "lethal", "vicious"},
		{"parking meters", "dust clouds", "ice bergs", "rock formations", "volcanoes"},
		{"plant", "tulip", "banana", "corn", "\xB2weed"},
		{"\xB2", "\xB1 \xB2", "\xB1 \x9B", "inhabitant", "\xB1 \xB2"},
		{"shrew", "beast", "bison", "snake", "wolf"},
		{"leopard", "cat", "monkey", "goat", "fish"},
		{"\x8C \x8B", "\xB1 \x9F \xA2", "its \x8D \xA0 \xA2", "\xA3 \xA4", "\x8C \x8B"},
		{"meat", "cutlet", "steak", "burgers", "soup"},
		{"ice", "mud", "Zero-G", "vacuum", "\xB1 ultra"},
		{"hockey", "cricket", "karate", "polo", "tennis"},
	}

	for {
		if len(source) == 0 {
			break
		}

		c := source[0:1]
		source = source[1:]

		cr := []byte(c)[0]

		if cr < 0x80 {
			fmt.Printf("%s", c)
		} else {
			if cr <= 0xA4 {
				rnd := g.prng.gen_rnd_number()

				a := 0
				b := 0
				c := 0
				d := 0

				if rnd >= 0x33 {
					a = 1
				}

				if rnd >= 0x66 {
					b = 1
				}

				if rnd >= 0x99 {
					c = 1
				}

				if rnd >= 0xCC {
					d = 1
				}

				g.goatSoup(desc[int(cr-0x81)][a+b+c+d], psy)

			} else {

				switch cr {
				case 0xB0: /* planet name */

					fmt.Printf("%s", psy.Name[0:1])
					for _, char := range psy.Name[1:] {
						fmt.Printf("%s", strings.ToLower(string(char)))
					}

				case 0xB1: /* <planet name>ian */
					i := 1
					fmt.Printf("%s", psy.Name[0:1])

					for _, char := range psy.Name[1:] {
						if (i+1 < len(psy.Name)) || ((char != 'E') && (char != 'I')) {
							fmt.Printf("%s", strings.ToLower(string(char)))
						}
					}
					fmt.Printf("ian")

				case 0xB2: /* random name */
					length := g.prng.gen_rnd_number() & 3

					for i := 0; uint16(i) <= length; i++ {
						x := g.prng.gen_rnd_number() & 0x3e
						if i == 0 {
							fmt.Printf("%c", g.dataTables.pairs0[x])
						} else {
							fmt.Printf("%s", strings.ToLower(string(g.dataTables.pairs0[x])))
						}

						fmt.Printf("%s", strings.ToLower(string(g.dataTables.pairs0[x+1])))
					}

				default:
					fmt.Printf("<bad char in data [%X]>", c)
					return
				}
			}
		}
	}
}

func (g *Galaxy) sgoatSoup(gs string, source string, psy *planetarySystem) string {

	desc := [][]string{
		{"fabled", "notable", "well known", "famous", "noted"},
		{"very", "mildly", "most", "reasonably", ""},
		{"ancient", "\x95", "great", "vast", "pink"},
		{"\x9E \x9D plantations", "mountains", "\x9C", "\x94 forests", "oceans"},
		{"shyness", "silliness", "mating traditions", "loathing of \x86", "love for \x86"},
		{"food blenders", "tourists", "poetry", "discos", "\x8E"},
		{"talking tree", "crab", "bat", "lobst", "\xB2"},
		{"beset", "plagued", "ravaged", "cursed", "scourged"},
		{"\x96 civil war", "\x9B \x98 \x99s", "a \x9B disease", "\x96 earthquakes", "\x96 solar activity"},
		{"its \x83 \x84", "the \xB1 \x98 \x99", "its inhabitants' \x9A \x85", "\xA1", "its \x8D \x8E"},
		{"juice", "brandy", "water", "brew", "gargle blasters"},
		{"\xB2", "\xB1 \x99", "\xB1 \xB2", "\xB1 \x9B", "\x9B \xB2"},
		{"fabulous", "exotic", "hoopy", "unusual", "exciting"},
		{"cuisine", "night life", "casinos", "sit coms", " \xA1 "},
		{"\xB0", "The planet \xB0", "The world \xB0", "This planet", "This world"},
		{"n unremarkable", " boring", " dull", " tedious", " revolting"},
		{"planet", "world", "place", "little planet", "dump"},
		{"wasp", "moth", "grub", "ant", "\xB2"},
		{"poet", "arts graduate", "yak", "snail", "slug"},
		{"tropical", "dense", "rain", "impenetrable", "exuberant"},
		{"funny", "wierd", "unusual", "strange", "peculiar"},
		{"frequent", "occasional", "unpredictable", "dreadful", "deadly"},
		{"\x82 \x81 for \x8A", "\x82 \x81 for \x8A and \x8A", "\x88 by \x89", "\x82 \x81 for \x8A but \x88 by \x89", "a\x90 \x91"},
		{"\x9B", "mountain", "edible", "tree", "spotted"},
		{"\x9F", "\xA0", "\x87oid", "\x93", "\x92"},
		{"ancient", "exceptional", "eccentric", "ingrained", "\x95"},
		{"killer", "deadly", "evil", "lethal", "vicious"},
		{"parking meters", "dust clouds", "ice bergs", "rock formations", "volcanoes"},
		{"plant", "tulip", "banana", "corn", "\xB2weed"},
		{"\xB2", "\xB1 \xB2", "\xB1 \x9B", "inhabitant", "\xB1 \xB2"},
		{"shrew", "beast", "bison", "snake", "wolf"},
		{"leopard", "cat", "monkey", "goat", "fish"},
		{"\x8C \x8B", "\xB1 \x9F \xA2", "its \x8D \xA0 \xA2", "\xA3 \xA4", "\x8C \x8B"},
		{"meat", "cutlet", "steak", "burgers", "soup"},
		{"ice", "mud", "Zero-G", "vacuum", "\xB1 ultra"},
		{"hockey", "cricket", "karate", "polo", "tennis"},
	}

	for {
		if len(source) == 0 {
			break
		}

		c := source[0:1]
		source = source[1:]

		cr := []byte(c)[0]

		if cr < 0x80 {
			gs = gs + c
		} else {
			if cr <= 0xA4 {
				rnd := g.prng.gen_rnd_number()

				a := 0
				b := 0
				c := 0
				d := 0

				if rnd >= 0x33 {
					a = 1
				}

				if rnd >= 0x66 {
					b = 1
				}

				if rnd >= 0x99 {
					c = 1
				}

				if rnd >= 0xCC {
					d = 1
				}

				gs = g.sgoatSoup(gs, desc[int(cr-0x81)][a+b+c+d], psy)

			} else {

				switch cr {
				case 0xB0: /* planet name */

					gs = gs + psy.Name[0:1]
					for _, char := range psy.Name[1:] {
						gs = gs + strings.ToLower(string(char))
					}

				case 0xB1: /* <planet name>ian */
					i := 1
					gs = gs + psy.Name[0:1]

					for _, char := range psy.Name[1:] {
						if (i+1 < len(psy.Name)) || ((char != 'E') && (char != 'I')) {
							gs = gs + strings.ToLower(string(char))
						}
					}

					gs = gs + "ian"

				case 0xB2: /* random name */
					length := g.prng.gen_rnd_number() & 3

					for i := 0; uint16(i) <= length; i++ {
						x := g.prng.gen_rnd_number() & 0x3e
						if i == 0 {
							gs = gs + fmt.Sprintf("%c", g.dataTables.pairs0[x])
						} else {
							gs = gs + strings.ToLower(string(g.dataTables.pairs0[x]))
						}

						gs = gs + strings.ToLower(string(g.dataTables.pairs0[x+1]))
					}

				default:
					gs = fmt.Sprintf("<bad char in data [%X]>", c)
					return gs
				}
			}
		}
	}

	return gs
}

// Convienince function to return the name of the planet that the player is at
func (g *Game) PlayerCurrentPlanetName() string {
	return g.Galaxy.Systems[g.Player.Ship.Location.CurrentPlanet].Name
}

func (g *Galaxy) PrintSystem(plsy planetarySystem, compressed bool) {

	if compressed {
		fmt.Printf("%10s\t", plsy.Name)
		fmt.Printf(" TL: %2d ", (plsy.Techlev)+1)
		fmt.Printf("%12s", g.dataTables.econnames[plsy.Economy])
		fmt.Printf(" %15s", g.dataTables.govnames[plsy.Govtype])
	} else {
		fmt.Printf("\n\nSystem:  ")
		fmt.Printf(plsy.Name)
		fmt.Printf("\nPosition (%d,", plsy.X)
		fmt.Printf("%d)", plsy.Y)
		fmt.Printf("\nEconomy: (%d) ", plsy.Economy)
		fmt.Printf(g.dataTables.econnames[plsy.Economy])
		fmt.Printf("\nGovernment: (%d) ", plsy.Govtype)
		fmt.Printf(g.dataTables.govnames[plsy.Govtype])
		fmt.Printf("\nTech Level: %2d", (plsy.Techlev)+1)
		fmt.Printf("\nTurnover: %d", (plsy.Productivity))
		fmt.Printf("\nRadius: %d", plsy.Radius)
		fmt.Printf("\nPopulation: %d Billion", (plsy.Population)>>3)

		g.prng.rnd_seed = plsy.goatsoupseed
		fmt.Println()
		g.goatSoup("\x8F is \x97.", &plsy)
		fmt.Println()

	}
}

func (g *Galaxy) SprintSystem(plsy planetarySystem, compressed bool) string {

	systemData := ""
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	g.prng.rnd_seed = plsy.goatsoupseed
	gs := g.sgoatSoup("", "\x8F is \x97.", &plsy)

	if compressed {
		systemData = fmt.Sprintf("%10s\t", plsy.Name)
		systemData = systemData + fmt.Sprintf(" %s %2d ", style.Render(" TL:"), (plsy.Techlev)+1)
		systemData = systemData + fmt.Sprintf("%12s", g.dataTables.econnames[plsy.Economy])
		systemData = systemData + fmt.Sprintf(" %15s", g.dataTables.govnames[plsy.Govtype])
	} else {
		systemData = systemData + fmt.Sprintf("%s %s\n", style.Render("System:"), plsy.Name)
		systemData = systemData + fmt.Sprintf("%s (%d,", style.Render("Position:"), plsy.X)
		systemData = systemData + fmt.Sprintf("%d)\n", plsy.Y)
		systemData = systemData + fmt.Sprintf("%s (%d) ", style.Render("Economy:"), plsy.Economy)
		systemData = systemData + fmt.Sprintf("%s\n", g.dataTables.econnames[plsy.Economy])
		systemData = systemData + fmt.Sprintf("%s (%d) ", style.Render("Government"), plsy.Govtype)
		systemData = systemData + fmt.Sprintf("%s\n", g.dataTables.govnames[plsy.Govtype])
		systemData = systemData + fmt.Sprintf("%s %2d\n", style.Render("Tech Level:"), (plsy.Techlev)+1)
		systemData = systemData + fmt.Sprintf("%s %d\n", style.Render("Turnover:"), (plsy.Productivity))
		systemData = systemData + fmt.Sprintf("%s %d\n", style.Render("Radius:"), plsy.Radius)
		systemData = systemData + fmt.Sprintf("%s %d Billion\n", style.Render("Population:"), (plsy.Population)>>3)

		systemData = systemData + gs + "\n"
		//systemData = systemData + "\n"

	}

	return systemData
}
