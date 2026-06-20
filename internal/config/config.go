package config

import "github.com/Camilo-845/type-cli/internal/words"

type Config struct {
	Mode      string `json:"mode"`
	Duration  int    `json:"duration"`
	WordCount int    `json:"word_count"`
	WordList  string `json:"word_list"`
}

func DefaultConfig() *Config {
	return &Config{
		Mode:      "time",
		Duration:  30,
		WordCount: 25,
		WordList:  "english",
	}
}

func AvailableLanguages() []string {
	return words.ListLanguages()
}

func (c *Config) Validate() *Config {
	if c == nil {
		return DefaultConfig()
	}

	validModes := map[string]bool{"time": true, "words": true}
	if !validModes[c.Mode] {
		c.Mode = "time"
	}

	validDurations := map[int]bool{15: true, 30: true, 60: true, 120: true}
	if !validDurations[c.Duration] {
		c.Duration = 30
	}

	validCounts := map[int]bool{10: true, 25: true, 50: true, 100: true}
	if !validCounts[c.WordCount] {
		c.WordCount = 25
	}

	available := make(map[string]bool)
	for _, lang := range words.ListLanguages() {
		available[lang] = true
	}
	if len(available) == 0 {
		available["english"] = true
	}
	if !available[c.WordList] {
		c.WordList = "english"
	}

	return c
}

func (c *Config) ToggleMode() {
	if c.Mode == "time" {
		c.Mode = "words"
	} else {
		c.Mode = "time"
	}
}

var durations = []int{15, 30, 60, 120}

func (c *Config) CycleDuration(forward bool) {
	n := len(durations)
	for i, d := range durations {
		if d == c.Duration {
			if forward {
				c.Duration = durations[(i+1)%n]
			} else {
				c.Duration = durations[(i-1+n)%n]
			}
			return
		}
	}
}

var counts = []int{10, 25, 50, 100}

func (c *Config) CycleWordCount(forward bool) {
	n := len(counts)
	for i, ct := range counts {
		if ct == c.WordCount {
			if forward {
				c.WordCount = counts[(i+1)%n]
			} else {
				c.WordCount = counts[(i-1+n)%n]
			}
			return
		}
	}
}

func (c *Config) CycleWordList(forward bool) {
	names := words.Names()
	if len(names) == 0 {
		return
	}
	n := len(names)
	for i, l := range names {
		if l == c.WordList {
			if forward {
				c.WordList = names[(i+1)%n]
			} else {
				c.WordList = names[(i-1+n)%n]
			}
			return
		}
	}
	c.WordList = names[0]
}

func (c *Config) Apply(cursor int, forward bool) {
	switch cursor {
	case 0:
		c.ToggleMode()
	case 1:
		if c.Mode == "time" {
			c.CycleDuration(forward)
		} else {
			c.CycleWordCount(forward)
		}
	case 2:
		c.CycleWordList(forward)
	}
}
