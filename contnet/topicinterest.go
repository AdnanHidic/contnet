package contnet

// Interest is a floating point value that represents the level of user's interest.
type Interest float64

// TopicInterest is a structure that represents the level of user's interest for specified topic.
type TopicInterest struct {
	Topic    Topic
	Interest Interest
}
type TopicInterestFactory struct{}

func (factory *TopicInterestFactory) New(topic Topic, interest Interest) *TopicInterest {
	return &TopicInterest{
		Topic:    topic,
		Interest: interest,
	}
}

func (topicInterest *TopicInterest) Clone() *TopicInterest {
	return Object.TopicInterest.New(topicInterest.Topic, topicInterest.Interest)
}

// Alias for a slice of pointers to topic interests.
type TopicInterests []*TopicInterest

func (topicInterests TopicInterests) Clone() TopicInterests {
	out := TopicInterests{}

	for i := 0; i < len(topicInterests); i++ {
		out = append(out, topicInterests[i].Clone())
	}

	return out
}
