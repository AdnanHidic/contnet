package contnet

// Interest is a floating point value that represents the level of user's interest.
type Interest float64

// TopicInterest is a structure that represents the level of user's interest for specified topic.
type TopicInterest struct {
	Topic              Topic
	Interest           Interest
	CumulativeInterest Interest ``
}
type TopicInterestFactory struct{}

type TopicInterestDescription struct {
	Keyword1           Keyword
	Keyword2           Keyword
	Interest           Interest
	CumulativeInterest Interest
}

func (topicInterest *TopicInterest) Describe() *TopicInterestDescription {
	k1, k2 := topicInterest.Topic.GetKeywords()
	return &TopicInterestDescription{
		Keyword1:           k1,
		Keyword2:           k2,
		Interest:           topicInterest.Interest,
		CumulativeInterest: topicInterest.CumulativeInterest,
	}
}

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

func (topicInterests TopicInterests) Describe() []*TopicInterestDescription {
	out := []*TopicInterestDescription{}

	for i := 0; i < len(topicInterests); i++ {
		out = append(out, topicInterests[i].Describe())
	}

	return out
}

// value - between -1.0 and 1.0
func (topicInterests TopicInterests) Apply(topics Topics, value float64) {

}
