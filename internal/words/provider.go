package words

import (
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

	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = pool[rand.IntN(len(pool))]
	}
	return words
}
