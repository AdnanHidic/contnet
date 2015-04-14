package contnet

import (
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

type ContentConfig struct {
	GravityStrength         float64
	MaxContentAge           time.Duration
	CheckContentAgeInterval time.Duration
}

type ContentStore struct {
	sync.RWMutex
	config   *ContentConfig
	bus      *EventBus.EventBus
	contents map[ID]*Content
}
type ContentStoreFactory struct{}

func (factory ContentStoreFactory) New(config *ContentConfig, bus *EventBus.EventBus) *ContentStore {
	store := &ContentStore{
		config:   config,
		bus:      bus,
		contents: map[ID]*Content{},
	}
	go store.__gravity()
	return store
}

func (store *ContentStore) Snapshot(path, filename string) error {
	store.RLock()
	defer store.RUnlock()

	return __snapshot(path, filename, store.contents)
}

func (store *ContentStore) RestoreFromSnapshot(path, filename string) error {
	store.Lock()
	defer store.Unlock()

	object, err := __restoreFromSnapshot(path, filename, store.contents)

	if err == nil {
		store.contents = object.(map[ID]*Content)
	}

	return err
}

func (store *ContentStore) Get(id ID) *Content {
	store.RLock()
	defer store.RUnlock()

	if content, ok := store.contents[id]; !ok {
		return nil
	} else {
		return content.Clone()
	}
}

func (store *ContentStore) GetMany(ids []ID) []*Content {
	store.RLock()
	defer store.RUnlock()

	out := []*Content{}

	for i := 0; i < len(ids); i++ {
		if content, ok := store.contents[ids[i]]; !ok {
			out = append(out, nil)
		} else {
			out = append(out, content.Clone())
		}
	}

	return out
}

func (store *ContentStore) Upsert(content *Content) {
	store.Lock()
	defer store.Unlock()

	// save it to map of contents
	old, existed := store.contents[content.ID]
    new := content.Clone()
    new.age = __age(time.Now(), *new, store.config.GravityStrength)
	store.contents[content.ID] = new

	if existed {
		store.bus.Publish("content:reindex", old, new)
	} else {
		store.bus.Publish("content:index", new)
	}

}

func (store *ContentStore) delete(content *Content) {
	delete(store.contents, content.ID)
	store.bus.Publish("content:removed", content)
}

func (store *ContentStore) __gravity() {
	referenceTime := time.Now()

	var contentsToRemove []*Content
	for {
		// lock for reading, we just want to calculate values and select candidates for deletion
		store.RLock()

		contentsToRemove = []*Content{}

		// for each content stored
		for _, content := range store.contents {
			// calculate age based on content parameters
			content.age = __age(referenceTime, *content, store.config.GravityStrength)
			// if content is considered stale and old, mark it for deletion
			if content.age.Before(referenceTime.Add(-store.config.MaxContentAge)) {
				contentsToRemove = append(contentsToRemove, content)
			}
		}
		// reading part is over, switch to full write lock to perform deletions
		store.RUnlock()
		store.Lock()
		log.Printf("GLOBAL PAUSE: Cleaning stale content from content store. Total: %d marked for deletion", len(contentsToRemove))
		for i := 0; i < len(contentsToRemove); i++ {
			store.delete(contentsToRemove[i])
		}
		log.Print("GLOBAL RESUME: Stale content cleaned from content store.")
		store.Unlock()

		time.Sleep(store.config.CheckContentAgeInterval)
	}
}
