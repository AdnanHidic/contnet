package contnet

import "sync"

// Keyword is represented by its unique id which is of type unsigned int32.
type Keyword uint32

// Keywords is a type representing a slice of keywords.
type Keywords []Keyword

// Gets all topics based on keywords, e.g.
// Keywords: A, B, C
// Topics: AB, AC, BC
func (keywords Keywords) GetTopics() Topics {
	// instantiate return object
	topics := Topics{}
	// iterate through keywords
	for i := 0; i < len(keywords); i++ {
		for j := i + 1; j < len(keywords); j++ {
			// instantiate topic
			topic := Object.Topic.New(keywords[i], keywords[j])
			// add it to the return object
			topics = append(topics, topic)
		}
	}
	// return generated topics
	return topics
}

func (keywords Keywords) Clone() Keywords {
	out := []Keyword{}
	for i := 0; i < len(keywords); i++ {
		out = append(out, keywords[i])
	}
	return out
}

// Content object contains information pertinent to content being indexed by ContNet.
type Content struct {
	sync.RWMutex
	ID         int64    // unique content ID
	Keywords   Keywords // slice of keywords for content
	Age        int32    // age of content is given in seconds since publishing
	Confidence float64  // confidence is an estimate about how precisely were keywords determined
	Popularity float64  // Popularity is an estimate about how popular the content is.
}
type ContentFactory struct{}

// Creates new content object
func (factory ContentFactory) New(id int64, keywords Keywords, age int32, confidence float64, popularity float64) *Content {
	return &Content{
		ID:         id,
		Keywords:   keywords.Clone(),
		Age:        age,
		Confidence: confidence,
		Popularity: popularity,
	}
}

// Thread-safe content object clone (deep copy)
func (content *Content) Clone() *Content {
	content.RLock()
	defer content.RUnlock()

	return Object.Content.New(
		content.ID,
		content.Keywords,
		content.Age,
		content.Confidence,
		content.Popularity,
	)
}

// Thread-safe content object update
func (old *Content) Update(new *Content) {
	old.Lock()
	defer old.Unlock()

	old.Keywords = new.Keywords
	old.Age = new.Age
	old.Confidence = new.Confidence
	old.Popularity = new.Popularity
}
