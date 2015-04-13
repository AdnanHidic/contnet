package contnet

import (
	"github.com/asaskevich/EventBus"
	"sync"
)

type Index struct {
	sync.RWMutex
	bus      EventBus.EventBus
	contents *ContentStore
	index    map[Topic][]ID
}
type IndexFactory struct{}

func (factory IndexFactory) New(bus *EventBus.EventBus, contentStore *ContentStore) *Index {
	return &Index{
		bus:      bus,
		index:    map[Topic][]ID{},
		contents: contentStore,
	}
}

func (index *Index) Index(content *Content) {

}

func (index *Index) Reindex(content *Content) {

}

func (index *Index) Select(profileID ID, page uint8) []ID {
	return nil
}
