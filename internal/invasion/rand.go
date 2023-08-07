package invasion

import (
	"math/rand"
	"time"
)

// rnd is package-level random number generator.
var rnd Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Rand is interface to random number generator.
type Rand interface {
	Intn(n int) int
}
