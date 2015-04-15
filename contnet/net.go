package contnet

import (
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

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

func (net *Net) Select(profileID int64, page uint8) ([]ID) {
	// get user's profile, if any
	profile := net.profileStore.Get(ID(profileID))

	// assign topic interests
	var interests TopicInterests
	if profile != nil {
		TopicInterestBy(topicInterestCriteria).Sort(profile.TopicInterests)
	} else {
		interests = TopicInterests{}
	}

    totalContentsToRetrieve := int(page) * int(net.config.ItemsPerPage)
    totalContentsByInterest := int((1-net.config.NoveltyPct)*float64(totalContentsToRetrieve))

    // select interesting contents
    interestingContents := net.__selectOfInterest(interests, totalContentsByInterest)

    // select popular contents (how many remaining / 2)
    totalContentsByTrend := int((totalContentsToRetrieve - len(interestingContents)) / 2)
    trendingContents := net.__selectOfTrending(interestingContents, totalContentsByTrend)

    // fill the remainder randomly
    totalContentsRemaining := totalContentsToRetrieve - len(interestingContents) - len(trendingContents)
    out := append(interestingContents, trendingContents...)
    remainingContents := net.__selectOfRemaining(out, totalContentsRemaining)

    // prep the final result
    out = append(out, remainingContents...)
    return out
}

func (net *Net) __selectOfInterest(interests TopicInterests, howMany int) []ID {
    // now we have sorted topic interests for this user.
    // next step is to create a map that will store information about how many contents have been used from a topic
    // eager loading for now
    interestTopics := Topics{}
    for i:=0;i<len(interests);i++ {
        interestTopics = append(interestTopics, &interests[i].Topic)
    }

    topicContents := net.index.GetForTopics(interestTopics)
    cache := map[Topic]*SelectionCacheInfo{}

    // fill it
    for i:=0;i<len(interests);i++ {
        interest := interests[i]
        cache[interest.Topic] = &SelectionCacheInfo{
            CurrentIndex: 0,
            MaximumIndex: len(topicContents[i]),
        }
    }

    // now we have the following data available:
    // 1. topics of interest for user sorted in descending order by how interesting they are
    // 2. for each topic of interest there is an array of ID's available for that topic
    // 3. map showing what index to choose for each topic

    // go go go
    // while there is anything interesting for user available & we haven't selected as many items as needed
    // for len(cache)>0 &&

    return nil
}

func (net *Net) __selectOfTrending(ignore []ID, howMany int) []ID {
    return nil
}

func (net *Net) __selectOfRemaining(ignore []ID, howMany int) []ID {
    return nil
}

func (net *Net) Describe() *NetDescription {
	return &NetDescription{
		Contents: net.contentStore.Describe(),
		Index:    net.index.Describe(),
		Profiles: net.profileStore.Describe(),
		Trends:   net.trendStore.Describe(),
	}
}
