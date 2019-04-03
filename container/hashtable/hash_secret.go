package hashtable

import (
	"math/rand"
	"time"
)

type HashSecret struct {
	prefix, suffix int
}

var hashSecret HashSecret

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	hashSecret = HashSecret{
		prefix: rand.Int(),
		suffix: rand.Int(),
	}
}
