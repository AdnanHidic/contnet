package contnet

import (
	"github.com/asaskevich/EventBus"
	"sync"
)

type Network struct {
	sync.RWMutex
	bus          *EventBus.EventBus
	index        *Index
	contentStore *ContentStore
	profileStore *ProfileStore
}
type NetworkFactory struct{}

func (factory NetworkFactory) New() *Network {
	bus := EventBus.New()
	contentStore := Object.ContentStore.New(bus)
	return &Network{
		contentStore: contentStore,
		profileStore: Object.ProfileStore.New(bus),
		index:        Object.Index.New(bus, contentStore),
	}
}

func (net *Network) Save() error {
	return Errors.NotImplemented
}

func (net *Network) Load() error {
	return Errors.NotImplemented
}

// Attempts to update network object with content specified.
// If content did not exist, it is added to the network.
// If content did exist, it is updated.
func (net *Network) SaveContent(content *Content) error {
	net.contentStore.Upsert(content)
	return nil
}

// Attempts to update network object with action for profile specified.
func (net *Network) SaveAction(action *Action) error {
	// get related content's copy, if any
	relatedContent := net.contentStore.Get(action.ContentID)

	// if content does not exist
	if relatedContent == nil {
		return Errors.ContentNotFound
	}

	// injected related content to action
	action.Content = relatedContent

	// save action to profile
	net.profileStore.SaveAction(action)

	// everything is okay
	return nil
}

func (net *Network) Select(profileID int64, page uint8) ([]*Content, error) {
	return nil, Errors.NotImplemented
}

func (net *Network) Describe() *NetworkDescription {
	return &NetworkDescription{
		Contents: len(net.contentStore.contents),
		Profiles: len(net.profileStore.profiles),
	}
}
