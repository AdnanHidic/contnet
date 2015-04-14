package contnet

import "sync"

// Profile contains information about topic interests for specified content consumer.
type Profile struct {
	sync.RWMutex
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

	switch action.Type {
	case ActionTypes.Read:
        profile.__saveReadAction(action)
    default:
        return
	}
}


func (profile *Profile) __saveReadAction(action *Action) {
    // TODO: algorithm for maintaining user profiles
}