package game

import (
	"fmt"
	"time"
)

type Result struct {
	Date           time.Time     `json:"date"`
	Mode           string        `json:"mode"`
	Duration       int           `json:"duration,omitempty"`
	WordCount      int           `json:"word_count,omitempty"`
	WordList       string        `json:"word_list"`
	WPM            float64       `json:"wpm"`
	RawWPM         float64       `json:"raw_wpm"`
	Accuracy       float64       `json:"accuracy"`
	Consistency    float64       `json:"consistency"`
	CharsCorrect   int           `json:"chars_correct"`
	CharsIncorrect int           `json:"chars_incorrect"`
	CharsExtra     int           `json:"chars_extra"`
	WordsTyped     int           `json:"words_typed"`
	TestDuration   time.Duration `json:"test_duration"`
}

func NewResult(stats *Stats, mode string, duration int, wordCount int, wordList string) Result {
	return Result{
		Date:           time.Now(),
		Mode:           mode,
		Duration:       duration,
		WordCount:      wordCount,
		WordList:       wordList,
		WPM:            stats.WPM,
		RawWPM:         stats.RawWPM,
		Accuracy:       stats.Accuracy,
		Consistency:    stats.Consistency,
		CharsCorrect:   stats.CharsCorrect,
		CharsIncorrect: stats.CharsIncorrect,
		CharsExtra:     stats.CharsExtra,
		WordsTyped:     stats.WordsTyped,
		TestDuration:   stats.Duration,
	}
}

func (r Result) ModeText() string {
	if r.Mode == "time" {
		return fmt.Sprintf("%ds", r.Duration)
	}
	return fmt.Sprintf("%d words", r.WordCount)
}
