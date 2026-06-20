package words

import "unicode"

type accentPair = [2]string

var commonAccents = []accentPair{
	{"áàâäåãąą́āą̄ă", "a"},
	{"éèêëẽęę́ēę̄ėě", "e"},
	{"íìîïĩįį́īį̄ı", "i"},
	{"óòôöøõóōǫǫ́ǭő", "o"},
	{"úùûüŭũúūůű", "u"},
	{"ńň", "n"},
	{"çĉčć", "c"},
	{"řŕṛ", "r"},
	{"ďđḍ", "d"},
	{"ťțṭ", "t"},
	{"ṃ", "m"},
	{"æ", "ae"},
	{"œ", "oe"},
	{"ẅŵ", "w"},
	{"ĝğg̃", "g"},
	{"ĥ", "h"},
	{"ĵ", "j"},
	{"ńṇṅ", "n"},
	{"ŝśšșşṣ", "s"},
	{"ß", "ss"},
	{"żźž", "z"},
	{"ÿỹýÿŷ", "y"},
	{"łľĺ", "l"},
	{"أإآ", "ا"},
	{"َ", ""},
	{"ُ", ""},
	{"ِ", ""},
	{"ْ", ""},
	{"ً", ""},
	{"ٌ", ""},
	{"ٍ", ""},
	{"ّ", ""},
	{"ё", "е"},
	{"ά", "α"},
	{"έ", "ε"},
	{"ί", "ι"},
	{"ύ", "υ"},
	{"ό", "ο"},
	{"ή", "η"},
	{"ώ", "ω"},
	{"þ", "th"},
}

var commonAccentsMap map[rune]string

func init() {
	commonAccentsMap = make(map[rune]string)
	for _, pair := range commonAccents {
		for _, r := range pair[0] {
			commonAccentsMap[r] = pair[1]
		}
	}
}

func buildAdditionalMap(additional [][2]string) map[rune]string {
	if len(additional) == 0 {
		return nil
	}
	m := make(map[rune]string)
	for _, pair := range additional {
		for _, r := range pair[0] {
			m[r] = pair[1]
		}
	}
	return m
}

func findAccent(r rune, additionalMap map[rune]string) (string, bool) {
	lower := unicode.ToLower(r)
	if additionalMap != nil {
		if repl, ok := additionalMap[lower]; ok {
			return repl, true
		}
	}
	repl, ok := commonAccentsMap[lower]
	return repl, ok
}

func ReplaceAccents(word string, additional [][2]string) string {
	if word == "" {
		return word
	}

	additionalMap := buildAdditionalMap(additional)

	runes := []rune(word)
	uppercase := make([]bool, len(runes))
	for i, r := range runes {
		uppercase[i] = unicode.IsUpper(r)
	}

	var result []rune
	offset := 0
	for i := 0; i < len(runes); i++ {
		idx := i + offset
		if idx >= len(runes) {
			break
		}

		replacement, found := findAccent(runes[idx], additionalMap)
		if found {
			replRunes := []rune(replacement)
			for j, cr := range replRunes {
				caseIdx := idx + j
				if caseIdx < len(uppercase) && uppercase[caseIdx] {
					result = append(result, unicode.ToUpper(cr))
				} else {
					result = append(result, cr)
				}
			}
			offset += len([]rune(string(runes[idx]))) - 1
		} else {
			if uppercase[idx] {
				result = append(result, unicode.ToUpper(runes[idx]))
			} else {
				result = append(result, runes[idx])
			}
		}
	}

	return string(result)
}
