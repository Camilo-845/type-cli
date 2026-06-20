package words

import (
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"
)

var (
	punctCheckRe   = regexp.MustCompile(`^\W*(\w+)\W*$`)
	punctReplaceRe = regexp.MustCompile(`^(\W*)(\w+)(\W*)$`)
)

type contractionGroup struct {
	base string
	cont []string
}

var contractionGroups = []contractionGroup{
	{"are", []string{"aren't"}},
	{"can", []string{"can't"}},
	{"could", []string{"couldn't"}},
	{"did", []string{"didn't"}},
	{"does", []string{"doesn't"}},
	{"do", []string{"don't"}},
	{"had", []string{"hadn't"}},
	{"has", []string{"hasn't"}},
	{"have", []string{"haven't"}},
	{"is", []string{"isn't"}},
	{"it", []string{"it's", "it'll"}},
	{"i", []string{"i'm", "i'll", "i've", "i'd"}},
	{"you", []string{"you'll", "you're", "you've", "you'd"}},
	{"that", []string{"that's", "that'll", "that'd"}},
	{"must", []string{"mustn't", "must've"}},
	{"there", []string{"there's", "there'll", "there'd"}},
	{"he", []string{"he's", "he'll", "he'd"}},
	{"she", []string{"she's", "she'll", "she'd"}},
	{"we", []string{"we're", "we'll", "we'd"}},
	{"they", []string{"they're", "they'll", "they'd"}},
	{"should", []string{"shouldn't", "should've"}},
	{"was", []string{"wasn't"}},
	{"were", []string{"weren't"}},
	{"will", []string{"won't"}},
	{"would", []string{"wouldn't", "would've"}},
	{"going", []string{"goin'"}},
}

func englishPunctuationCheck(word string) bool {
	lower := strings.ToLower(word)
	matches := punctCheckRe.FindStringSubmatch(lower)
	if matches == nil {
		return false
	}
	baseWord := matches[1]
	for _, group := range contractionGroups {
		if baseWord == group.base {
			return true
		}
	}
	return false
}

func englishPunctuationReplace(word string) string {
	lower := strings.ToLower(word)
	matches := punctReplaceRe.FindStringSubmatch(lower)
	if matches == nil {
		return word
	}
	prefix := matches[1]
	baseWord := matches[2]
	suffix := matches[3]

	for _, group := range contractionGroups {
		if baseWord == group.base {
			replacement := group.cont[rand.IntN(len(group.cont))]

			origBase := findOriginalBase(word, group.base)
			cased := applyCase(origBase, replacement)
			return prefix + cased + suffix
		}
	}
	return word
}

func findOriginalBase(word, base string) string {
	re := regexp.MustCompile("(?i)^(\\W*)(" + regexp.QuoteMeta(base) + ")(\\W*)$")
	matches := re.FindStringSubmatch(word)
	if matches == nil {
		return base
	}
	return matches[2]
}

func applyCase(original, replacement string) string {
	if original == "I" {
		return replacement
	}
	isAllUpper := true
	for _, r := range original {
		if !unicode.IsUpper(r) {
			isAllUpper = false
			break
		}
	}
	if isAllUpper {
		return strings.ToUpper(replacement)
	}
	firstUpper := unicode.IsUpper([]rune(original)[0])
	if firstUpper {
		runes := []rune(replacement)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}
	return replacement
}
