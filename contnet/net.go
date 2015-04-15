package contnet

import (
	"github.com/asaskevich/EventBus"
	"log"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type NetConfig struct {
	MaxContentAge           time.Duration
	CheckContentAgeInterval time.Duration
	ItemsPerPage            uint8
	NoveltyPct              float64
	SnapshotPath            string
	SnapshotInterval        time.Duration
}

type Net struct {
	sync.RWMutex
	config       *NetConfig
	bus          *EventBus.EventBus
	contentStore *ContentStore
	profileStore *ProfileStore
	trendStore   *TrendStore
	index        *Index
}
type NetFactory struct{}

func (factory NetFactory) New(config *NetConfig) *Net {
	bus := EventBus.New()
	contentStore := Object.ContentStore.New(config, bus)
	net := &Net{
		config:       config,
		bus:          bus,
		contentStore: contentStore,
		profileStore: Object.ProfileStore.New(),
		trendStore:   Object.TrendStore.New(bus),
		index:        Object.Index.New(config, bus, contentStore),
	}
	if err := net.Restore(); err != nil {
		log.Print("Failed to restore net object, proceeding empty.")
	}
	go net.__snapshot()
	go net.index.__refresh()

	return net
}

func (net *Net) __snapshot() {
	for {
		net.Snapshot()
		time.Sleep(net.config.SnapshotInterval)
	}
}

func (net *Net) Snapshot() error {
	if err := net.contentStore.Snapshot(net.config.SnapshotPath, "content"); err != nil {
		return err
	}

	if err := net.index.Snapshot(net.config.SnapshotPath, "index"); err != nil {
		return err
	}

	if err := net.profileStore.Snapshot(net.config.SnapshotPath, "profiles"); err != nil {
		return err
	}

	if err := net.trendStore.Snapshot(net.config.SnapshotPath, "trends"); err != nil {
		return err
	}

	return nil
}

func (net *Net) Restore() error {
	if err := net.contentStore.RestoreFromSnapshot(net.config.SnapshotPath, "content"); err != nil {
		return err
	}

	if err := net.index.RestoreFromSnapshot(net.config.SnapshotPath, "index"); err != nil {
		return err
	}

	if err := net.profileStore.RestoreFromSnapshot(net.config.SnapshotPath, "profiles"); err != nil {
		return err
	}

	if err := net.trendStore.RestoreFromSnapshot(net.config.SnapshotPath, "trends"); err != nil {
		return err
	}

	return nil
}

// Attempts to update network object with content specified.
// If content did not exist, it is added to the network.
// If content did exist, it is updated.
func (net *Net) SaveContent(content *Content) {
	net.contentStore.Upsert(content)
}

// Attempts to update network object with action for profile specified.
func (net *Net) SaveAction(action *Action) error {
	// get related content's copy, if any
	relatedContent := net.contentStore.Get(action.ContentID)

	// if content does not exist
	if relatedContent == nil {
		return Errors.ContentNotFound
	}

	// inject related content to action
	action.Content = relatedContent

	// save action to profile
	net.profileStore.Save(action)

	// everything is okay
	return nil
}

type SelectionCacheInfo struct {
	CurrentIndex int
	MaximumIndex int
}

func (net *Net) Select(profileID int64, page uint8) []ID {
	// get user's profile, if any
	profile := net.profileStore.Get(ID(profileID))

	// assign topic interests & sort
	var interests TopicInterests
	if profile != nil {
		TopicInterestBy(topicInterestCriteria).Sort(profile.TopicInterests)
		interests = profile.TopicInterests
	} else {
		interests = TopicInterests{}
	}

	totalContentsToRetrieve := int(page) * int(net.config.ItemsPerPage)
	totalContentsByInterest := int((1 - net.config.NoveltyPct) * float64(totalContentsToRetrieve))

	// select interesting contents
	interestingContents := net.__selectOfInterest(interests, totalContentsByInterest)
	//log.Println("Interesting: ", len(interestingContents))

	// select popular contents (how many remaining / 2)
	totalContentsByTrend := int((totalContentsToRetrieve - len(interestingContents)) / 2)
	trendingContents := net.__selectOfTrending(interestingContents, totalContentsByTrend)
	//log.Println("Trending: ", len(trendingContents))

	// fill the remainder randomly
	totalContentsRemaining := totalContentsToRetrieve - len(interestingContents) - len(trendingContents)
	out := append(interestingContents, trendingContents...)
	remainingContents := net.__selectOfRemaining(out, totalContentsRemaining)
	//log.Println("Remaining: ", len(remainingContents))

	// prep the final result
	out = append(out, remainingContents...)
	return out
}

func (net *Net) __selectOfInterest(interests TopicInterests, howMany int) []ID {
	// abort early if nothing to select..
	if len(interests) == 0 {
		return []ID{}
	}

	// we have sorted topic interests for this user.
	// now we extract topics from topic-interest objects for querying the index
	// and we use dem iterations to form cumulative probabilities slice to be used later for randomly drawing topics
	interestTopics := Topics{}
	cumulativeProbabilities := []float64{}
	for i := 0; i < len(interests); i++ {
		interestTopics = append(interestTopics, &interests[i].Topic)
		cumulativeProbabilities = append(cumulativeProbabilities, float64(interests[i].Interest))
		if i > 0 {
			cumulativeProbabilities[i] += cumulativeProbabilities[i-1]
		}
	}

	// use the previously created topics slice for querying the index
	topicContents := net.index.GetForTopics(interestTopics)

	// now we have all content IDs for all topics interesting to this user
	// to control where we're at, we will create temporary cache that will mark extraction location for each topic
	cache := map[Topic]*SelectionCacheInfo{}
	// filling in the cache
	for i := 0; i < len(interests); i++ {
		interest := interests[i]
		cache[interest.Topic] = &SelectionCacheInfo{
			CurrentIndex: 0,
			MaximumIndex: len(topicContents[i]) - 1,
		}
	}

	// now we have the following data available:
	// 1. topics of interest for user sorted in descending order by how interesting they are
	// 2. for each topic of interest there is an array of ID's available for that topic
	// 3. map showing what index to choose for each topic
	// 4. cumulative probabilities

	// while there is anything interesting for user available & we haven't selected as many items as needed
	out := []ID{}
	for len(cache) > 0 && len(out) < howMany {
		// select a random topic of interest
		topic, i := __drawRandomTopicInterest(interests, cumulativeProbabilities)

		// select its best content (first index available)
		indexToSelect := cache[topic].CurrentIndex

		// current index can be ahead of the max index in some special cases
		if indexToSelect <= cache[topic].MaximumIndex {
			selectedID := topicContents[i][indexToSelect]
			if !__idsContainID(out, selectedID) {
				out = append(out, selectedID)
			}
		}

		// advance current index in cache
		cache[topic].CurrentIndex++

		// if we overshot the maximum index for this topic, remove it
		if cache[topic].CurrentIndex > cache[topic].MaximumIndex {
			// remove from cache
			delete(cache, topic)
			// remove from interest topics
			interestTopics = append(interestTopics[:i], interestTopics[i+1:]...)
			// remove from topic contents
			topicContents = append(topicContents[:i], topicContents[i+1:]...)
			// remove from interests
			interests = append(interests[:i], interests[i+1:]...)
			// recalculate interests
			cumulativeProbabilities = __recalculateCumulativeProbabilities(interests)
		}
	}

	return out
}

func __idsContainID(ids []ID, id ID) bool {
	for i := 0; i < len(ids); i++ {
		if ids[i] == id {
			return true
		}
	}
	return false
}

func __drawRandomTopicInterest(topicInterests TopicInterests, cumulativeProbabilities []float64) (Topic, int) {
	randomFloat := rand.Float64()
	for i := 0; i < len(cumulativeProbabilities); i++ {
		if cumulativeProbabilities[i] >= randomFloat {
			return topicInterests[i].Topic, i
		}
	}
	// this is guaranteed not to happen
	return topicInterests[len(topicInterests)-1].Topic, len(topicInterests) - 1
}

func __recalculateCumulativeProbabilities(topicInterests TopicInterests) []float64 {
	// calculate sum of cumulative interests
	sum := float64(0)
	for i := 0; i < len(topicInterests); i++ {
		sum += float64(topicInterests[i].CumulativeInterest)
	}

	output := []float64{}

	for i := 0; i < len(topicInterests); i++ {
		output = append(output, float64(topicInterests[i].CumulativeInterest)/sum)
		if i > 0 {
			output[i] += output[i-1]
		}
	}

	return output
}

func (net *Net) __selectOfTrending(ignore []ID, howMany int) []ID {
	// select top 10 trending topics
	trendingTopics := net.trendStore.GetTopN(10)

	// nothing trending, bail early
	if len(trendingTopics) == 0 {
		return []ID{}
	}

	// get content IDs for trending topics
	contentIDs := net.index.GetForTopics(trendingTopics)

	// cache indices foreach contentID
	cache := map[Topic]*SelectionCacheInfo{}
	// filling in the cache
	for i := 0; i < len(trendingTopics); i++ {
		cache[*trendingTopics[i]] = &SelectionCacheInfo{
			CurrentIndex: 0,
			MaximumIndex: len(contentIDs[i]) - 1,
		}
	}

	nextTrendInd := 0
	out := []ID{}
	for len(cache) > 0 && len(out) < howMany {
		trendingTopic := *trendingTopics[nextTrendInd]

		// select its best content (first index available)
		indexToSelect := cache[trendingTopic].CurrentIndex

		// current index can be ahead of the max index in some special cases
		if indexToSelect <= cache[trendingTopic].MaximumIndex {
			selectedID := contentIDs[nextTrendInd][indexToSelect]

			if !__idsContainID(out, selectedID) && !__idsContainID(ignore, selectedID) {
				out = append(out, selectedID)
			}
		}

		// advance current index in cache
		cache[trendingTopic].CurrentIndex++

		// if we overshot the maximum index for this topic, remove it
		if cache[trendingTopic].CurrentIndex > cache[trendingTopic].MaximumIndex {
			// remove from cache
			delete(cache, trendingTopic)
			// remove from trending topics
			trendingTopics = append(trendingTopics[:nextTrendInd], trendingTopics[nextTrendInd+1:]...)
			// go back one position
			nextTrendInd--
		}

		// reset the counter
		nextTrendInd++
		if nextTrendInd == len(trendingTopics) {
			nextTrendInd = 0
		}
	}

	return out
}

func (net *Net) __selectOfRemaining(ignore []ID, howMany int) []ID {
	// if we are asking for more content items then there are available in content store
	if howMany >= (net.contentStore.Count() - len(ignore)) {
		ids := net.contentStore.GetAllContentIDs()
		// filter out ignored ones
		for i := 0; i < len(ids); i++ {
			if __idsContainID(ignore, ids[i]) {
				ids = append(ids[:i], ids[i+1:]...)
				i--
			}
		}
		return ids

	}

	// otherwise, draw random elements until howMany is satisifed. ignore duplicates of course.
	out := []ID{}

	for len(out) < howMany {
		randomID := net.contentStore.GetAnyContentID()
		if !__idsContainID(out, randomID) && !__idsContainID(ignore, randomID) {
			out = append(out, randomID)
		}
	}

	return out

}

func (net *Net) Describe() *NetDescription {
	return &NetDescription{
		Contents: net.contentStore.Describe(),
		Index:    net.index.Describe(),
		Profiles: net.profileStore.Describe(),
		Trends:   net.trendStore.Describe(),
	}
}
