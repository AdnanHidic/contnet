package contnet

import "sync"

type ContentStorage struct {
	sync.RWMutex
	content map[int64]*Content
}
type ContentStorageFactory struct{}

func (factory ContentStorageFactory) New() *ContentStorage {
	return &ContentStorage{
		content: map[int64]*Content{},
	}
}

func (storage *ContentStorage) Get(id int64) *Content {
	storage.RLock()
	defer storage.RUnlock()

	if content, ok := storage.content[id]; !ok {
		return nil
	} else {
		return content.Clone()
	}
}

func (storage *ContentStorage) Save(content *Content) {
	storage.Lock()
	defer storage.Unlock()

	storage.content[content.ID] = content.Clone()
}

func (storage *ContentStorage) Delete(id int64) {
	storage.Lock()
	defer storage.Unlock()

	delete(storage.content, id)
}
