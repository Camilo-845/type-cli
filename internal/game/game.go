package game

import (
	"time"
)

type State int

const (
	Idle State = iota
	Running
	Paused
	Complete
)

type WordState struct {
	Word    string
	Typed   string
	Correct []bool
}

type Game struct {
	Words    []WordState
	Current  int
	State    State
	Elapsed  time.Duration

	totalKeystrokes     int
	correctKeystrokes   int
	incorrectKeystrokes int
	extraKeystrokes     int
	completedWords      int
	firstKeyTime        time.Time

	TimeMode  bool
	Duration  time.Duration
	WordCount int

	burstKeystrokes []int
	burstCorrect    []int
	lastBurstTick   time.Time
	pendingKeystrokes int
	pendingCorrect    int
}

type Stats struct {
	WPM             float64
	RawWPM          float64
	Accuracy        float64
	Consistency     float64
	CharsCorrect    int
	CharsIncorrect  int
	CharsExtra      int
	WordsTyped      int
	Duration        time.Duration
}

const wordPoolSize = 500

func NewTimeGame(durationSec int, words []string) *Game {
	batch := make([]WordState, wordPoolSize)
	for i, w := range words {
		if i >= wordPoolSize {
			break
		}
		batch[i] = WordState{Word: w, Typed: "", Correct: nil}
	}
	return &Game{
		Words:      batch,
		Current:    0,
		State:      Idle,
		TimeMode:   true,
		Duration:   time.Duration(durationSec) * time.Second,
		WordCount:  wordPoolSize,
	}
}

func NewWordGame(count int, words []string) *Game {
	batch := make([]WordState, count)
	for i, w := range words {
		if i >= count {
			break
		}
		batch[i] = WordState{Word: w, Typed: "", Correct: nil}
	}
	return &Game{
		Words:      batch,
		Current:    0,
		State:      Idle,
		TimeMode:   false,
		WordCount:  count,
	}
}

func (g *Game) CurrentWord() *WordState {
	if g.Current >= len(g.Words) {
		return nil
	}
	return &g.Words[g.Current]
}

func (g *Game) Remaining() time.Duration {
	if !g.TimeMode {
		return 0
	}
	remaining := g.Duration - g.Elapsed
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}
