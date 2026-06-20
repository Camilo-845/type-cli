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
	}

	ws := g.CurrentWord()
	if ws == nil {
		return
	}

	if char == ' ' {
		g.submitWord()
		return
	}

	ws.Typed = append(ws.Typed, char)
	ws.Correct = append(ws.Correct, false)

	idx := len(ws.Typed) - 1
	if idx < len(ws.Word) {
		ws.Correct[idx] = ws.Typed[idx] == ws.Word[idx]
	} else {
		ws.Correct[idx] = false
	}

	if idx >= len(ws.Word) {
		g.keystrokes.RecordExtra()
		g.bursts.AddPending(false)
	} else if ws.Correct[idx] {
		g.keystrokes.RecordCorrect()
		g.bursts.AddPending(true)
	} else {
		g.keystrokes.RecordIncorrect()
		g.bursts.AddPending(false)
	}

	if g.Current == g.WordCount-1 && len(ws.Typed) == len(ws.Word) {
		g.keystrokes.RecordWord()
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

	runes := []rune(s)
	if len(runes) == 1 {
		g.handleKeystroke(runes[0])
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
		if g.Current > 0 {
			g.Current--
			g.keystrokes.UndoWord()
		}
		return
	}

	idx := len(ws.Typed) - 1
	wasCorrect := ws.Correct[idx]
	wasExtra := idx >= len(ws.Word)

	ws.Typed = ws.Typed[:len(ws.Typed)-1]
	ws.Correct = ws.Correct[:len(ws.Correct)-1]

	if wasExtra {
		g.keystrokes.UndoExtra()
	} else if wasCorrect {
		g.keystrokes.UndoCorrect()
	} else {
		g.keystrokes.UndoIncorrect()
	}
	g.bursts.UndoPending(wasCorrect)
}

func (g *Game) submitWord() {
	if g.State != Running {
		return
	}

	g.keystrokes.RecordWord()

	g.Current++
	if g.Current >= len(g.Words) {
		g.State = Complete
		return
	}

	if !g.TimeMode && g.Current >= g.WordCount {
		g.State = Complete
	}
}
