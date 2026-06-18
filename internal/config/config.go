package config

type Config struct {
	Mode        string `json:"mode"`
	Duration    int    `json:"duration"`
	WordCount   int    `json:"word_count"`
	WordList    string `json:"word_list"`
}

func DefaultConfig() *Config {
	return &Config{
		Mode:      "time",
		Duration:  30,
		WordCount: 25,
		WordList:  "english",
	}
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

	validLists := map[string]bool{"english": true, "english_1k": true}
	if !validLists[c.WordList] {
		c.WordList = "english"
	}

	return c
}
