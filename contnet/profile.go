package contnet

import "sync"

// Interest is a floating point value that represents the level of user's interest.
type Interest float64

// TopicInterest is a structure that represents the level of user's interest for specified topic.
type TopicInterest struct {
	Topic    Topic
	Interest Interest
}
type TopicInterestFactory struct{}

func (factory *TopicInterestFactory) New(topic Topic, interest Interest) *TopicInterest {
	return &TopicInterest{
		Topic:    topic,
		Interest: interest,
	}
}

func (topicInterest *TopicInterest) Clone() *TopicInterest {
	return Object.TopicInterest.New(topicInterest.Topic, topicInterest.Interest)
}

// Alias for a slice of pointers to topic interests.
type TopicInterests []*TopicInterest

func (topicInterests TopicInterests) Clone() TopicInterests {
	out := TopicInterests{}

	for i := 0; i < len(topicInterests); i++ {
		out = append(out, topicInterests[i].Clone())
	}

	return out
}

// Profile contains information about topic interests for specified content consumer.
type Profile struct {
	sync.RWMutex
	ID             int64
	TopicInterests TopicInterests
}
type ProfileFactory struct{}

// Creates new Profile object
func (factory ProfileFactory) New(id int64, topicInterests TopicInterests) *Profile {
	return &Profile{
		ID:             id,
		TopicInterests: topicInterests.Clone(),
	}
}

// Thread-safe profile object clone (deep copy)
func (profile *Profile) Clone() *Profile {
	profile.RLock()
	defer profile.RUnlock()

	return Object.Profile.New(profile.ID, profile.TopicInterests)
}

// Thread-safe profile object update with action
func (profile *Profile) SaveAction(action *Action) {
	profile.Lock()
	defer profile.Unlock()

	// update with action
}
