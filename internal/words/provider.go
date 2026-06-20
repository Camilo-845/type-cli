package words

import "slices"

var lists = map[string][]string{}

func Register(name string, words []string) {
	lists[name] = words
}

func SortedNames() []string {
	names := make([]string, 0, len(lists))
	for name := range lists {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func Names() []string {
	return SortedNames()
}

func Words(listName string) ([]string, bool) {
	w, ok := lists[listName]
	return w, ok
}

func Generate(count int, listName string, cfg GeneratorConfig) []string {
	meta, ok := languageMeta[listName]
	if !ok {
		if len(lists) == 0 {
			return []string{"the", "quick", "brown", "fox"}
		}
		for _, v := range lists {
			meta = LanguageObject{Name: listName, Words: v}
			break
		}
	}

	gen := NewGenerator(&meta, cfg)
	return gen.GenerateWords(count)
}
