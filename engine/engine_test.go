package eliteEngine

import (
	"strings"
	"testing"
)

func TestGalaxy(t *testing.T) {
	galaxy := InitGalaxy(1)

	if strings.Compare(strings.ToUpper(galaxy.Systems[7].Name), "LAVE") == 0 {
		t.Error("System number 7 in galaxy 1 should be LAVE. Got:", galaxy.Systems[7].Name)
	}

	if strings.Compare(strings.ToUpper(galaxy.Systems[147].Name), "DISO") == 0 {
		t.Error("System number 147 in galaxy 1 should be DISO. Got:", galaxy.Systems[147].Name)
	}

	lave := galaxy.Systems[7]

	if lave.Economy != 5 {
		t.Error("LAVE has the wrong economy type. Got: ", lave.Economy)
	}

	if lave.Govtype != 3 {
		t.Error("LAVE has the wrong government type. Got: ", lave.Govtype)
	}

	if lave.Techlev != 4 {
		t.Error("LAVE has the wrong tech level. Got: ", lave.Techlev)
	}

	if lave.Productivity != 7000 {
		t.Error("LAVE has the wrong productivity. Got: ", lave.Productivity)
	}

	if lave.Population != 25 {
		t.Error("LAVE has the wrong population. Got: ", lave.Population)
	}

	if lave.X != 20 {
		t.Error("LAVE has the wrong X position. Got: ", lave.X)
	}

	if lave.Y != 173 {
		t.Error("LAVE has the wrong Y position. Got: ", lave.Y)
	}

	if lave.Radius != 4116 {
		t.Error("LAVE has the wrong radius. Got: ", lave.Radius)
	}
}
