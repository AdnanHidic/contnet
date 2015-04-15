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

func (store *ProfileStore) Describe() []*ProfileDescription {
	store.RLock()
	defer store.RUnlock()

	out := []*ProfileDescription{}

	for _, profile := range store.profiles {
		out = append(out, profile.Describe())
	}

	return out
}

func (store *ProfileStore) Snapshot(path, filename string) error {
	store.RLock()
	defer store.RUnlock()

	return __snapshot(path, filename, &store.profiles)
}

func (store *ProfileStore) RestoreFromSnapshot(path, filename string) error {
	store.Lock()
	defer store.Unlock()

	_, err := __restoreFromSnapshot(path, filename, &store.profiles)
	return err
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

func (store *ProfileStore) Delete(id ID) {
	store.Lock()
	defer store.Unlock()

	delete(store.profiles, id)
}

func (store *ProfileStore) Save(action *Action) {
	store.Lock()
	defer store.Unlock()

	var (
		profile       *Profile
		profileExists bool
	)
	// get profile, if any
	profile, profileExists = store.profiles[action.ProfileID]
	// if profile does not exist, create a new one and save it
	if !profileExists {
		profile = Object.Profile.New(action.ProfileID, TopicInterests{})
		store.profiles[profile.ID] = profile
	}

	// update profile with this action
	profile.SaveAction(action)
}
