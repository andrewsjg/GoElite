package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// Go implementation  txtelite

// four byte random number used for planet description
type fastseed struct {
	a uint
	b uint
	c uint
	d uint
}

// six byte random number used as seed for planets
type seed struct {
	w0 uint16
	w1 uint16
	w2 uint16
}

type planetarySystem struct {
	x            uint16
	y            uint16 /* One byte unsigned */
	economy      uint16 /* These two are actually only 0-7  */
	govtype      uint16
	techlev      uint16 /* 0-16 i think */
	population   uint16 /* One byte */
	productivity uint16 /* Two byte */
	radius       uint16 /* Two byte (not used by game at all) */
	goatsoupseed fastseed
	name         string
}

var mainseed seed
var rnd_seed fastseed

const galSize = 256

var galaxy = [galSize]planetarySystem{}

const base0 = 0x5A4A
const base1 = 0x0248
const base2 = 0xB753 /* Base seed for galaxy 1 */

var pairs0 = []byte("ABOUSEITILETSTONLONUTHNOALLEXEGEZACEBISOUSESARMAINDIREA.ERATENBERALAVETIEDORQUANTEISRION")
var pairs = []byte("..LEXEGEZACEBISO" +
	"USESARMAINDIREA." +
	"ERATENBERALAVETI" +
	"EDORQUANTEISRION") /* Dots should be nullprint characters */

var govnames = []string{"Anarchy", "Feudal", "Multi-gov", "Dictatorship",
	"Communist", "Confederacy", "Democracy", "Corporate State"}

var econnames = []string{"Rich Ind", "Average Ind", "Poor Ind", "Mainly Ind",
	"Mainly Agri", "Rich Agri", "Average Agri", "Poor Agri"}

var unitnames = []string{"t", "kg", "g"}

var lastrand int64 = 0

var galaxyNum uint = 1

var currentPlanet int

func gen_rnd_number() uint {
	var a, x uint
	x = (rnd_seed.a * 2) & 0xFF
	a = x + rnd_seed.c

	if rnd_seed.a > 127 {
		a++
	}

	rnd_seed.a = a & 0xFF
	rnd_seed.c = x

	a = a / 256 /* a = any carry left from above */
	x = rnd_seed.b
	a = (a + x + rnd_seed.d) & 0xFF
	rnd_seed.b = a
	rnd_seed.d = x

	return a
}

func mysrand(seed int64) {
	rand.Seed(seed)

	lastrand = seed - 1
}

func tweakseed(s *seed) {

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
func nextgalaxy(s *seed) {
	s.w0 = twist(s.w0)
	s.w1 = twist(s.w1)
	s.w2 = twist(s.w2)
}

func buildGalaxy(galaxyNum int) {
	var syscount, galcount int

	mainseed.w0 = base0
	mainseed.w1 = base1
	mainseed.w2 = base2 /* Initialise seed for galaxy 1 */

	for galcount = 1; galcount < galaxyNum; galcount++ {
		nextgalaxy(&mainseed)
	}

	/* Put galaxy data into array of structures */
	for syscount = 0; syscount < galSize; syscount++ {
		galaxy[syscount] = makeSystem(&mainseed)
	}
}

func makeSystem(s *seed) planetarySystem {

	var pair1, pair2, pair3, pair4 uint16
	var longnameflag uint16 = (s.w0) & 64

	planSys := planetarySystem{}

	planSys.x = ((s.w1) >> 8)
	planSys.y = ((s.w0) >> 8)

	planSys.govtype = (((s.w1) >> 3) & 7) /* bits 3,4 &5 of w1 */

	planSys.economy = (((s.w0) >> 8) & 7) /* bits 8,9 &A of w0 */
	if planSys.govtype <= 1 {
		planSys.economy = ((planSys.economy) | 2)
	}

	planSys.techlev = (((s.w1) >> 8) & 3) + ((planSys.economy) ^ 7)
	planSys.techlev += ((planSys.govtype) >> 1)

	if ((planSys.govtype) & 1) == 1 {
		planSys.techlev += 1
	}
	/* C simulation of 6502's LSR then ADC */

	planSys.population = 4*(planSys.techlev) + (planSys.economy)
	planSys.population += (planSys.govtype) + 1

	planSys.productivity = (((planSys.economy) ^ 7) + 3) * ((planSys.govtype) + 4)
	planSys.productivity *= (planSys.population) * 8

	planSys.radius = 256*((((s.w2)>>8)&15)+11) + planSys.x

	planSys.goatsoupseed.a = uint(s.w1 & 0xFF)
	planSys.goatsoupseed.b = uint(s.w1 >> 8)
	planSys.goatsoupseed.c = uint(s.w2 & 0xFF)
	planSys.goatsoupseed.d = uint(s.w2 >> 8)

	pair1 = 2 * (((s.w2) >> 8) & 31)
	tweakseed(s)

	pair2 = 2 * (((s.w2) >> 8) & 31)
	tweakseed(s)

	pair3 = 2 * (((s.w2) >> 8) & 31)
	tweakseed(s)

	pair4 = 2 * (((s.w2) >> 8) & 31)
	tweakseed(s)
	/* Always four iterations of random number */

	name := make([]byte, 8)

	name[0] = pairs[pair1]
	name[1] = pairs[pair1+1]
	name[2] = pairs[pair2]
	name[3] = pairs[pair2+1]
	name[4] = pairs[pair3]
	name[5] = pairs[pair3+1]

	/* bit 6 of ORIGINAL w0 flags a four-pair name */
	if longnameflag == 1 {
		name[6] = pairs[pair4]
		name[7] = pairs[pair4+1]
		name[8] = 0
	} else {
		name[6] = 0
	}

	planSys.name = strings.ReplaceAll(string(name), ".", "")

	return planSys
}

/* Return id of the planet whose name matches passed strinmg
   closest to currentplanet - if none return currentplanet */

func matchsys(platnetName string) int {

	p := currentPlanet
	d := 9999

	for syscount := 0; syscount < galSize; syscount++ {
		if strings.HasPrefix(galaxy[syscount].name, platnetName) {
			if distance(galaxy[syscount], galaxy[currentPlanet]) < d {
				d = distance(galaxy[syscount], galaxy[currentPlanet])
				p = syscount
			}
		}
	}
	return p
}

// Seperation between two planets (4*sqrt(X*X+Y*Y/4))
func distance(planetA planetarySystem, planetB planetarySystem) int {

	val := (float64((planetA.x-planetB.x)*(planetA.x-planetB.x) + (planetA.y-planetB.y)*(planetA.y-planetB.y))) / 4.0

	return int(4 * math.Sqrt(val))
}

func goatSoup(source string, psy *planetarySystem) {

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
				rnd := gen_rnd_number()

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

				goatSoup(desc[int(cr-0x81)][a+b+c+d], psy)

			} else {

				switch cr {
				case 0xB0: /* planet name */

					fmt.Printf("%s", psy.name[0:1])
					for _, char := range psy.name[1:] {
						fmt.Printf("%s", strings.ToLower(string(char)))
					}
					break
				case 0xB1: /* <planet name>ian */
					i := 1
					fmt.Printf("%s", psy.name[0:1])

					for _, char := range psy.name[1:] {
						if (i+1 < len(psy.name)) || ((char != 'E') && (char != 'I')) {
							fmt.Printf("%s", strings.ToLower(string(char)))
						}
					}
					fmt.Printf("ian")
					break
				case 0xB2: /* random name */
					length := gen_rnd_number() & 3

					for i := 0; uint(i) <= length; i++ {
						x := gen_rnd_number() & 0x3e
						if i == 0 {
							fmt.Printf("%c", pairs0[x])
						} else {
							fmt.Printf("%s", strings.ToLower(string(pairs0[x])))
						}

						fmt.Printf("%s", strings.ToLower(string(pairs0[x+1])))
					}
					break
				default:
					fmt.Printf("<bad char in data [%X]>", c)
					return
				}
			}
		}
	}
}

func printSystem(plsy planetarySystem, compressed bool) {

	if compressed {
		fmt.Printf("%10s", plsy.name)
		fmt.Printf(" TL: %2d ", (plsy.techlev)+1)
		fmt.Printf("%12s", econnames[plsy.economy])
		fmt.Printf(" %15s", govnames[plsy.govtype])
	} else {
		fmt.Printf("\n\nSystem:  ")
		fmt.Printf(plsy.name)
		fmt.Printf("\nPosition (%d,", plsy.x)
		fmt.Printf("%d)", plsy.y)
		fmt.Printf("\nEconomy: (%d) ", plsy.economy)
		fmt.Printf(econnames[plsy.economy])
		fmt.Printf("\nGovernment: (%d) ", plsy.govtype)
		fmt.Printf(govnames[plsy.govtype])
		fmt.Printf("\nTech Level: %2d", (plsy.techlev)+1)
		fmt.Printf("\nTurnover: %d", (plsy.productivity))
		fmt.Printf("\nRadius: %d", plsy.radius)
		fmt.Printf("\nPopulation: %d Billion", (plsy.population)>>3)

		rnd_seed = plsy.goatsoupseed
		fmt.Println()
		goatSoup("\x8F is \x97.", &plsy)
		fmt.Println()

	}
}

func main() {
	// Init things
	mysrand(12345)

	galaxynum := 1
	buildGalaxy(galaxynum)
	currentPlanet = 7 // Lave
	debugTests()

}

func debugTests() {

	fmt.Printf("The current system is: %s", galaxy[currentPlanet].name)
	// test current planet (Lave at start)
	printSystem(galaxy[currentPlanet], false)
	fmt.Println()
	// test matchsys
	fmt.Printf("DISO is system numner: %d", matchsys("DISO")) // 147

	diso := matchsys("DISO")
	printSystem(galaxy[diso], false)
}
