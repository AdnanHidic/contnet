package contnet

import (
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
)

var __gravitySleep, _ = time.ParseDuration("60s")

type ContentStore struct {
	sync.RWMutex
	bus      *EventBus.EventBus
	contents map[ID]*Content
}
type ContentStoreFactory struct{}

func (factory ContentStoreFactory) New(bus *EventBus.EventBus) *ContentStore {
	store := &ContentStore{
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
	store.contents[content.ID] = content.Clone()

	if existed {
		store.bus.Publish("content:reindex", old, content)
	} else {
		store.bus.Publish("content:index", content)
	}

}

func (store *ContentStore) delete(content *Content) {
	delete(store.contents, content.ID)
	store.bus.Publish("content:removed", content)
}

func (store *ContentStore) __gravity() {
	var gravity float64
	var contentsToRemove []*Content
	for {
		// lock for reading, we just want to calculate values and select candidates for deletion
		store.RLock()

		contentsToRemove = []*Content{}

		for _, content := range store.contents {
			gravity = __applyGravity(content)
			if gravity > GRAVITY_TRESHOLD {
				contentsToRemove = append(contentsToRemove, content)
			}
		}
		// reading part is over, switch to full write lock to perform deletions
		store.RUnlock()
		store.Lock()
		for i := 0; i < len(contentsToRemove); i++ {
			store.delete(contentsToRemove[i])
		}
		store.Unlock()

		time.Sleep(__gravitySleep)
	}
}

const GRAVITY_TRESHOLD = 100

func __applyGravity(content *Content) float64 {
	return 0
}
