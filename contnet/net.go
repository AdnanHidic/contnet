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

func (net *Network) Init() error {
	return Errors.NotImplemented
}

func (net *Network) SaveContent(c *Content) error {
	return Errors.NotImplemented
}

func (net *Network) SaveAction(action *Action) error {
	// get related content's copy, if any
	relatedContent := net.ContentStorage.Get(action.ContentID)

	// if content does not exist
	if relatedContent == nil {
		return Errors.ContentNotFound
	}

	// injected related content to action
	action.Content = relatedContent

	// save action to profile
	net.ProfileStorage.SaveAction(action)

	// everything is okay
	return nil
}

func (net *Network) Select(consumerID, page int8) ([]*Content, error) {
	return nil, Errors.NotImplemented
}
