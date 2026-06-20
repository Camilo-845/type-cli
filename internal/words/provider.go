package words

import (
	"math"
	"math/rand/v2"
)

var lists = map[string][]string{}

func Register(name string, words []string) {
	lists[name] = words
}

func Names() []string {
	names := make([]string, 0, len(lists))
	for name := range lists {
		names = append(names, name)
	}
	return names
}

func Words(listName string) ([]string, bool) {
	w, ok := lists[listName]
	return w, ok
}

func Generate(count int, listName string) []string {
	pool, ok := lists[listName]
	if !ok {
		if len(lists) == 0 {
			return []string{"the", "quick", "brown", "fox"}
		}
		for _, v := range lists {
			pool = v
			break
		}
	}

	meta, hasMeta := languageMeta[listName]
	zipf := hasMeta && meta.OrderedByFrequency

	words := make([]string, count)
	for i := 0; i < count; i++ {
		if zipf {
			words[i] = pool[zipfIndex(len(pool))]
		} else {
			words[i] = pool[rand.IntN(len(pool))]
		}
	}
	return words
}

func zipfIndex(n int) int {
	u := rand.Float64()
	idx := int(math.Pow(u, 0.85) * float64(n))
	if idx >= n {
		idx = n - 1
	}
	return idx
}
