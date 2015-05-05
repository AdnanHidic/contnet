package contnet

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

var randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

func contains(arr Keywords, k Keyword) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == k {
			return true
		}
	}
	return false
}

func random5Keywords(from Keywords) Keywords {
	keywords := Keywords{}
	for len(keywords) < 5 {
		randomIndex := randomizer.Intn(len(from))
		keyword := from[randomIndex]
		if !contains(keywords, keyword) {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func randomQuality() float64 {
	if randomizer.Intn(2) == 0 {
		return randomizer.Float64()
	} else {
		return -randomizer.Float64()
	}
}

func randomPopularity() float64 {
	return 400.0 * randomizer.Float64()
}

func grade(keywords Keywords, content *Content) float64 {
	kwf, popf, qf := 7.0, 2.0, 1.0
	matchingKeywords := 0
	for i := 0; i < len(content.Keywords); i++ {
		if contains(keywords, content.Keywords[i]) {
			matchingKeywords++
		}
	}
	propF := float64(matchingKeywords) / float64(len(content.Keywords))
	propP := (content.Quality + 1.0) / 2.0
	propQ := content.Popularity / 400.0
	return kwf*propF + popf*propP + qf*propQ

}

func TestNetCreate(t *testing.T) {
	// create n keywords:
	n := 1000
	keywords := Keywords{}
	for i := 0; i < n; i++ {
		keywords = append(keywords, Keyword(i))
	}

	// create m contents with random keyword combinations
	m := 15000
	contents := []*Content{}
	for i := 0; i < m; i++ {
		rK := random5Keywords(keywords)
		rQ := randomQuality()
		rP := randomPopularity()
		// log.Println("K, Q, P: ", rK, rQ, rP)
		content := &Content{
			ID:         ID(i),
			Keywords:   rK,
			CreatedAt:  time.Now(),
			Quality:    rQ,
			Popularity: rP,
		}
		contents = append(contents, content)
	}

	// create user with id = 1 whose real interests are x given keywords
	x := 60
	interestingKeywords := Keywords{}
	for i := 0; i < n; i++ {
		interestingKeywords = append(interestingKeywords, Keyword(randomizer.Intn(x)))
	}

	// grade every generated content based on how it matches user's interests
	grades := map[ID]float64{}
	bestGrade, bestGradeID := 0.0, ID(0)
	worstGrade, worstGradeID := 11.0, ID(0)
	sumGrade := 0.0

	for i := 0; i < m; i++ {
		grades[contents[i].ID] = grade(interestingKeywords, contents[i])
		if grades[contents[i].ID] > grades[bestGradeID] {
			bestGrade = grades[contents[i].ID]
			bestGradeID = contents[i].ID
		}
		if grades[contents[i].ID] < grades[worstGradeID] {
			worstGrade = grades[contents[i].ID]
			worstGradeID = contents[i].ID
		}
		sumGrade += grades[contents[i].ID]
	}

	// index everything
	conf := &NetConfig{
		MaxContentAge:           72 * time.Hour,
		CheckContentAgeInterval: 10 * time.Second,
		ItemsPerPage:            uint8(20),
		NoveltyPct:              0.25,
		SnapshotPath:            "",
		SnapshotInterval:        60 * time.Second,
	}

	net := Object.Net.New(conf)

	for i := 0; i < m; i++ {
		net.SaveContent(contents[i])
	}

	// F user requests for 2 pages of content
	f := 200
	visited := []ID{}
	for i := 0; i < f; i++ {
		// select page
		ids := net.Select(1, 2)
		// mark and select best one to read that was not visited previously
		read := mark(grades, ids, visited)
		// read
		action := &Action{
			ProfileID: 1,
			ContentID: read,
			Type:      ActionTypes.Read,
			Arguments: ActionArguments{
				&ActionArgument{
					Name:  "duration-seconds",
					Type:  ActionArgumentTypes.Integer,
					Value: 10 + randomizer.Intn(160),
				},
			},
		}
		net.SaveAction(action)
		visited = append(visited, read)
	}
	log.Printf("Best grade: %.2f \t Worst grade: %.2f \t Average grade: %.2f ", bestGrade, worstGrade, sumGrade/float64(m))
}

func mark(grades map[ID]float64, ids, visited []ID) ID {
	// filter all visited from ids
	filtered := filterVisited(ids, visited)

	totalGrade := 0.0
	bestID := ID(ids[0])
	for i := 0; i < len(ids); i++ {
		totalGrade += grades[ids[i]]
		if grades[ids[i]] > grades[bestID] {
			bestID = ids[i]
		}
	}

	bestUnvisitedID := ID(filtered[0])
	for i := 0; i < len(filtered); i++ {
		if grades[filtered[i]] > grades[bestUnvisitedID] {
			bestUnvisitedID = filtered[i]
		}
	}

	log.Printf("Total: %.2f \t Best: %.2f \t Best unvisited: %.2f \t Average: %.2f", totalGrade, grades[bestID], grades[bestUnvisitedID], totalGrade/float64(len(ids)))
	return bestUnvisitedID
}

func filterVisited(ids, visited []ID) []ID {
	out := []ID{}

	for i := 0; i < len(ids); i++ {
		if !hasVisited(visited, ids[i]) {
			out = append(out, ids[i])
		}
	}

	return out
}

func hasVisited(ids []ID, id ID) bool {
	for i := 0; i < len(ids); i++ {
		if ids[i] == id {
			return true
		}
	}
	return false
}
