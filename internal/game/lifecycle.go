package game

import (
	"time"
)

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
		g.Words[i].Typed = nil
		g.Words[i].Correct = nil
	}
}
