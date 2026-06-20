package words

import (
	"math"
	"math/rand/v2"
)

func zipfRandomIndex(dictLength int) int {
	const gamma = 0.5772156649015329
	hN := math.Log(float64(dictLength)+0.5) + gamma
	r := rand.Float64()
	inverseCDF := math.Exp(r*hN-gamma) - 0.5
	idx := int(math.Floor(inverseCDF))
	if idx < 0 {
		idx = 0
	}
	if idx >= dictLength {
		idx = dictLength - 1
	}
	return idx
}
