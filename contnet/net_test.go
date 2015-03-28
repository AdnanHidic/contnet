package contnet

import (
	"testing"
	"time"
)

var (
	backupPath        = "test"
	itemsPerPage      = int32(10)
	gravity           = 1.0
	novelty           = 0.3
	autosavePeriod, _ = time.ParseDuration("10s")
)

func TestNetworkCreate(t *testing.T) {
	netConfig := Object.NetworkConfig.New(
		backupPath,
		itemsPerPage,
		gravity,
		novelty,
		autosavePeriod,
	)
	net := Object.Network.New(netConfig)

	description := net.Describe()
	if description.Contents != 0 {
		t.Errorf("Got %d but expected %d contents", description.Contents, 0)
	}
	if description.Profiles != 0 {
		t.Errorf("Got %d but expected %d profiles", description.Profiles, 0)
	}

}
