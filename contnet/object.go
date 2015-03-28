package contnet

// registry of all object factories
var Object = struct {
	Topic          TopicFactory
	Content        ContentFactory
	ContentStorage ContentStorageFactory
	TopicInterest  TopicInterestFactory
	Profile        ProfileFactory
	ProfileStorage ProfileStorageFactory
	NetworkConfig  NetworkConfigFactory
	Network        NetworkFactory
}{
	Topic:          TopicFactory{},
	Content:        ContentFactory{},
	ContentStorage: ContentStorageFactory{},
	TopicInterest:  TopicInterestFactory{},
	Profile:        ProfileFactory{},
	ProfileStorage: ProfileStorageFactory{},
	NetworkConfig:  NetworkConfigFactory{},
	Network:        NetworkFactory{},
}
