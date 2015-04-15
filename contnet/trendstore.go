package contnet

import (
	"github.com/asaskevich/EventBus"
	"sort"
	"sync"
)

type Trend struct {
	Topic      Topic
	Popularity float64
}

var trendPopularityCriteria = func(t1, t2 *Trend) bool {
	return t1.Popularity > t2.Popularity
}

// function that defines ordering between trend objects
type TrendBy func(t1, t2 *Trend) bool

// method on the function type, sorts the argument slice according  to the function
func (trendBy TrendBy) Sort(trends []*Trend) {
	ts := &trendSorter{
		trends:  trends,
		trendBy: trendBy,
	}
	sort.Sort(ts)
}

type trendSorter struct {
	trends  []*Trend
	trendBy func(t1, t2 *Trend) bool
}

// Len is part of sort.Interface.
func (ts *trendSorter) Len() int {
	return len(ts.trends)
}

// Swap is part of sort.Interface.
func (ts *trendSorter) Swap(i, j int) {
	ts.trends[i], ts.trends[j] = ts.trends[j], ts.trends[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ts *trendSorter) Less(i, j int) bool {
	return ts.trendBy(ts.trends[i], ts.trends[j])
}

type TrendDescription struct {
	Keyword1   Keyword
	Keyword2   Keyword
	Popularity float64
}

func (trend *Trend) Describe() *TrendDescription {
	k1, k2 := trend.Topic.GetKeywords()
	return &TrendDescription{
		Keyword1:   k1,
		Keyword2:   k2,
		Popularity: trend.Popularity,
	}
}

type TrendStore struct {
	sync.RWMutex
	bus    *EventBus.EventBus
	cache  map[Topic]*Trend
	trends []*Trend
}
type TrendStoreFactory struct{}

func (factory TrendStoreFactory) New(bus *EventBus.EventBus) *TrendStore {
	store := &TrendStore{
		bus:    bus,
		cache:  map[Topic]*Trend{},
		trends: []*Trend{},
	}

	bus.SubscribeAsync("topics:mentioned", store.Register, false)
	bus.SubscribeAsync("topics:unmentioned", store.Unregister, false)

	return store
}

func (store *TrendStore) Describe() []*TrendDescription {
	store.RLock()
	defer store.RUnlock()

	out := []*TrendDescription{}
	for i := 0; i < len(store.trends); i++ {
		out = append(out, store.trends[i].Describe())
	}

	return out
}

func (store *TrendStore) Snapshot(path, filename string) error {
	store.RLock()
	defer store.RUnlock()

	return __snapshot(path, filename, &store.trends)
}

func (store *TrendStore) RestoreFromSnapshot(path, filename string) error {
	store.Lock()
	defer store.Unlock()

	_, err := __restoreFromSnapshot(path, filename, &store.trends)
	// fill the cache
	for i := 0; i < len(store.trends); i++ {
		store.cache[store.trends[i].Topic] = store.trends[i]
	}

	return err
}

func (store *TrendStore) GetTopN(n int) Topics {
	store.RLock()
	defer store.RUnlock()

	topNTrends := store.trends[0:n]

	out := Topics{}

	for i := 0; i < len(topNTrends); i++ {
		out = append(out, &topNTrends[i].Topic)
	}

	return out

}

func (store *TrendStore) Register(topics Topics, popularity float64) {
	store.Lock()
	defer store.Unlock()

	// foreach topic
	for i := 0; i < len(topics); i++ {
		// check if topic is cached. If not cached, this is this topic's first mention..
		if trend, exists := store.cache[*topics[i]]; !exists {
			trend = &Trend{
				Topic:      *topics[i],
				Popularity: popularity,
			}

			store.cache[*topics[i]] = trend
			store.trends = append(store.trends, trend)
		} else {
			trend.Popularity += popularity
		}
	}

	// sort trends
	TrendBy(trendPopularityCriteria).Sort(store.trends)
}

func (store *TrendStore) Unregister(topics Topics, popularity float64) {
	store.Lock()
	defer store.Unlock()

	// foreach topic
	for i := 0; i < len(topics); i++ {
		// check if topic is cached. If cached, decrease popularity
		if trend, exists := store.cache[*topics[i]]; exists {
			trend.Popularity -= popularity
			// remove if popularity becomes too low
			if trend.Popularity <= 0 {
				// delete from cache
				delete(store.cache, *topics[i])
				// delete from trends
				for j := 0; j < len(store.trends); j++ {
					if store.trends[j].Topic == *topics[i] {
						store.trends = append(store.trends[:j], store.trends[j+1:]...)
						break
					}
				}
			}
		}
	}

	// sort trends
	TrendBy(trendPopularityCriteria).Sort(store.trends)
}
