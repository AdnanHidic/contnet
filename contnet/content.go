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
	ID         ID
	Keywords   Keywords
	CreatedAt  time.Time
	Quality    float64
	Popularity float64
}
type ContentFactory struct{}

// Creates new content object
func (factory ContentFactory) New(id ID, keywords Keywords, createdAt time.Time, quality, popularity float64) *Content {
	return &Content{
		ID:         id,
		Keywords:   keywords.Clone(),
		CreatedAt:  createdAt,
		Quality:    quality,
		Popularity: popularity,
	}
}

// Content object clone (deep copy)
func (content *Content) Clone() *Content {
	return Object.Content.New(
		content.ID,
		content.Keywords,
		content.CreatedAt,
		content.Quality,
		content.Popularity,
	)
}
