package game

import (
	"math"
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

func (g *Game) handleKeystroke(char rune) {
	if g.State == Complete || g.State == Paused {
		return
	}

	if g.State == Idle {
		g.State = Running
		g.firstKeyTime = time.Now()
		g.lastBurstTick = time.Now()
	}

	ws := g.CurrentWord()
	if ws == nil {
		return
	}

	if char == ' ' {
		g.submitWord()
		return
	}

	ws.Typed += string(char)
	ws.Correct = append(ws.Correct, false)

	idx := len(ws.Typed) - 1
	if idx < len(ws.Word) {
		ws.Correct[idx] = ws.Typed[idx] == ws.Word[idx]
	} else {
		ws.Correct[idx] = false
	}

	if ws.Correct[idx] {
		g.correctKeystrokes++
	} else {
		g.incorrectKeystrokes++
	}
	if idx >= len(ws.Word) {
		g.extraKeystrokes++
	}

	g.totalKeystrokes++
	g.pendingKeystrokes++
	if ws.Correct[idx] {
		g.pendingCorrect++
	}

	if g.Current == g.WordCount-1 && len(ws.Typed) == len(ws.Word) {
		g.completedWords++
		g.Current++
		g.State = Complete
	}
}

func (g *Game) HandleKey(s string) {
	if s == "" {
		return
	}

	if s == "backspace" {
		g.handleBackspace()
		return
	}

	if s == " " {
		g.submitWord()
		return
	}

	if len(s) == 1 {
		g.handleKeystroke(rune(s[0]))
	}
}

func (g *Game) handleBackspace() {
	if g.State != Running {
		return
	}

	ws := g.CurrentWord()
	if ws == nil {
		return
	}

	if len(ws.Typed) == 0 {
		return
	}

	idx := len(ws.Typed) - 1
	wasCorrect := ws.Correct[idx]

	ws.Typed = ws.Typed[:len(ws.Typed)-1]
	ws.Correct = ws.Correct[:len(ws.Correct)-1]

	g.totalKeystrokes--
	g.pendingKeystrokes--
	if wasCorrect {
		g.correctKeystrokes--
		g.pendingCorrect--
	} else {
		g.incorrectKeystrokes--
		if idx >= len(ws.Word) {
			g.extraKeystrokes--
		}
	}
}

func (g *Game) submitWord() {
	if g.State != Running {
		return
	}

	g.completedWords++

	g.Current++
	if g.Current >= len(g.Words) {
		g.State = Complete
		return
	}

	if !g.TimeMode && g.Current >= g.WordCount {
		g.State = Complete
	}
}

func (g *Game) Tick() {
	if g.State != Running {
		return
	}

	now := time.Now()
	g.Elapsed = now.Sub(g.firstKeyTime)

	g.burstKeystrokes = append(g.burstKeystrokes, g.pendingKeystrokes)
	g.burstCorrect = append(g.burstCorrect, g.pendingCorrect)
	g.pendingKeystrokes = 0
	g.pendingCorrect = 0
	g.lastBurstTick = now

	if g.TimeMode && g.Elapsed >= g.Duration {
		g.State = Complete
	}
}

func (g *Game) Pause() {
	if g.State == Running {
		g.State = Paused
	}
}

func (g *Game) Resume() {
	if g.State == Paused {
		g.State = Running
	}
}

func (g *Game) Reset() {
	g.State = Idle
	g.Current = 0
	g.Elapsed = 0
	g.totalKeystrokes = 0
	g.correctKeystrokes = 0
	g.incorrectKeystrokes = 0
	g.extraKeystrokes = 0
	g.completedWords = 0
	g.burstKeystrokes = nil
	g.burstCorrect = nil
	g.pendingKeystrokes = 0
	g.pendingCorrect = 0

	for i := range g.Words {
		g.Words[i].Typed = ""
		g.Words[i].Correct = nil
	}
}

func (g *Game) LiveWPM() float64 {
	if g.State == Idle || g.totalKeystrokes == 0 {
		return 0
	}
	elapsed := time.Since(g.firstKeyTime).Minutes()
	if elapsed < 0.001 {
		elapsed = 0.001
	}
	return float64(g.correctKeystrokes) / 5.0 / elapsed
}

func (g *Game) LiveAccuracy() float64 {
	if g.totalKeystrokes == 0 {
		return 100
	}
	return float64(g.correctKeystrokes) / float64(g.totalKeystrokes) * 100
}

func (g *Game) CorrectChars() int {
	return g.correctKeystrokes
}

func (g *Game) TotalChars() int {
	total := 0
	for _, ws := range g.Words {
		total += len(ws.Word)
	}
	return total
}

func (g *Game) Stats() *Stats {
	elapsed := g.Elapsed.Minutes()
	if elapsed < 0.001 {
		elapsed = 0.001
	}

	wpm := float64(g.correctKeystrokes) / 5.0 / elapsed
	rawWPM := float64(g.totalKeystrokes) / 5.0 / elapsed
	accuracy := 100.0
	if g.totalKeystrokes > 0 {
		accuracy = float64(g.correctKeystrokes) / float64(g.totalKeystrokes) * 100
	}

	consistency := g.calculateConsistency()

	return &Stats{
		WPM:        math.Round(wpm*10) / 10,
		RawWPM:     math.Round(rawWPM*10) / 10,
		Accuracy:   math.Round(accuracy*10) / 10,
		Consistency: math.Round(consistency*10) / 10,
		CharsCorrect: g.correctKeystrokes,
		CharsIncorrect: g.incorrectKeystrokes,
		CharsExtra: g.extraKeystrokes,
		WordsTyped: g.completedWords,
		Duration:   g.Elapsed,
	}
}

func (g *Game) calculateConsistency() float64 {
	if len(g.burstKeystrokes) < 2 {
		return 100
	}

	var total float64
	for _, k := range g.burstKeystrokes {
		total += float64(k)
	}
	mean := total / float64(len(g.burstKeystrokes))
	if mean == 0 {
		return 100
	}

	var sumSqDiff float64
	for _, k := range g.burstKeystrokes {
		diff := float64(k) - mean
		sumSqDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSqDiff / float64(len(g.burstKeystrokes)))
	cv := stdDev / mean

	consistency := (1 - cv) * 100
	if consistency < 0 {
		consistency = 0
	}
	if consistency > 100 {
		consistency = 100
	}
	return consistency
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
