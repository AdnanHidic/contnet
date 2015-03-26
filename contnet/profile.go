package contnet

import "sync"

type Interest struct {
	Topic Topic
}

type Profile struct {
	ID int64
}

func (profile *Profile) Clone() *Profile {
	return nil
}

type ProfileStorage struct {
	sync.RWMutex
	profiles map[int64]*Profile
}

func NewProfileStorage() *ProfileStorage {
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
