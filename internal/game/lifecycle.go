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

	g.bursts.Flush()

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
	g.firstKeyTime = time.Time{}
	g.keystrokes.Reset()
	g.bursts.Reset()

	for i := range g.Words {
		g.Words[i].Typed = nil
		g.Words[i].Correct = nil
	}
}
