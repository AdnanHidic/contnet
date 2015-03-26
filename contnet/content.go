package contnet

import "sync"

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

func NewContent(id int64, keywords []Keyword, age int32, confidence float64, relevance float64) *Content {
	return &Content{
		ID:         id,
		Keywords:   keywords,
		Age:        age,
		Confidence: confidence,
		Relevance:  relevance,
	}
}

func (content *Content) Clone() *Content {
	return &Content{
		ID:         content.ID,
		Keywords:   content.Keywords.Clone(),
		Age:        content.Age,
		Confidence: content.Confidence,
		Relevance:  content.Relevance,
	}
}

type ContentStorage struct {
	sync.RWMutex
	content map[int64]*Content
}

func NewContentStorage() *ContentStorage {
	return &ContentStorage{
		content: map[int64]*Content{},
	}
}

func (storage *ContentStorage) Get(id int64) *Content {
	storage.RLock()
	defer storage.RUnlock()

	if content, ok := storage.content[id]; !ok {
		return nil
	} else {
		return content.Clone()
	}
}

func (storage *ContentStorage) Save(content *Content) {
	storage.Lock()
	defer storage.Unlock()

	storage.content[content.ID] = content.Clone()
}

func (storage *ContentStorage) Delete(id int64) {
	storage.Lock()
	defer storage.Unlock()

	delete(storage.content, id)
}
