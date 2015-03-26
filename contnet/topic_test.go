package contnet

import (
	"testing"
)

func TestPackUnpack(t *testing.T) {
	keywordPairs := [][]Keyword{}

	for i := 1; i < 100; i++ {
		for j := 1; j < 100; j++ {
			keywordPairs = append(keywordPairs, []Keyword{Keyword(i), Keyword(j)})
		}
	}

	for i := 0; i < len(keywordPairs); i++ {
		keyword1 := keywordPairs[i][0]
		keyword2 := keywordPairs[i][1]

		topic := TopicFactory{}.FromKeywords(keyword1, keyword2)

		keyword1FromTopic, keyword2FromTopic := topic.ToKeywords()

		if keyword1 != keyword1FromTopic {
			t.Errorf("Expected %d but got %d as keyword 1", keyword1, keyword1FromTopic)
		}
		if keyword2 != keyword2FromTopic {
			t.Errorf("Expected %d but got %d as keyword 2", keyword2, keyword2FromTopic)
		}
	}
}
