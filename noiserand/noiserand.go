// Package noiserand implements a random number generator that is based on multiple dimensional noise, using one of the
// dimensions as a seed, which allows slight alterations to the seed to slightly alter the generated sequence.
package noiserand

import (
	"math/rand"
)

const defaultRngSeed = 123456

// NoiseFunc is a function that takes two float64 values and returns a float64 value.
type NoiseFunc func(x, y float64) float64

// NoiseRand is a random number generator that is based on multidimensional noise, using one of the
// dimensions as a seed, which allows slight alterations to the seed to slightly alter the generated sequence.
type NoiseRand struct {
	seed   float64    // The offset used as a dimension in the noise function.
	offset float64    // Offset used to jump ahead in the noise function sequence.
	rng    *rand.Rand // This is used for randomly jumping ahead in the sequence with a somewhat deterministic result.
	noise  NoiseFunc  // The noise function used to generate the random numbers.
}

// New creates a new NoiseRand instance with the given seed and noise function.
func New(seed float64, f NoiseFunc) *NoiseRand {
	return &NoiseRand{
		seed:  seed,
		rng:   rand.New(rand.NewSource(defaultRngSeed)),
		noise: f,
	}
}

// SetSeed sets the seed used for the random number generator that
// is used to modify the values of the random numbers.
func (r *NoiseRand) SetSeed(seed float64) {
	r.seed = seed
}

// SetRngSeed sets the seed used for the random number generator that
// is used to jump ahead in the noise sequence.
//
// NOTE: This is NOT the seed you should rely on to alter the noise sequence.
func (r *NoiseRand) SetRngSeed(seed int64) {
	r.rng.Seed(seed)
}

// JumpAhead jumps ahead in the random number sequence by the given amount.
func (r *NoiseRand) JumpAhead(amount float64) {
	r.offset += amount
}

// Intn returns a random integer in the range [0, n).
func (r *NoiseRand) Intn(n int) int {
	nv := r.noise(r.seed, r.offset)
	r.JumpAhead(r.rng.Float64())
	return int(nv * float64(n))
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *NoiseRand) Int63() int64 {
	nv := r.noise(r.seed, r.offset)
	r.JumpAhead(r.rng.Float64())
	return int64(nv * float64(1<<63))
}

// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
func (r *NoiseRand) Float64() float64 {
	nv := r.noise(r.seed, r.offset)
	r.JumpAhead(r.rng.Float64())
	return nv
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
func (r *NoiseRand) Perm(n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}
