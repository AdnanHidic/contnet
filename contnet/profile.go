package contnet

import (
	"math"
)

// Profile contains information about topic interests for specified content consumer.
type Profile struct {
	ID             ID
	TopicInterests TopicInterests
}
type ProfileFactory struct{}

// Creates new Profile object
func (factory ProfileFactory) New(id ID, topicInterests TopicInterests) *Profile {
	return &Profile{
		ID:             id,
		TopicInterests: topicInterests.Clone(),
	}
}

func (profile *Profile) Clone() *Profile {
	return Object.Profile.New(profile.ID, profile.TopicInterests)
}

type ProfileDescription struct {
	ID             ID                          `json:"id"`
	TopicInterests []*TopicInterestDescription `json:"topicInterests"`
}

func (profile *Profile) Describe() *ProfileDescription {
	return &ProfileDescription{
		ID:             profile.ID,
		TopicInterests: profile.TopicInterests.Describe(),
	}
}

// Thread-safe profile object update with action
func (profile *Profile) SaveAction(action *Action) {
	switch action.Type {
	case ActionTypes.Read:
		profile.__saveReadAction(action)
	case ActionTypes.Upvote:
		profile.__saveUpvoteAction(action)
	case ActionTypes.Downvote:
		profile.__saveDownvoteAction(action)
	default:
		return
	}
}

const (
	__READ_SECOND_VALUE = 0.15
	__UPVOTE_VALUE      = 1.0
	__DOWNVOTE_VALUE    = -1.0
)

func (profile *Profile) __saveReadAction(action *Action) {
	// get topics for read content
	topics := action.Content.Keywords.GetTopics()
	// now we calculate the "value" of this action, e.g. based on the length of the reading session
	value := __READ_SECOND_VALUE * math.Log10(float64(action.Arguments.GetArgumentByName("duration-seconds").GetValueInteger()))
	if value > 1.0 {
		value = 1.0
	}

	profile.TopicInterests = profile.TopicInterests.Apply(topics, value)
}

func (profile *Profile) __saveUpvoteAction(action *Action) {
	// get topics for upvoted content
	topics := action.Content.Keywords.GetTopics()
	value := __UPVOTE_VALUE

	profile.TopicInterests.Apply(topics, value)
}

func (profile *Profile) __saveDownvoteAction(action *Action) {
	// get topics for downvoted content
	topics := action.Content.Keywords.GetTopics()
	value := __DOWNVOTE_VALUE

	profile.TopicInterests.Apply(topics, value)
}
