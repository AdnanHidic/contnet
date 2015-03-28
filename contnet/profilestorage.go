package contnet

import "sync"

type ProfileStorage struct {
	sync.RWMutex
	profiles map[int64]*Profile
}
type ProfileStorageFactory struct{}

func (factory ProfileStorageFactory) New() *ProfileStorage {
	return &ProfileStorage{
		profiles: map[int64]*Profile{},
	}
}

func (storage *ProfileStorage) Get(id int64) *Profile {
	storage.RLock()
	defer storage.RUnlock()

	if profile, ok := storage.profiles[id]; !ok {
		return nil
	} else {
		return profile.Clone()
	}
}

func (storage *ProfileStorage) Save(profile *Profile) {
	storage.Lock()
	defer storage.Unlock()

	storage.profiles[profile.ID] = profile.Clone()
}

func (storage *ProfileStorage) Delete(id int64) {
	storage.Lock()
	defer storage.Unlock()

	delete(storage.profiles, id)
}

func (storage *ProfileStorage) SaveAction(action Action) {
	storage.RLock()
	defer storage.RUnlock()

	var (
		profile       *Profile
		profileExists bool
	)
	// get profile, if any
	profile, profileExists = storage.profiles[action.ProfileID]

	// if profile does not exist, create a new one and save it
	if !profileExists {
		profile = Object.Profile.New(action.ProfileID, TopicInterests{})
		storage.Save(profile)
	}

	// update profile with this action
	profile.SaveAction(action)
}
