package contnet

import "time"

type NetworkConfig struct {
	BackupPath     string // full path to backup file
	ItemsPerPage   int32
	Gravity        float64       // gravity is a constant used to drag old content down and obliterate it from the network
	Novelty        float64       // novelty factor is a constant used to determine the amount of new content to recommend
	AutosavePeriod time.Duration // number of seconds until network autosave is invoked
}
type NetworkConfigFactory struct{}

func (factory NetworkConfigFactory) New(backupPath string, itemsPerPage int32, gravity, novelty float64, autosavePeriod time.Duration) *NetworkConfig {
	return &NetworkConfig{
		BackupPath:     backupPath,
		ItemsPerPage:   itemsPerPage,
		Gravity:        gravity,
		Novelty:        novelty,
		AutosavePeriod: autosavePeriod,
	}
}

func (config *NetworkConfig) Clone() *NetworkConfig {
	return Object.NetworkConfig.New(
		config.BackupPath,
		config.ItemsPerPage,
		config.Gravity,
		config.Novelty,
		config.AutosavePeriod,
	)
}
