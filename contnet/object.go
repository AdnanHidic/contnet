package contnet

// registry of all object factories
var Object = struct {
	Topic          TopicFactory
	Content        ContentFactory
	ContentStorage ContentStorageFactory
	TopicInterest  TopicInterestFactory
	Profile        ProfileFactory
	ProfileStorage ProfileStorageFactory
}{
	Topic:          TopicFactory{},
	Content:        ContentFactory{},
	ContentStorage: ContentStorageFactory{},
	TopicInterest:  TopicInterestFactory{},
	Profile:        ProfileFactory{},
	ProfileStorage: ProfileStorageFactory{},
}
