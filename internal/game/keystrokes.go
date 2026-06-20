package game

type KeystrokeTracker struct {
	total     int
	correct   int
	incorrect int
	extra     int
	words     int
}

func (k *KeystrokeTracker) RecordCorrect() {
	k.total++
	k.correct++
}

func (k *KeystrokeTracker) RecordIncorrect() {
	k.total++
	k.incorrect++
}

func (k *KeystrokeTracker) RecordExtra() {
	k.total++
	k.incorrect++
	k.extra++
}

func (k *KeystrokeTracker) RecordWord() {
	k.words++
}

func (k *KeystrokeTracker) UndoCorrect() {
	k.total--
	k.correct--
}

func (k *KeystrokeTracker) UndoIncorrect() {
	k.total--
	k.incorrect--
}

func (k *KeystrokeTracker) UndoExtra() {
	k.total--
	k.incorrect--
	k.extra--
}

func (k *KeystrokeTracker) UndoWord() {
	if k.words > 0 {
		k.words--
	}
}

func (k *KeystrokeTracker) Total() int     { return k.total }
func (k *KeystrokeTracker) Correct() int   { return k.correct }
func (k *KeystrokeTracker) Incorrect() int { return k.incorrect }
func (k *KeystrokeTracker) Extra() int     { return k.extra }
func (k *KeystrokeTracker) Words() int     { return k.words }

func (k *KeystrokeTracker) Reset() {
	k.total = 0
	k.correct = 0
	k.incorrect = 0
	k.extra = 0
	k.words = 0
}

func (k *KeystrokeTracker) Snapshot() KeystrokeSnapshot {
	return KeystrokeSnapshot{
		Total:     k.total,
		Correct:   k.correct,
		Incorrect: k.incorrect,
		Extra:     k.extra,
		Words:     k.words,
	}
}

type KeystrokeSnapshot struct {
	Total     int
	Correct   int
	Incorrect int
	Extra     int
	Words     int
}
