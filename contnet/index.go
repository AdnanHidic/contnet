package contnet

import (
	"github.com/asaskevich/EventBus"
	"sync"
)

// TODO: index should periodically re-sort content IDs for topic
type Index struct {
	sync.RWMutex
	config   *NetConfig
	bus      *EventBus.EventBus
	contents *ContentStore
	index    map[Topic][]ID
}
type IndexFactory struct{}

func (factory IndexFactory) New(config *NetConfig, bus *EventBus.EventBus, contentStore *ContentStore) *Index {
	index := &Index{
		config:   config,
		bus:      bus,
		index:    map[Topic][]ID{},
		contents: contentStore,
	}

	bus.SubscribeAsync("content:index", index.Index, false)
	bus.SubscribeAsync("content:reindex", index.Reindex, false)
	bus.SubscribeAsync("content:removed", index.Remove, false)

	return index
}

func (index *Index) Snapshot(path, filename string) error {
	index.RLock()
	defer index.RUnlock()

	return __snapshot(path, filename, &index.index)
}

func (index *Index) RestoreFromSnapshot(path, filename string) error {
	index.Lock()
	defer index.Unlock()

	_, err := __restoreFromSnapshot(path, filename, &index.index)

	return err
}

// Index previously unindexed content.
func (index *Index) Index(content *Content) {
	// get topics for this content
	topics := content.Keywords.GetTopics()
	// foreach topic
	for i := 0; i < len(topics); i++ {
		// find topic in index, if any
		index.RLock()
		mentions, exists := index.index[*topics[i]]
		index.RUnlock()

		index.Lock()
		if !exists {
			// if topic not indexed yet, index it and add this content as its first mention
			mentions = []ID{content.ID}
			index.index[*topics[i]] = mentions
		} else {
			// if topic is indexed, add this content to topic mentions
			index.index[*topics[i]] = index.addMention(mentions, content)
		}
		index.Unlock()
	}
	// notify any listener that these topics have been mentioned.
	index.bus.Publish("topics:mentioned", topics)
}

func (index *Index) Reindex(old, new *Content) {
	// remove old
	index.Remove(old)
	// add new
	index.Index(new)
}

func (index *Index) Remove(content *Content) {
	// get topics for removed content
	topics := content.Keywords.GetTopics()
	// foreach topic
	for i := 0; i < len(topics); i++ {
		// find topic in index, if any
		index.RLock()
		mentions, exists := index.index[*topics[i]]
		index.RUnlock()
		if !exists {
			return
		}

		index.Lock()
		index.index[*topics[i]] = index.removeMention(mentions, content.ID)
		index.Unlock()
	}
	// notify anly listener that these topics have been unmentioned
	index.bus.Publish("topics:unmentioned", topics)
}

func (index *Index) addMention(mentions []ID, content *Content) []ID {
	// get contents from store by ids provided (in order as specified) and extend with the new one
	contents := index.contents.GetMany(mentions)
	contents = append(contents, content)

	// sort all contents plus the new content & sort (best first)
	ContentBy(contentAgeCriteria).Sort(contents)

	// project slice to extract IDs
	return __extractIDsFromContents(contents)
}

func (index *Index) removeMention(mentions []ID, mention ID) []ID {
	// find & remove id from mentions
	for i := 0; i < len(mentions); i++ {
		// if mention was found
		if mentions[i] == mention {
			// remove it
			mentions = append(mentions[:i], mentions[i+1:]...)
			break
		}
	}

	// if 0 or 1 mentions remain, there's nothing else to do
	if len(mentions) < 2 {
		return mentions
	}

	// otherwise, we have to recalculate everything!
	// get contents from store by ids provided and in-order specified
	contents := index.contents.GetMany(mentions)

	// calculate age for all contents  & sort (best first)
	ContentBy(contentAgeCriteria).Sort(contents)

	// project slice to extract IDs
	return __extractIDsFromContents(contents)

}

func __extractIDsFromContents(contents []*Content) []ID {
	ids := []ID{}
	for i := 0; i < len(contents); i++ {
		ids = append(ids, contents[i].ID)
	}
	return ids
}
