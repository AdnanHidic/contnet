package contnet

type Topic uint64

type TopicFactory struct{}

func (t Topic) ToKeywords() (Keyword, Keyword) {
	k1 := t & 0xFFFFFFFF
	k2 := t >> 32

	return Keyword(k1), Keyword(k2)
}

func (factory TopicFactory) FromKeywords(k1, k2 Keyword) Topic {
	if k1 > k2 {
		_swap(&k1, &k2)
	}

	x := (int64(k1)) << 0
	y := (int64(k2)) << 32

	ans := x | y
	return Topic(ans)
}

func _swap(k1, k2 *Keyword) {
	tmp := k1
	k1 = k2
	k2 = tmp
}
