package game

import (
	"testing"
	"time"
)

func TestHandleKey_IdleToRunning(t *testing.T) {
	g := NewTimeGame(30, []string{"hello", "world"})
	if g.State != Idle {
		t.Fatal("expected idle")
	}

	g.HandleKey("h")
	if g.State != Running {
		t.Error("expected running after first keystroke")
	}
}

func TestHandleKey_CorrectChar(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def"})
	g.HandleKey("a")

	ws := g.CurrentWord()
	if ws.Typed != "a" {
		t.Errorf("expected 'a', got %q", ws.Typed)
	}
	if !ws.Correct[0] {
		t.Error("expected correct")
	}
	if g.correctKeystrokes != 1 {
		t.Errorf("expected 1 correct, got %d", g.correctKeystrokes)
	}
}

func TestHandleKey_IncorrectChar(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def"})
	g.HandleKey("z")

	ws := g.CurrentWord()
	if ws.Correct[0] {
		t.Error("expected incorrect")
	}
	if g.incorrectKeystrokes != 1 {
		t.Errorf("expected 1 incorrect, got %d", g.incorrectKeystrokes)
	}
}

func TestHandleKey_SpaceAdvancesWord(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def", "ghi"})
	g.HandleKey("a")
	g.HandleKey("b")
	g.HandleKey("c")
	g.HandleKey(" ")

	if g.Current != 1 {
		t.Errorf("expected current 1, got %d", g.Current)
	}
	if g.completedWords != 1 {
		t.Errorf("expected 1 completed, got %d", g.completedWords)
	}
}

func TestHandleKey_Backspace(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def"})
	g.HandleKey("a")
	g.HandleKey("x")
	g.HandleKey("backspace")

	ws := g.CurrentWord()
	if ws.Typed != "a" {
		t.Errorf("expected 'a', got %q", ws.Typed)
	}
	if len(ws.Correct) != 1 {
		t.Errorf("expected 1 correct bool, got %d", len(ws.Correct))
	}
}

func TestHandleKey_ExtraChars(t *testing.T) {
	g := NewTimeGame(30, []string{"a", "b"})
	g.HandleKey("a")
	g.HandleKey("b")
	g.HandleKey("c")

	ws := g.CurrentWord()
	if len(ws.Typed) != 3 {
		t.Errorf("expected 3 typed chars, got %d", len(ws.Typed))
	}
	if g.extraKeystrokes != 2 {
		t.Errorf("expected 2 extra, got %d", g.extraKeystrokes)
	}
}

func TestWordGame_CompletesAtCount(t *testing.T) {
	g := NewWordGame(2, []string{"a", "b"})
	g.HandleKey("a")
	g.HandleKey(" ")
	g.HandleKey("b")
	g.HandleKey(" ")

	if g.State != Complete {
		t.Error("expected complete after all words")
	}
}

func TestTick_TimeExpires(t *testing.T) {
	g := NewTimeGame(1, []string{"hello", "world"})
	g.HandleKey("h")
	g.State = Running
	g.firstKeyTime = time.Now().Add(-2 * time.Second)
	g.Elapsed = 2 * time.Second
	g.Tick()

	if g.State != Complete {
		t.Error("expected complete when time expires")
	}
}

func TestLiveAccuracy(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def"})
	if g.LiveAccuracy() != 100 {
		t.Errorf("expected 100%%, got %.1f", g.LiveAccuracy())
	}

	g.HandleKey("a")
	g.HandleKey("x")
	if g.LiveAccuracy() != 50 {
		t.Errorf("expected 50%%, got %.1f", g.LiveAccuracy())
	}
}

func TestPauseResume(t *testing.T) {
	g := NewTimeGame(30, []string{"abc"})
	g.HandleKey("a")
	g.Pause()
	if g.State != Paused {
		t.Error("expected paused")
	}

	g.HandleKey("b")
	if len(g.CurrentWord().Typed) != 1 {
		t.Error("should not accept input while paused")
	}

	g.Resume()
	if g.State != Running {
		t.Error("expected running after resume")
	}

	g.HandleKey("b")
	if len(g.CurrentWord().Typed) != 2 {
		t.Error("should accept input after resume")
	}
}

func TestReset(t *testing.T) {
	g := NewTimeGame(30, []string{"abc", "def"})
	g.HandleKey("a")
	g.HandleKey(" ")

	g.Reset()
	if g.State != Idle {
		t.Error("expected idle after reset")
	}
	if g.Current != 0 {
		t.Error("expected current 0 after reset")
	}
	if g.totalKeystrokes != 0 {
		t.Error("expected 0 keystrokes after reset")
	}
}

func TestStats(t *testing.T) {
	t.Run("idle stats are zero", func(t *testing.T) {
		g := NewTimeGame(30, []string{"hello"})
		s := g.Stats()
		if s.WPM != 0 {
			t.Error("expected 0 wpm")
		}
		if s.Accuracy != 100 {
			t.Errorf("expected 100%% accuracy, got %.1f", s.Accuracy)
		}
	})

	t.Run("stats after typing", func(t *testing.T) {
		g := NewTimeGame(30, []string{"abc", "def"})
		g.HandleKey("a")
		g.HandleKey("b")
		g.HandleKey("x")
		g.HandleKey(" ")
		g.Elapsed = time.Minute
		s := g.Stats()
		if s.WordsTyped != 1 {
			t.Errorf("expected 1 word, got %d", s.WordsTyped)
		}
		if s.Accuracy <= 0 {
			t.Error("expected positive accuracy")
		}
	})
}
