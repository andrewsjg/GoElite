package main

import (
	"strings"
	"testing"
)

func TestGalaxy(t *testing.T) {
	galaxy := initGalaxy(1)

	if strings.Compare(strings.ToUpper(galaxy.systems[7].name), "LAVE") == 0 {
		t.Error("System number 7 in galaxy 1 should be LAVE. Got:", galaxy.systems[7].name)
	}

	if strings.Compare(strings.ToUpper(galaxy.systems[147].name), "DISO") == 0 {
		t.Error("System number 147 in galaxy 1 should be DISO. Got:", galaxy.systems[147].name)
	}

	lave := galaxy.systems[7]

	if lave.economy != 5 {
		t.Error("LAVE has the wrong economy type. Got: ", lave.economy)
	}

	if lave.govtype != 3 {
		t.Error("LAVE has the wrong government type. Got: ", lave.govtype)
	}

	if lave.techlev != 4 {
		t.Error("LAVE has the wrong tech level. Got: ", lave.techlev)
	}

	if lave.productivity != 7000 {
		t.Error("LAVE has the wrong productivity. Got: ", lave.productivity)
	}

	if lave.population != 25 {
		t.Error("LAVE has the wrong population. Got: ", lave.population)
	}

	if lave.x != 20 {
		t.Error("LAVE has the wrong X position. Got: ", lave.x)
	}

	if lave.y != 173 {
		t.Error("LAVE has the wrong Y position. Got: ", lave.y)
	}

	if lave.radius != 4116 {
		t.Error("LAVE has the wrong radius. Got: ", lave.radius)
	}
}
