package contnet

// registry of all object factories
var Object = struct {
	Topic          TopicFactory
	Content        ContentFactory
	ContentStore   ContentStoreFactory
	TopicInterest  TopicInterestFactory
	Profile        ProfileFactory
	ProfileStore   ProfileStoreFactory
	Index          IndexFactory
	Net            NetFactory
	Action         ActionFactory
	ActionArgument ActionArgumentFactory
	TrendStore     TrendStoreFactory
}{
	Topic:          TopicFactory{},
	Content:        ContentFactory{},
	ContentStore:   ContentStoreFactory{},
	TopicInterest:  TopicInterestFactory{},
	Profile:        ProfileFactory{},
	ProfileStore:   ProfileStoreFactory{},
	Index:          IndexFactory{},
	Net:            NetFactory{},
	Action:         ActionFactory{},
	ActionArgument: ActionArgumentFactory{},
	TrendStore:     TrendStoreFactory{},
}
