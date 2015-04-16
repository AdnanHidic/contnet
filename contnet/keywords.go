package contnet

import (
	"log"
	"math"
	"sort"
	"strings"
)

type keywordExtractor struct{}

var ContentKeywordExtractor = keywordExtractor{}

type ContentKeywordExtractionInput struct {
	Title       string
	Description string
	Comments    []string
}

func (keywordExtractor keywordExtractor) Extract(input *ContentKeywordExtractionInput, maxKeywords int) []string {
	joinedComments := strings.Join(input.Comments, " ")

	// get histograms for every element of the input
	titleHistogram, titleCache := __stringHistogram(&input.Title)
	descriptionHistogram, descriptionCache := __stringHistogram(&input.Description)
	commentHistogram, commentCache := __stringHistogram(&joinedComments)

	// now we have histograms (sorted in descending order by word occurrence count) for all elements of the input
	var descriptionAvg, commentAvg float64
	titleHistogram, titleCache, _ = __removeOutliers(titleHistogram, titleCache)
	descriptionHistogram, descriptionCache, descriptionAvg = __removeOutliers(descriptionHistogram, descriptionCache)
	commentHistogram, commentCache, commentAvg = __removeOutliers(commentHistogram, commentCache)

	// now all common or rare words have been removed. Remaining words are all candidates for keyword status.
	// if a word occurs in: title, description and comment it is the strongest candidate
	// combination of a two is still good
	// only in one is pretty bad but still can be keyword
	// join maps
	global := map[string]float64{}
	for tK, tV := range titleCache {
		global[tK] = float64(tV) + descriptionAvg + commentAvg
	}

	for dK, dV := range descriptionCache {
		global[dK] += float64(dV) + commentAvg
	}

	for cK, cV := range commentCache {
		global[cK] += float64(cV)
	}

	// extract array and sort
	globalArr := []*WordCount{}
	for k, v := range global {
		globalArr = append(globalArr, &WordCount{Word: k, Count: v})
	}

	WordCountBy(wordCountCriteria).Sort(globalArr)

	lb := maxKeywords
	if maxKeywords > len(globalArr) {
		lb = len(globalArr) - 1
	}

	out := []string{}

	for i := 0; i < lb; i++ {
		out = append(out, globalArr[i].Word)
	}

	return out
}

func __removeOutliers(histogram []*WordCount, cache map[string]float64) ([]*WordCount, map[string]float64, float64) {
	countSum := 0.0
	for i := 0; i < len(histogram); i++ {
		countSum += histogram[i].Count
	}

	meanValue := countSum / float64(len(histogram))

	sqDiffSum := 0.0
	for i := 0; i < len(histogram); i++ {
		sqDiffSum += math.Pow(histogram[i].Count-meanValue, 2.0)
	}

	variance := sqDiffSum / float64(len(histogram))

	std := math.Sqrt(variance)

	out := []*WordCount{}
	outCache := map[string]float64{}

	lb := meanValue - std
	ub := meanValue + std

	tmpSum := 0.0
	for i := 0; i < len(histogram); i++ {
		if histogram[i].Count >= lb && histogram[i].Count <= ub && !__simplewords.is(histogram[i].Word) {
			out = append(out, histogram[i])
			outCache[histogram[i].Word] = histogram[i].Count
			tmpSum += histogram[i].Count
		} else {
			log.Printf("Removing word %s occ = %f", histogram[i].Word, histogram[i].Count)
		}
	}

	return out, outCache, tmpSum / float64(len(out))
}

type runes map[rune]bool
type simplewords map[string]bool

func (r runes) contain(rune rune) bool {
	_, contains := r[rune]
	return contains
}

func (sw simplewords) is(word string) bool {
	_, contains := sw[word]
	return contains
}

var __stopRunes = runes{
	' ':  true,
	'.':  true,
	'\n': true,
	'(':  true,
	')':  true,
	'"':  true,
	',':  true,
	'*':  true,
}

var __simplewords = simplewords{
	"": true,
}

func __stringHistogram(str *string) ([]*WordCount, map[string]float64) {
	words := strings.FieldsFunc(*str, func(r rune) bool {
		return __stopRunes.contain(r)
	})

	cache := map[string]float64{}

	for i := 0; i < len(words); i++ {
		cache[words[i]]++
	}

	out := []*WordCount{}

	for k, v := range cache {
		out = append(out, &WordCount{Word: k, Count: v})
	}

	WordCountBy(wordCountCriteria).Sort(out)

	return out, cache
}

type WordCount struct {
	Word  string
	Count float64
}

var wordCountCriteria = func(c1, c2 *WordCount) bool {
	return c1.Count > c2.Count
}

// function that defines ordering between content objects
type WordCountBy func(c1, c2 *WordCount) bool

// method on the function type, sorts the argument slice according  to the function
func (wordCountBy WordCountBy) Sort(wordCounts []*WordCount) {
	ws := &wordCountSorter{
		wordCounts:  wordCounts,
		wordCountBy: wordCountBy,
	}
	sort.Sort(ws)
}

type wordCountSorter struct {
	wordCounts  []*WordCount
	wordCountBy func(c1, c2 *WordCount) bool
}

// Len is part of sort.Interface.
func (ws *wordCountSorter) Len() int {
	return len(ws.wordCounts)
}

// Swap is part of sort.Interface.
func (ws *wordCountSorter) Swap(i, j int) {
	ws.wordCounts[i], ws.wordCounts[j] = ws.wordCounts[j], ws.wordCounts[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ws *wordCountSorter) Less(i, j int) bool {
	return ws.wordCountBy(ws.wordCounts[i], ws.wordCounts[j])
}
