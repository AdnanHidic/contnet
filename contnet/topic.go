package contnet

// Topic is represented by a hash of content keywords.
type Topic uint64
type TopicFactory struct{}

// Returns keywords from topic object. Keywords are always ordered in ASC.
func (t Topic) GetKeywords() (Keyword, Keyword) {
	k1 := t & 0xFFFFFFFF
	k2 := t >> 32

	return Keyword(k1), Keyword(k2)
}

// Creates new Topic from keywords specified.
func (factory TopicFactory) New(k1, k2 Keyword) *Topic {
	// if first keyword is greater than second keyword then swap to maintain consistency.
	if k1 > k2 {
		_swap(&k1, &k2)
	}

	x := (int64(k1)) << 0
	y := (int64(k2)) << 32

	ans := Topic(x | y)
	return &ans
}

// Utility for swapping keywords.
func _swap(k1, k2 *Keyword) {
	tmp := k1
	k1 = k2
	k2 = tmp
}
