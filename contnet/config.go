package contnet

type pkgConfig struct {
	ItemsPerPage int32   // number of results per page
	Gravity      float64 // gravity is a constant used to drag old content down and obliterate it from the network
	Novelty      float64 // novelty factor is a constant used to determine the amount of new content to recommend
}

var config *pkgConfig

func init() {
	config = &pkgConfig{
		ItemsPerPage: 30,
		Gravity:      1.0,
		Novelty:      0.3,
	}
}

func Configure(itemsPerPage int32, gravity, novelty float64) {
	config = &pkgConfig{
		ItemsPerPage: itemsPerPage,
		Gravity:      gravity,
		Novelty:      novelty,
	}
}
