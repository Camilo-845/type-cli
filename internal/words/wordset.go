package words

import "math/rand/v2"

type Wordset struct {
	Words           []string
	Length          int
	orderedIndex    int
	shuffledIndexes []int
}

func NewWordset(words []string) *Wordset {
	return &Wordset{
		Words:  words,
		Length: len(words),
	}
}

func (w *Wordset) ResetIndexes() {
	w.orderedIndex = 0
	w.shuffledIndexes = nil
}

func (w *Wordset) RandomWord(zipf bool) string {
	if zipf {
		return w.Words[zipfRandomIndex(w.Length)]
	}
	return w.Words[rand.IntN(w.Length)]
}

func (w *Wordset) NextWord() string {
	if w.orderedIndex >= w.Length {
		w.orderedIndex = 0
	}
	word := w.Words[w.orderedIndex]
	w.orderedIndex++
	return word
}

func (w *Wordset) ShuffledWord() string {
	if len(w.shuffledIndexes) == 0 {
		w.generateShuffledIndexes()
	}
	idx := w.shuffledIndexes[len(w.shuffledIndexes)-1]
	w.shuffledIndexes = w.shuffledIndexes[:len(w.shuffledIndexes)-1]
	return w.Words[idx]
}

func (w *Wordset) generateShuffledIndexes() {
	w.shuffledIndexes = make([]int, w.Length)
	for i := 0; i < w.Length; i++ {
		w.shuffledIndexes[i] = i
	}
	rand.Shuffle(w.Length, func(i, j int) {
		w.shuffledIndexes[i], w.shuffledIndexes[j] = w.shuffledIndexes[j], w.shuffledIndexes[i]
	})
}
