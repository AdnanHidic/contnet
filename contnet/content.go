package contnet

import "time"

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

type ID int64

// Content object contains information pertinent to content being indexed by ContNet.
type Content struct {
	ID         ID       // unique content ID
	Keywords   Keywords // slice of keywords for content
	CreatedAt  time.Time
	Confidence float64 // confidence is an estimate about how precisely were keywords determined
	Popularity float64 // Popularity is an estimate about how popular the content is.
}
type ContentFactory struct{}

// Creates new content object
func (factory ContentFactory) New(id ID, keywords Keywords, createdAt time.Time, confidence float64, popularity float64) *Content {
	return &Content{
		ID:         id,
		Keywords:   keywords.Clone(),
		CreatedAt:  createdAt,
		Confidence: confidence,
		Popularity: popularity,
	}
}

// Content object clone (deep copy)
func (content *Content) Clone() *Content {
	return Object.Content.New(
		content.ID,
		content.Keywords,
		content.CreatedAt,
		content.Confidence,
		content.Popularity,
	)
}
