package words

import (
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"
)

const (
	maxRegenAttempts   = 100
	numberProbability  = 0.1
	maxNumberLength    = 4
)

type GeneratorConfig struct {
	Punctuation bool
	LazyMode    bool
	Numbers     bool
}

type Generator struct {
	wordset        *Wordset
	language       *LanguageObject
	cfg            GeneratorConfig
	prevWordRaw    string
	prevWord2Raw   string
	spanishTracker string
	sectionIndex   int
	currentSection []string
	sectionHistory []string
}

func NewGenerator(lang *LanguageObject, cfg GeneratorConfig) *Generator {
	return &Generator{
		wordset:  NewWordset(lang.Words),
		language: lang,
		cfg:      cfg,
	}
}

func (g *Generator) GenerateWords(count int) []string {
	g.wordset.ResetIndexes()
	g.prevWordRaw = ""
	g.prevWord2Raw = ""
	g.spanishTracker = ""
	g.sectionIndex = 0
	g.currentSection = nil
	g.sectionHistory = nil

	words := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var prevWord string
		var prevWord2 string
		if len(words) > 0 {
			prevWord = words[len(words)-1]
		}
		if len(words) > 1 {
			prevWord2 = words[len(words)-2]
		}

		word := g.nextWord(i, count, prevWord, prevWord2)
		words = append(words, word)
	}
	return words
}

func (g *Generator) nextWord(wordIndex, wordsBound int, prevWord, prevWord2 string) string {
	zipf := g.language.OrderedByFrequency

	prevWordRaw := stripForCompare(prevWord)
	prevWord2Raw := stripForCompare2(prevWord2)

	var randomWord string
	if len(g.currentSection) == 0 {
		randomWord = g.wordset.RandomWord(zipf)

		firstAfterSplit := strings.ToLower(firstWord(randomWord))
		firstAfterSplitLazy := g.applyLazyMode(firstAfterSplit)

		regenCount := 0
		for regenCount < maxRegenAttempts &&
			(prevWordRaw == firstAfterSplitLazy ||
				prevWord2Raw == firstAfterSplitLazy ||
				(!g.cfg.Punctuation && randomWord == "I") ||
				(!g.cfg.Punctuation && !g.isCodeLanguage() && hasSymbols(randomWord)) ||
				(!g.cfg.Numbers && hasDigits(randomWord))) {
			regenCount++
			randomWord = g.wordset.RandomWord(zipf)
			firstAfterSplit = strings.ToLower(firstWord(randomWord))
			firstAfterSplitLazy = g.applyLazyMode(firstAfterSplit)
		}

		randomWord = collapseSpaces(randomWord)

		g.currentSection = strings.Split(randomWord, " ")
		g.sectionHistory = append(g.sectionHistory, randomWord)
		randomWord = g.currentSection[0]
		g.currentSection = g.currentSection[1:]
		g.sectionIndex++
	} else {
		randomWord = g.currentSection[0]
		g.currentSection = g.currentSection[1:]
	}

	if randomWord == "" {
		return ""
	}

	langFamily := g.languageFamily()

	shouldLower := !g.cfg.Punctuation &&
		containsUpper(randomWord) &&
		!strings.HasPrefix(langFamily, "german") &&
		!strings.HasPrefix(langFamily, "swiss_german") &&
		!g.isCodeLanguage() &&
		!strings.HasPrefix(langFamily, "klingon")

	if shouldLower {
		randomWord = strings.ToLower(randomWord)
	}

	randomWord = collapseSpaces(randomWord)
	randomWord = g.applyLazyMode(randomWord)

	if strings.HasPrefix(langFamily, "swiss_german") {
		randomWord = strings.ReplaceAll(randomWord, "ß", "ss")
	}

	if g.cfg.Punctuation && !g.language.OriginalPunctuation {
		randomWord = g.punctuateWord(prevWord, randomWord, wordIndex, wordsBound)
	}

	if g.cfg.Numbers && rand.Float64() < numberProbability {
		randomWord = g.generateNumbers(maxNumberLength)
	}

	return randomWord
}

func (g *Generator) applyLazyMode(word string) string {
	if !g.cfg.LazyMode {
		return word
	}
	allow := !g.language.NoLazyMode
	if allow {
		return ReplaceAccents(word, g.language.AdditionalAccents)
	}
	return word
}

func (g *Generator) languageFamily() string {
	lang := g.language.Name
	if idx := strings.Index(lang, "_"); idx != -1 {
		return lang[:idx]
	}
	return lang
}

func (g *Generator) isCodeLanguage() bool {
	return strings.HasPrefix(g.language.Name, "code_")
}

func (g *Generator) generateNumbers(length int) string {
	n := 1 + rand.IntN(length)
	result := make([]byte, n)
	result[0] = byte('1' + rand.IntN(9))
	for i := 1; i < n; i++ {
		result[i] = byte('0' + rand.IntN(10))
	}
	return string(result)
}

func (g *Generator) punctuateWord(prevWord, word string, wordIndex, wordsBound int) string {
	langFamily := g.languageFamily()
	lastChar := lastRune(prevWord)

	if !g.isCodeLanguage() && langFamily != "georgian" &&
		(wordIndex == 0 || shouldCapitalizeAfter(lastChar)) {

		word = capitalizeFirst(word)

		if langFamily == "turkish" {
			word = strings.ReplaceAll(word, "I", "İ")
		}
		if langFamily == "spanish" {
			r := rand.Float64()
			if r > 0.9 {
				word = "¿" + word
				g.spanishTracker = "?"
			} else if r > 0.8 {
				word = "¡" + word
				g.spanishTracker = "!"
			}
		}
		return word
	}

	if (rand.Float64() < 0.1 && lastChar != '.' && lastChar != ',' && wordIndex != wordsBound-2) ||
		wordIndex == wordsBound-1 {

		if langFamily == "spanish" && (g.spanishTracker == "?" || g.spanishTracker == "!") {
			word += g.spanishTracker
			g.spanishTracker = ""
			return word
		}

		r := rand.Float64()
		if r <= 0.8 {
			word += periodForLang(langFamily)
		} else if r > 0.8 && r < 0.9 {
			if langFamily == "french" {
				word = "?"
			} else {
				word += questionForLang(langFamily)
			}
		} else {
			if langFamily == "french" {
				word = "!"
			} else {
				word += exclaimForLang(langFamily)
			}
		}
		return word
	}

	r := rand.Float64()
	switch {
	case r < 0.01 && lastChar != ',' && lastChar != '.' && langFamily != "russian":
		word = `"` + word + `"`
	case r < 0.011 && lastChar != ',' && lastChar != '.' &&
		langFamily != "russian" && langFamily != "ukrainian" && langFamily != "slovak":
		word = `'` + word + `'`
	case r < 0.012 && lastChar != ',' && lastChar != '.':
		word = g.wrapBrackets(word)
	case r < 0.013 && lastChar != ',' && lastChar != '.' &&
		lastChar != ';' && lastChar != '؛' && lastChar != ':' && lastChar != '；' && lastChar != '：':
		word = g.appendColon(word)
	case r < 0.014 && lastChar != ',' && lastChar != '.' && prevWord != "-":
		word = "-"
	case r < 0.015 && lastChar != ',' && lastChar != '.' &&
		lastChar != ';' && lastChar != '؛' && lastChar != '；' && lastChar != '：':
		word = g.appendSemicolon(word)
	case r < 0.2 && lastChar != ',':
		word += commaForLang(langFamily)
	case r < 0.25 && g.isCodeLanguage():
		word = g.codeSpecial()
	case r < 0.5 && langFamily == "english" && englishPunctuationCheck(word):
		word = englishPunctuationReplace(word)
	}

	if strings.ContainsRune(word, '\t') {
		word = strings.ReplaceAll(word, "\t", "")
		word += "\t"
	}
	if strings.ContainsRune(word, '\n') {
		word = strings.ReplaceAll(word, "\n", "")
		word += "\n"
	}

	return word
}

func (g *Generator) wrapBrackets(word string) string {
	if g.isCodeLanguage() {
		brackets := []string{"()", "{}", "[]", "<>"}
		if strings.HasPrefix(g.language.Name, "code_javascript") {
			brackets = append(brackets, "``")
		}
		b := brackets[rand.IntN(len(brackets))]
		runes := []rune(b)
		return string(runes[0]) + word + string(runes[1])
	}
	if g.languageFamily() == "japanese" || g.languageFamily() == "chinese" {
		return "（" + word + "）"
	}
	return "(" + word + ")"
}

func (g *Generator) appendColon(word string) string {
	if g.languageFamily() == "french" {
		return ":"
	}
	if g.languageFamily() == "chinese" {
		return word + "："
	}
	return word + ":"
}

func (g *Generator) appendSemicolon(word string) string {
	switch g.languageFamily() {
	case "french":
		return ";"
	case "greek":
		return "."
	case "arabic", "kurdish":
		return word + "؛"
	case "chinese":
		return word + "；"
	default:
		return word + ";"
	}
}

func (g *Generator) codeSpecial() string {
	specials := []string{"{", "}", "[", "]", "(", ")", ";", "=", "+", "%", "/"}
	specialsC := []string{
		"{", "}", "[", "]", "(", ")", ";", "=", "+", "%", "/",
		"/*", "*/", "//", "!=", "==", "<=", ">=", "||", "&&",
		"<<", ">>", "%=", "&=", "*=", "++", "+=", "--", "-=", "/=", "^=", "|=",
	}
	langName := g.language.Name
	if (strings.HasPrefix(langName, "code_c") && !strings.HasPrefix(langName, "code_css")) ||
		strings.HasPrefix(langName, "code_arduino") {
		return specialsC[rand.IntN(len(specialsC))]
	}
	if strings.HasPrefix(langName, "code_javascript") {
		return append(specials, "`")[rand.IntN(len(specials)+1)]
	}
	return specials[rand.IntN(len(specials))]
}

func stripForCompare(word string) string {
	return strings.ToLower(stripForCompareRe.ReplaceAllString(word, ""))
}

func stripForCompare2(word string) string {
	return strings.ToLower(stripCompare2Re.ReplaceAllString(word, ""))
}

func firstWord(s string) string {
	parts := strings.Split(s, " ")
	if len(parts) > 0 {
		return parts[0]
	}
	return s
}

func collapseSpaces(s string) string {
	return strings.TrimSpace(collapseSpacesRe.ReplaceAllString(s, " "))
}

func lastRune(s string) rune {
	if s == "" {
		return 0
	}
	runes := []rune(s)
	return runes[len(runes)-1]
}

func shouldCapitalizeAfter(last rune) bool {
	return last == '.' || last == '!' || last == '?' || last == '؟'
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func containsUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

var (
	symbolRe          = regexp.MustCompile(`[-=_+\[\]{};'\\:"|,./<>?]`)
	digitRe           = regexp.MustCompile(`[0-9]`)
	collapseSpacesRe  = regexp.MustCompile(` +`)
	stripForCompareRe = regexp.MustCompile(`[.?!":\-,]`)
	stripCompare2Re   = regexp.MustCompile(`[.?!":\-,']`)
)

func hasSymbols(s string) bool {
	return symbolRe.MatchString(s)
}

func hasDigits(s string) bool {
	return digitRe.MatchString(s)
}

func periodForLang(lang string) string {
	switch lang {
	case "nepali", "bangla", "hindi":
		return "।"
	case "japanese", "chinese":
		return "。"
	default:
		return "."
	}
}

func questionForLang(lang string) string {
	switch lang {
	case "arabic", "persian", "urdu", "kurdish":
		return "؟"
	case "greek":
		return ";"
	case "japanese", "chinese":
		return "？"
	default:
		return "?"
	}
}

func exclaimForLang(lang string) string {
	switch lang {
	case "japanese", "chinese":
		return "！"
	default:
		return "!"
	}
}

func commaForLang(lang string) string {
	switch lang {
	case "arabic", "urdu", "persian", "kurdish":
		return "،"
	case "japanese":
		return "、"
	case "chinese":
		return "，"
	default:
		return ","
	}
}
