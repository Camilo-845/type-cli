package game

import (
	"time"
)

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
