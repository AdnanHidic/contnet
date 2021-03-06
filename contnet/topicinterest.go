package contnet

import "sort"

// Interest is a floating point value that represents the level of user's interest.
type Interest float64

// TopicInterest is a structure that represents the level of user's interest for specified topic.
type TopicInterest struct {
	Topic              Topic
	Interest           Interest
	CumulativeInterest Interest
}
type TopicInterestFactory struct{}

type TopicInterestDescription struct {
	Keyword1           Keyword
	Keyword2           Keyword
	Interest           Interest
	CumulativeInterest Interest
}

var topicInterestCriteria = func(t1, t2 *TopicInterest) bool {
	return t1.Interest > t2.Interest
}

// function that defines ordering between trend objects
type TopicInterestBy func(t1, t2 *TopicInterest) bool

// method on the function type, sorts the argument slice according  to the function
func (topicInterestBy TopicInterestBy) Sort(topicInterests TopicInterests) {
	ts := &topicInterestSorter{
		topicInterests:  topicInterests,
		topicInterestBy: topicInterestBy,
	}
	sort.Sort(ts)
}

type topicInterestSorter struct {
	topicInterests  TopicInterests
	topicInterestBy func(t1, t2 *TopicInterest) bool
}

// Len is part of sort.Interface.
func (ts *topicInterestSorter) Len() int {
	return len(ts.topicInterests)
}

// Swap is part of sort.Interface.
func (ts *topicInterestSorter) Swap(i, j int) {
	ts.topicInterests[i], ts.topicInterests[j] = ts.topicInterests[j], ts.topicInterests[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ts *topicInterestSorter) Less(i, j int) bool {
	return ts.topicInterestBy(ts.topicInterests[i], ts.topicInterests[j])
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

func (factory *TopicInterestFactory) New(topic Topic, cumulativeInterest Interest) *TopicInterest {
	return &TopicInterest{
		Topic:              topic,
		CumulativeInterest: cumulativeInterest,
	}
}

func (topicInterest *TopicInterest) Clone() *TopicInterest {
	return &TopicInterest{
		Topic:              topicInterest.Topic,
		CumulativeInterest: topicInterest.CumulativeInterest,
		Interest:           topicInterest.Interest,
	}
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

const (
	__TRIM_INTEREST_LOWER_BOUND = 0.0001 // 0.001% interested
	__FAIRNESS_CONSTANT         = 0.0015
)

// value - between -1.0 and 1.0
func (topicInterests TopicInterests) Apply(topics Topics, value float64) TopicInterests {
	// foreach topic, add value to respective cumulative interest.
	// if no topic is registered, register it.
	// if cumulative interest falls below 0, remove interest
	interestSum := Interest(0.0)

	for i := 0; i < len(topics); i++ {
		found := false
		for j := 0; j < len(topicInterests); j++ {
			if topicInterests[j].Topic == *topics[i] {
				found = true
				topicInterests[j].CumulativeInterest += Interest(value)
				// if interest has become negative, remove interest (also if base interest is too low)
				if topicInterests[j].CumulativeInterest < 0 || topicInterests[j].Interest < __TRIM_INTEREST_LOWER_BOUND {
					topicInterests = append(topicInterests[:j], topicInterests[j+1:]...)
					break
				}

				interestSum += topicInterests[j].CumulativeInterest

				break
			}
		}
		// topic not found, create new interest if interest is positive, otherwise ignore it
		if !found && value > 0 {
			// apply fairness factor to give newer entries a better chance for survival
			fair := value + __FAIRNESS_CONSTANT
			interest := Object.TopicInterest.New(*topics[i], Interest(fair))
			topicInterests = append(topicInterests, interest)

			interestSum += Interest(fair)
		}
	}

	// now recalculate base interests (proportion)
	for i := 0; i < len(topicInterests); i++ {
		topicInterests[i].Interest = topicInterests[i].CumulativeInterest / interestSum
	}

	return topicInterests
}
