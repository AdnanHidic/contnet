package contnet

import "sync"

type Network struct {
	sync.RWMutex
	Config         *NetworkConfig
	ContentStorage *ContentStorage
	ProfileStorage *ProfileStorage
}
type NetworkFactory struct{}

func (factory NetworkFactory) New(config *NetworkConfig) *Network {
	return &Network{
		Config:         config.Clone(),
		ContentStorage: Object.ContentStorage.New(),
		ProfileStorage: Object.ProfileStorage.New(),
	}
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
