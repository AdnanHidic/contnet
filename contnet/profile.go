package contnet

type Interest float64

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

type TopicInterests []*TopicInterest

func (topicInterests TopicInterests) Clone() TopicInterests {
	out := TopicInterests{}

	for i := 0; i < len(topicInterests); i++ {
		out = append(out, topicInterests[i].Clone())
	}

	return out
}

type Profile struct {
	ID             int64
	TopicInterests TopicInterests
}
type ProfileFactory struct{}

func (factory ProfileFactory) New(id int64, topicInterests TopicInterests) *Profile {
	return &Profile{
		ID:             id,
		TopicInterests: topicInterests.Clone(),
	}
}

func (profile *Profile) Clone() *Profile {
	return Object.Profile.New(profile.ID, profile.TopicInterests)
}
