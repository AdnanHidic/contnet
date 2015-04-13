package contnet

import (
	"github.com/asaskevich/EventBus"
	"sync"
)

type NetConfig struct {
	ItemsPerPage    uint8
	GravityStrength float64
	NoveltyStrength float64
	SnapshotPath    string
}

type Net struct {
	sync.RWMutex
	config       *NetConfig
	bus          *EventBus.EventBus
	contentStore *ContentStore
	profileStore *ProfileStore
	index        *Index
}
type NetFactory struct{}

func (factory NetFactory) New(config *NetConfig) *Net {
	bus := EventBus.New()
	contentStore := Object.ContentStore.New(bus)
	return &Net{
		config:       config,
		bus:          bus,
		contentStore: contentStore,
		profileStore: Object.ProfileStore.New(),
		index:        Object.Index.New(bus, contentStore),
	}
}

func (net *Net) Snapshot() error {
	err := net.contentStore.Snapshot(net.config.SnapshotPath, "content")
	if err != nil {
		return err
	}

	return nil
}

func (net *Net) Restore() error {
	err := net.contentStore.RestoreFromSnapshot(net.config.SnapshotPath, "content")
	if err != nil {
		return err
	}

	return nil
}

// Attempts to update network object with content specified.
// If content did not exist, it is added to the network.
// If content did exist, it is updated.
func (net *Net) SaveContent(content *Content) error {
	net.contentStore.Upsert(content)
	return nil
}

// Attempts to update network object with action for profile specified.
func (net *Net) SaveAction(action *Action) error {
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

func (net *Net) Select(profileID int64, page uint8) ([]*Content, error) {
	return nil, Errors.NotImplemented
}

func (net *Net) Describe() *NetDescription {
	return &NetDescription{
		Contents: len(net.contentStore.contents),
		Profiles: len(net.profileStore.profiles),
	}
}
