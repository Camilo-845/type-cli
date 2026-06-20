package game

import (
	"math"
	"time"
)

func (g *Game) LiveWPM() float64 {
	total := g.keystrokes.Total()
	if g.State == Idle || total == 0 {
		return 0
	}
	elapsed := time.Since(g.firstKeyTime).Minutes()
	if elapsed < 0.001 {
		elapsed = 0.001
	}
	return float64(g.keystrokes.Correct()) / 5.0 / elapsed
}

func (g *Game) LiveAccuracy() float64 {
	total := g.keystrokes.Total()
	if total == 0 {
		return 100
	}
	return float64(g.keystrokes.Correct()) / float64(total) * 100
}

func (g *Game) CorrectChars() int {
	return g.keystrokes.Correct()
}

func (g *Game) TotalChars() int {
	return g.totalChars
}

func (g *Game) Stats() *Stats {
	elapsed := g.Elapsed.Minutes()
	if elapsed < 0.001 {
		elapsed = 0.001
	}

	correct := g.keystrokes.Correct()
	total := g.keystrokes.Total()

	wpm := float64(correct) / 5.0 / elapsed
	rawWPM := float64(total) / 5.0 / elapsed
	accuracy := 100.0
	if total > 0 {
		accuracy = float64(correct) / float64(total) * 100
	}

	consistency := g.calculateConsistency()

	return &Stats{
		WPM:            math.Round(wpm*10) / 10,
		RawWPM:         math.Round(rawWPM*10) / 10,
		Accuracy:       math.Round(accuracy*10) / 10,
		Consistency:    math.Round(consistency*10) / 10,
		CharsCorrect:   correct,
		CharsIncorrect: g.keystrokes.Incorrect(),
		CharsExtra:     g.keystrokes.Extra(),
		WordsTyped:     g.keystrokes.Words(),
		Duration:       g.Elapsed,
	}
}

func (g *Game) calculateConsistency() float64 {
	bursts := g.bursts.Bursts()
	if len(bursts) < 2 {
		return 100
	}

	var total float64
	for _, k := range bursts {
		total += float64(k)
	}
	mean := total / float64(len(bursts))
	if mean == 0 {
		return 100
	}

	var sumSqDiff float64
	for _, k := range bursts {
		diff := float64(k) - mean
		sumSqDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSqDiff / float64(len(bursts)))
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
