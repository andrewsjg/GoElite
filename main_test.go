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

}
