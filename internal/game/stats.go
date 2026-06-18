package game

import (
	"math"
	"time"
)

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
		WPM:            math.Round(wpm*10) / 10,
		RawWPM:         math.Round(rawWPM*10) / 10,
		Accuracy:       math.Round(accuracy*10) / 10,
		Consistency:    math.Round(consistency*10) / 10,
		CharsCorrect:   g.correctKeystrokes,
		CharsIncorrect: g.incorrectKeystrokes,
		CharsExtra:     g.extraKeystrokes,
		WordsTyped:     g.completedWords,
		Duration:       g.Elapsed,
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
