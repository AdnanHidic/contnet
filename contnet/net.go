package contnet

import (
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

type NetConfig struct {
	MaxContentAge           time.Duration
	CheckContentAgeInterval time.Duration
	ItemsPerPage            uint8
	NoveltyStrength         float64
	SnapshotPath            string
	SnapshotInterval        time.Duration
}

// TODO: implement trend store
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
	contentStore := Object.ContentStore.New(config, bus)
	net := &Net{
		config:       config,
		bus:          bus,
		contentStore: contentStore,
		profileStore: Object.ProfileStore.New(),
		index:        Object.Index.New(config, bus, contentStore),
	}
	if err := net.Restore(); err != nil {
		log.Print("Failed to restore net object, proceeding empty..")
	}
	go net.__snapshot()
    go net.index.__refresh()

	return net
}

func (net *Net) __snapshot() {
	for {
		net.Snapshot()
		time.Sleep(net.config.SnapshotInterval)
	}
}

func (net *Net) Snapshot() error {
	if err := net.contentStore.Snapshot(net.config.SnapshotPath, "content"); err != nil {
		return err
	}

	if err := net.index.Snapshot(net.config.SnapshotPath, "index"); err != nil {
		return err
	}

	if err := net.profileStore.Snapshot(net.config.SnapshotPath, "profiles"); err != nil {
		return err
	}

	// TODO: snapshot trends

	return nil
}

func (net *Net) Restore() error {
	if err := net.contentStore.RestoreFromSnapshot(net.config.SnapshotPath, "content"); err != nil {
		return err
	}

	if err := net.index.RestoreFromSnapshot(net.config.SnapshotPath, "index"); err != nil {
		return err
	}

	if err := net.profileStore.RestoreFromSnapshot(net.config.SnapshotPath, "profiles"); err != nil {
		return err
	}

	// TODO: restore trends

	return nil
}

// Attempts to update network object with content specified.
// If content did not exist, it is added to the network.
// If content did exist, it is updated.
func (net *Net) SaveContent(content *Content) {
	net.contentStore.Upsert(content)
}

// Attempts to update network object with action for profile specified.
func (net *Net) SaveAction(action *Action) error {
	// get related content's copy, if any
	relatedContent := net.contentStore.Get(action.ContentID)

	// if content does not exist
	if relatedContent == nil {
		return Errors.ContentNotFound
	}

	// inject related content to action
	action.Content = relatedContent

	// save action to profile
	net.profileStore.Save(action)

	// everything is okay
	return nil
}

func (net *Net) Select(profileID int64, page uint8) ([]*Content, error) {
	// TODO implement selector
	return nil, Errors.NotImplemented
}

func (net *Net) Describe() *NetDescription {
	return &NetDescription{
		Contents: net.contentStore.Describe(),
		Index:    net.index.Describe(),
        Profiles: net.profileStore.Describe(),
	}
}
