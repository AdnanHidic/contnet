package contnet

import "sync"

type ProfileStore struct {
	sync.RWMutex
	profiles map[ID]*Profile
}
type ProfileStoreFactory struct{}

func (factory ProfileStoreFactory) New() *ProfileStore {
	return &ProfileStore{
		profiles: map[ID]*Profile{},
	}
}

func (store *ProfileStore) Get(id ID) *Profile {
	store.RLock()
	defer store.RUnlock()

	if profile, ok := store.profiles[id]; !ok {
		return nil
	} else {
		return profile.Clone()
	}
}

func (store *ProfileStore) Save(profile *Profile) {
	store.Lock()
	defer store.Unlock()

	store.profiles[profile.ID] = profile.Clone()
}

func (store *ProfileStore) Delete(id ID) {
	store.Lock()
	defer store.Unlock()

	delete(store.profiles, id)
}

func (store *ProfileStore) SaveAction(action *Action) {
	store.RLock()
	defer store.RUnlock()

	var (
		profile       *Profile
		profileExists bool
	)
	// get profile, if any
	profile, profileExists = store.profiles[action.ProfileID]

	// if profile does not exist, create a new one and save it
	if !profileExists {
		profile = Object.Profile.New(action.ProfileID, TopicInterests{})
		store.Save(profile)
	}

	// update profile with this action
	profile.SaveAction(action)
}
