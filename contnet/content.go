package contnet

type Keyword uint32
type Keywords []Keyword

func (keywords Keywords) Clone() Keywords {
	out := []Keyword{}
	for i := 0; i < len(keywords); i++ {
		out = append(out, keywords[i])
	}
	return out
}

type Content struct {
	ID         int64
	Keywords   Keywords
	Age        int32
	Confidence float64
	Relevance  float64
}
type ContentFactory struct{}

func (factory ContentFactory) New(id int64, keywords Keywords, age int32, confidence float64, relevance float64) *Content {
	return &Content{
		ID:         id,
		Keywords:   keywords.Clone(),
		Age:        age,
		Confidence: confidence,
		Relevance:  relevance,
	}
}

func (content *Content) Clone() *Content {
	return Object.Content.New(
		content.ID,
		content.Keywords,
		content.Age,
		content.Confidence,
		content.Relevance,
	)
}
