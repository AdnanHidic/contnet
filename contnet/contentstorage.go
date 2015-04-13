package contnet

import (
	"sync"
	"time"
    "github.com/asaskevich/EventBus"
)

var __gravitySleep, _ = time.ParseDuration("60s")

type ContentStore struct {
	sync.RWMutex
    bus *EventBus.EventBus
	contents map[ID]*Content
}
type ContentStoreFactory struct{}

func (factory ContentStoreFactory) New(bus *EventBus.EventBus) *ContentStore {
	store := &ContentStore{
        bus: bus,
		contents: map[ID]*Content{},
	}
	go store.__gravity()
	return store
}

func (store *ContentStore) Snapshot(path, filename string) {

}

func (store *ContentStore) RestoreFromSnapshot(path, filename string) {

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
	store.contents[content.ID] = content.Clone()

    store.bus.Publish("content:reindex", content)
}

func (store *ContentStore) delete(content *Content) {
    delete(store.contents, content.ID)
    store.bus.Publish("content:remove", content)
}

func (store *ContentStore) __gravity() {
	var gravity float64
    for {
		store.Lock()

		for _, content := range store.contents {
            gravity = __applyGravity(content)
            if gravity > GRAVITY_TRESHOLD {
                store.delete(content)
            }
		}

		store.Unlock()
		time.Sleep(__gravitySleep)
	}
}

const GRAVITY_TRESHOLD = 100

func __applyGravity(content *Content) float64 {
    return 0
}
