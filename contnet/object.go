package contnet

// registry of all object factories
var Object = struct {
	Topic          TopicFactory
	Content        ContentFactory
	ContentStore ContentStoreFactory
	TopicInterest  TopicInterestFactory
	Profile        ProfileFactory
	ProfileStore ProfileStoreFactory
	Network        NetworkFactory
	Action         ActionFactory
	ActionArgument ActionArgumentFactory
}{
	Topic:          TopicFactory{},
	Content:        ContentFactory{},
	ContentStore: ContentStoreFactory{},
	TopicInterest:  TopicInterestFactory{},
	Profile:        ProfileFactory{},
	ProfileStore: ProfileStoreFactory{},
	Network:        NetworkFactory{},
	Action:         ActionFactory{},
	ActionArgument: ActionArgumentFactory{},
}
