package main

import (
	"testing"
)

func TestGalaxy(t *testing.T) {
	galaxy := initGalaxy(1)

	if galaxy.systems[7].name != "LAVE" {
		t.Error("System number 7 in galaxy 1 should be LAVE. Got: ", galaxy.systems[7].name)
	}

	if galaxy.systems[147].name != "DISO" {
		t.Error("System number 7 in galaxy 1 should be DISO. Got: ", galaxy.systems[147].name)
	}
}
