package contnet

import "sync"

type Network struct {
	sync.RWMutex
	contentStore *ContentStore
	profileStorage *ProfileStore
}
type NetworkFactory struct{}

func (factory NetworkFactory) New() *Network {
	return &Network{
		contentStore: Object.ContentStore.New(),
		profileStorage: Object.ProfileStore.New(),
	}
}

func (net *Network) Save() error {
	return Errors.NotImplemented
}

func (net *Network) Load() error {
	return Errors.NotImplemented
}

// Attempts to update network object with content specified.
// If content did not exist, it is added to the network.
// If content did exist, it is updated.
func (net *Network) SaveContent(content *Content) error {
	contentInStorage := net.contentStore.Get(content.ID)
	// check if content object already exists in storage
	switch contentInStorage {
	case nil:
		// if content does not exist in storage
		net.contentStore.Create(content)
		return nil
	default:
		// if content does exist in storage
		net.contentStore.Update(contentInStorage, content)
		return nil
	}
}

// Attempts to update network object with action for profile specified.
func (net *Network) SaveAction(action *Action) error {
	// get related content's copy, if any
	relatedContent := net.contentStore.Get(action.ContentID)

	// if content does not exist
	if relatedContent == nil {
		return Errors.ContentNotFound
	}

	// injected related content to action
	action.Content = relatedContent

	// save action to profile
	net.profileStorage.SaveAction(action)

	// everything is okay
	return nil
}

func (net *Network) Select(profileID int64, page uint8) ([]*Content, error) {
	profile := net.profileStorage.Get(profileID)
    // if no profile found, return error
    if profile == nil {
        return nil, Errors.ProfileNotFound
    }
    // select content for profile specified
    content := net.contentStore.Select(profile, page)
    return content, nil
}

func (net *Network) Describe() *NetworkDescription {
	return &NetworkDescription{
		Contents: len(net.contentStore.contents),
		Profiles: len(net.profileStorage.profiles),
	}
}
