package contnet

import (
	"errors"
	"sync"
	"time"
)

var Errors = struct {
	NotImplemented error
}{
	NotImplemented: errors.New("Not implemented."),
}

type NetworkConfig struct {
	BackupPath     string // full path to backup file
	ItemsPerPage   int32
	Gravity        float64       // gravity is a constant used to drag old content down and obliterate it from the network
	Novelty        float64       // novelty factor is a constant used to determine the amount of new content to recommend
	AutosavePeriod time.Duration // number of seconds until network autosave is invoked
}

type Network struct {
	sync.RWMutex
	Config         *NetworkConfig
	ContentStorage *ContentStorage
	ProfileStorage *ProfileStorage
}

func NewNetwork(config *NetworkConfig) *Network {
	return &Network{
		Config:         config,
		ContentStorage: Object.ContentStorage.New(),
		ProfileStorage: Object.ProfileStorage.New(),
	}
}

func (n *Network) Init() error {
	return Errors.NotImplemented
}

func (n *Network) SaveContent(c *Content) error {
	return Errors.NotImplemented
}

func (n *Network) SaveAction(a *Action) error {
	return Errors.NotImplemented
}

func (n *Network) Select(consumerID, page int8) ([]*Content, error) {
	return nil, Errors.NotImplemented
}
