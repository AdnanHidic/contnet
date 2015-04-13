package contnet

import "sync"

type ContentStore struct {
	sync.RWMutex
	contents map[int64]*Content
	index    map[Topic]Contents
}
type ContentStoreFactory struct{}

func (factory ContentStoreFactory) New() *ContentStore {
	return &ContentStore{
		contents: map[int64]*Content{},
		index:    map[Topic]Contents{},
	}
}

func (storage *ContentStore) Get(id int64) *Content {
	storage.RLock()
	defer storage.RUnlock()

	if content, ok := storage.contents[id]; !ok {
		return nil
	} else {
		return content.Clone()
	}
}

func (storage *ContentStore) Create(content *Content) {
	storage.Lock()
	defer storage.Unlock()

	// save it to map of contents
	storage.contents[content.ID] = content.Clone()

	// get all topics
	topics := content.Keywords.GetTopics()
	// foreach topic, add this content

	storage.addContentToIndex(topics, content)
}

func (storage *ContentStore) Update(old, new *Content) {
	storage.Lock()
	defer storage.Unlock()

	// update contents object
	storage.contents[old.ID] = new

	// now update index
	// get all topics for old and new content
	oldTopics := old.Keywords.GetTopics()
	newTopics := new.Keywords.GetTopics()

	// remove content from all old topics
	storage.removeContentFromIndex(oldTopics, old)

	// add content to all new topics
	storage.addContentToIndex(newTopics, new)
}

func (storage *ContentStore) Delete(id int64) {
	storage.Lock()
	defer storage.Unlock()

	// remove content from contents
	content, exists := storage.contents[id]

	if !exists {
		return
	}

	// remove from index
	topics := content.Keywords.GetTopics()

	// remove content from topics
	storage.removeContentFromIndex(topics, content)

	// finally, remove from contents
	delete(storage.contents, content.ID)

}

func (storage *ContentStore) Select(profile *Profile, page uint8) []*Content {
    storage.RLock()
    defer storage.RUnlock()

    return nil
}

func (storage *ContentStore) addContentToIndex(topics Topics, content *Content) {
	// foreach topic
	for i := 0; i < len(topics); i++ {
		// try to get indexed topic contents
		contents, isTopicRegistered := storage.index[*topics[i]]
		// if no topic contents indexed
		if !isTopicRegistered {
			// create topic contents
			contents = Contents{}
			// save
			storage.index[*topics[i]] = contents
		}
		// finally, add content to topic contents
		contents.Add(content)
	}
}

func (storage *ContentStore) removeContentFromIndex(topics Topics, content *Content) {
	// foreach topic
	for i := 0; i < len(topics); i++ {
		// try to get indexed topic contents
		contents, isTopicRegistered := storage.index[*topics[i]]
		// if no topic contents indexed
		if isTopicRegistered {
			// finally, remove content from topic contents
			contents.Remove(content)
		}
	}
}
