package game

type BurstTracker struct {
	bursts        []int
	correctBursts []int
	pending       int
	pendingCorrect int
}

func (b *BurstTracker) AddPending(correct bool) {
	b.pending++
	if correct {
		b.pendingCorrect++
	}
}

func (b *BurstTracker) UndoPending(correct bool) {
	b.pending--
	if correct {
		b.pendingCorrect--
	}
}

func (b *BurstTracker) Flush() {
	b.bursts = append(b.bursts, b.pending)
	b.correctBursts = append(b.correctBursts, b.pendingCorrect)
	b.pending = 0
	b.pendingCorrect = 0
}

func (b *BurstTracker) Bursts() []int        { return b.bursts }
func (b *BurstTracker) CorrectBursts() []int { return b.correctBursts }

func (b *BurstTracker) Reset() {
	b.bursts = nil
	b.correctBursts = nil
	b.pending = 0
	b.pendingCorrect = 0
}
