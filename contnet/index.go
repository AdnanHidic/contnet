package contnet

import (
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

type Index struct {
	sync.RWMutex
	config   *NetConfig
	bus      *EventBus.EventBus
	contents *ContentStore
	index    map[Topic][]ID
}
type IndexFactory struct{}

type IndexNodeDescription struct {
	Keyword1 Keyword
	Keyword2 Keyword
	IDs      []ID
}

func (index *Index) GetForTopics(topics Topics) [][]ID {
	index.RLock()
	defer index.RUnlock()

	out := [][]ID{}

	for i := 0; i < len(topics); i++ {
		out = append(out, []ID{})

		if ids, exists := index.index[*topics[i]]; exists {
			out[i] = append(out[i], ids...)
		}
	}

	return out
}

func (index *Index) Describe() []*IndexNodeDescription {
	index.RLock()
	defer index.RUnlock()

	out := []*IndexNodeDescription{}

	for topic, mentions := range index.index {
		k1, k2 := topic.GetKeywords()
		desc := &IndexNodeDescription{
			Keyword1: k1,
			Keyword2: k2,
			IDs:      mentions,
		}
		out = append(out, desc)
	}

	return out
}

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

func (index *Index) __refresh() {
	for {
		index.Lock()
		log.Println("GLOBAL PAUSE: refreshing index ordering...")

		// for every entry
		for topic, mentions := range index.index {
			// get all contents
			contents := index.contents.GetMany(mentions)
			// sort all contents plus the new content & sort (best first)
			ContentBy(contentAgeCriteria).Sort(contents)
			// project slice to extract IDs
			index.index[topic] = __extractIDsFromContents(contents)
		}

		log.Println("GLOBAL UNPAUSE: refreshing index ordering finished.")
		index.Unlock()
		time.Sleep(index.config.CheckContentAgeInterval)
	}
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
	index.bus.Publish("topics:mentioned", topics, content.Popularity)
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
	index.bus.Publish("topics:unmentioned", topics, content.Popularity)
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
