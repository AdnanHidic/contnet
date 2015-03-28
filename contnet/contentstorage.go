package contnet

import "sync"

type ContentStorage struct {
	sync.RWMutex
	contents map[int64]*Content
}
type ContentStorageFactory struct{}

func (factory ContentStorageFactory) New() *ContentStorage {
	return &ContentStorage{
		contents: map[int64]*Content{},
	}
}

func (storage *ContentStorage) Get(id int64) *Content {
	storage.RLock()
	defer storage.RUnlock()

	if content, ok := storage.contents[id]; !ok {
		return nil
	} else {
		return content.Clone()
	}
}

func (storage *ContentStorage) Save(content *Content) {
	storage.Lock()
	defer storage.Unlock()

	storage.contents[content.ID] = content.Clone()
}

func (storage *ContentStorage) Delete(id int64) {
	storage.Lock()
	defer storage.Unlock()

	delete(storage.contents, id)
}
