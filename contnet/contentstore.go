package contnet

import (
	"bufio"
	"encoding/json"
	"github.com/asaskevich/EventBus"
	"log"
	"os"
	"sync"
	"time"
)

var __gravitySleep, _ = time.ParseDuration("60s")

const (
	__errSnapshotFile = "Failed to take a content store snapshot because snapshot file could not be created."
	__errSnapshotJson = "Failed to take a content store snapshot because store object failed to serialize."
	__snapshotSaved   = "Content store snapshot successfully created."

	__errRestoreFile   = "Failed to restore content store from snapshot file because file could not be opened."
	__errRestoreJson   = "Failed to restore content store from snapshot file because its JSON content failed to deserialize."
	__snapshotRestored = "Content store snapshot successfully loaded."
)

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

func __fullPath(path, filename string) string {
	return path + "/" + filename
}

func (store *ContentStore) Snapshot(path, filename string) error {
	// create new snapshot file
	fullpath := __fullPath(path, filename)
	file, err := os.Create(fullpath)
	if err != nil {
		log.Print(__errSnapshotFile, err.Error())
		return err
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	store.RLock()
	defer store.RUnlock()

	err = json.NewEncoder(bufferedWriter).Encode(store.contents)
	if err != nil {
		log.Print(__errSnapshotJson, err.Error())
		return err
	}

	log.Print(__snapshotSaved)
    return nil
}

func (store *ContentStore) RestoreFromSnapshot(path, filename string) error {
	fullpath := __fullPath(path, filename)
	file, err := os.Open(fullpath)
	if err != nil {
		log.Print(__errRestoreFile, err.Error())
		return err
	}
	defer file.Close()

	bufferedReader := bufio.NewReader(file)
	store.Lock()
	defer store.Unlock()

	err = json.NewDecoder(bufferedReader).Decode(store.contents)
	if err != nil {
		log.Print(__errRestoreJson, err.Error())
		return err
	}

    log.Print(__snapshotRestored)
    return nil
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
	_, existed := store.contents[content.ID]
	store.contents[content.ID] = content.Clone()

	if existed {
		store.bus.Publish("content:reindex", content)
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
