package dfstyle

import (
	"math"
	"unsafe"
)

type NoiseType int

const (
	NoiseTypePerlin NoiseType = 1 << iota
	NoiseTypeSimplex
	NoiseTypeWavelet
	NoiseTypeDefault NoiseType = 0
)

const (
	TCOD_NOISE_MAX_OCTAVES        = 128
	TCOD_NOISE_MAX_DIMENSIONS     = 4
	TCOD_NOISE_DEFAULT_HURST      = 0.5
	TCOD_NOISE_DEFAULT_LACUNARITY = 2.0
)

type Noise struct {
	ndim            int
	Map             [256]uint8
	buffer          [256][TCOD_NOISE_MAX_DIMENSIONS]float64
	H               float64
	lacunarity      float64
	exponent        [TCOD_NOISE_MAX_OCTAVES]float64
	waveletTileData []float64
	noiseType       NoiseType
}

const WAVELET_TILE_SIZE = 32
const WAVELET_ARAD = 16

const SIMPLEX_SCALE = 0.5
const WAVELET_SCALE = 2.0

// Common noise function pointer.
// Right now `TCOD_noise_wavelet` prevents `noise` from being const.
type TCOD_noise_func_t func(noise *Noise, f []float64) float64

// Return a floating point value clamped between -1.0f and 1.0f exclusively.
// The return value excludes -1.0f and 1.0f to avoid rounding issues.

func clamp_signed_f(value float64) float64 {
	const LOW = -1.0 + math.SmallestNonzeroFloat64
	const HIGH = 1.0 - math.SmallestNonzeroFloat64
	if value < LOW {
		return LOW
	}
	if value > HIGH {
		return HIGH
	}
	return value
}

func lattice(data *Noise, ix int, fx float64, iy int, fy float64, iz int, fz float64, iw int, fw float64) float64 {
	n := [4]int{ix, iy, iz, iw}
	f := [4]float64{fx, fy, fz, fw}
	nIndex := 0
	for i := 0; i < data.ndim; i++ {
		nIndex = int(data.Map[(nIndex+n[i])&0xFF])
	}
	value := float64(0)
	for i := 0; i < data.ndim; i++ {
		value += data.buffer[nIndex][i] * f[i]
	}
	return value
}

const DEFAULT_SEED uint32 = 0x15687436
const DELTA float64 = 1e-6

func FLOOR(a float64) int {
	if a > 0 {
		return int(a)
	}
	return int(a) - 1
}

func CUBIC(a float64) float64 {
	return a * a * (3 - 2*a)
}

func GENERIC_SWAP(x, y interface{}) {
	x, y = y, x
}

func normalize(data *Noise, f []float64) {
	magnitude := float64(0)
	for i := 0; i < data.ndim; i++ {
		magnitude += f[i] * f[i]
	}
	magnitude = 1.0 / float64(math.Sqrt(float64(magnitude)))
	for i := 0; i < data.ndim; i++ {
		f[i] *= magnitude
	}
}

func TCOD_noise_new(ndim int, hurst, lacunarity float64) *Noise {
	data := &Noise{}
	data.ndim = ndim
	for i := 0; i < 256; i++ {
		data.Map[i] = uint8(i)
		for j := 0; j < data.ndim; j++ {
			data.buffer[i][j] = TCOD_random_get_float(-0.5, 0.5)
		}
		normalize(data, data.buffer[i][:])
	}
	for i := 255; i >= 0; i-- {
		j := TCOD_random_get_int(0, 255)
		GENERIC_SWAP(&data.Map[i], &data.Map[j])
	}
	data.H = hurst
	data.lacunarity = lacunarity
	f := float64(1)
	for i := 0; i < TCOD_NOISE_MAX_OCTAVES; i++ {
		data.exponent[i] = 1.0 / f
		f *= lacunarity
	}
	data.noiseType = NoiseTypePerlin

	return data
}

func cubic(a float64) float64 {
	return a * a * (3 - 2*a)
}

func lerp(a, b, x float64) float64 {
	return a + x*(b-a)
}

func noisePerlin(data *Noise, f []float64) float64 {
	n := [TCOD_NOISE_MAX_DIMENSIONS]int{}
	r := [TCOD_NOISE_MAX_DIMENSIONS]float64{}
	w := [TCOD_NOISE_MAX_DIMENSIONS]float64{}
	for i := 0; i < data.ndim; i++ {
		n[i] = int(math.Floor(float64(f[i])))
		r[i] = f[i] - float64(n[i])
		w[i] = cubic(r[i])
	}

	var value float64
	switch data.ndim {
	case 1:
		value = lerp(
			lattice(data, n[0], r[0], 0, 0, 0, 0, 0, 0),
			lattice(data, n[0]+1, r[0]-1, 0, 0, 0, 0, 0, 0),
			w[0],
		)
	case 2:
		value = lerp(
			lerp(
				lattice(data, n[0], r[0], n[1], r[1], 0, 0, 0, 0),
				lattice(data, n[0]+1, r[0]-1, n[1], r[1], 0, 0, 0, 0),
				w[0],
			),
			lerp(
				lattice(data, n[0], r[0], n[1]+1, r[1]-1, 0, 0, 0, 0),
				lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, 0, 0, 0, 0),
				w[0],
			),
			w[1],
		)
	case 3:
		value = lerp(
			lerp(
				lerp(
					lattice(data, n[0], r[0], n[1], r[1], n[2], r[2], 0, 0),
					lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2], r[2], 0, 0),
					w[0],
				),
				lerp(
					lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2], r[2], 0, 0),
					lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2], r[2], 0, 0),
					w[0],
				),
				w[1],
			),
			lerp(
				lerp(
					lattice(data, n[0], r[0], n[1], r[1], n[2]+1, r[2]-1, 0, 0),
					lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2]+1, r[2]-1, 0, 0),
					w[0],
				),
				lerp(
					lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2]+1, r[2]-1, 0, 0),
					lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2]+1, r[2]-1, 0, 0),
					w[0],
				),
				w[1],
			),
			w[2],
		)
	case 4:
		value = lerp(
			lerp(
				lerp(
					lerp(
						lattice(data, n[0], r[0], n[1], r[1], n[2], r[2], n[3], r[3]),
						lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2], r[2], n[3], r[3]),
						w[0],
					),
					lerp(
						lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2], r[2], n[3], r[3]),
						lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2], r[2], n[3], r[3]),
						w[0],
					),
					w[1],
				),
				lerp(
					lerp(
						lattice(data, n[0], r[0], n[1], r[1], n[2]+1, r[2]-1, n[3], r[3]),
						lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2]+1, r[2]-1, n[3], r[3]),
						w[0],
					),
					lerp(
						lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2]+1, r[2]-1, 0, 0),
						lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2]+1, r[2]-1, n[3], r[3]),
						w[0],
					),
					w[1],
				),
				w[2],
			),
			lerp(
				lerp(
					lerp(
						lattice(data, n[0], r[0], n[1], r[1], n[2], r[2], n[3]+1, r[3]-1),
						lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2], r[2], n[3]+1, r[3]-1),
						w[0],
					),
					lerp(
						lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2], r[2], n[3]+1, r[3]-1),
						lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2], r[2], n[3]+1, r[3]-1),
						w[0],
					),
					w[1],
				),
				lerp(
					lerp(
						lattice(data, n[0], r[0], n[1], r[1], n[2]+1, r[2]-1, n[3]+1, r[3]-1),
						lattice(data, n[0]+1, r[0]-1, n[1], r[1], n[2]+1, r[2]-1, n[3]+1, r[3]-1),
						w[0],
					),
					lerp(
						lattice(data, n[0], r[0], n[1]+1, r[1]-1, n[2]+1, r[2]-1, 0, 0),
						lattice(data, n[0]+1, r[0]-1, n[1]+1, r[1]-1, n[2]+1, r[2]-1, n[3]+1, r[3]-1),
						w[0],
					),
					w[1],
				),
				w[2],
			),
			w[3],
		)
	}
	return clampSignedF(value)
}

func clampSignedF(x float64) float64 {
	if x < -1.0 {
		return -1.0
	}
	if x > 1.0 {
		return 1.0
	}
	return x
}
func absmod(x, n int) int {
	m := x % n
	if m < 0 {
		m += n
	}
	return m
}

// simplex noise, adapted from Ken Perlin's presentation at Siggraph 2001
// and Stefan Gustavson implementation
func simplexGradient1D(h int, x float64) float64 {
	h &= 0xF
	grad := 1.0 + float64(h&7)
	if h&8 != 0 {
		grad = -grad
	}
	return grad * x
}

func simplexGradient2D(h int, x, y float64) float64 {
	var u, v float64
	h &= 0x7
	if h < 4 {
		u = x
		v = 2.0 * y
	} else {
		u = y
		v = 2.0 * x
	}
	var n float64
	if h&1 != 0 {
		n = -u
	}
	if h&2 != 0 {
		n += -v
	}
	// n := ((h & 1) != 0) * -u
	// n += ((h & 2) != 0) * -v
	return n
}

func simplexGradient3D(h int, x, y, z float64) float64 {
	h &= 0xF
	u := x
	v := y
	if h < 4 {
		v = 2.0 * y
	} else if h == 12 || h == 14 {
		u = y
		v = x
	} else {
		u = z
		v = y
	}
	var n float64
	if h&1 != 0 {
		n = -u
	}
	if h&2 != 0 {
		n += -v
	}
	// n := ((h & 1) != 0) * -u
	// n += ((h & 2) != 0) * -v
	return n
}

func simplexGradient4D(h int, x, y, z, t float64) float64 {
	h &= 0x1F
	u := x
	v := y
	w := z
	if h < 16 {
		v = y
		w = z
	} else if h < 24 {
		u = y
		w = t
	} else {
		u = z
		v = t
	}
	var n float64
	if h&1 != 0 {
		n = -u
	}
	if h&2 != 0 {
		n += -v
	}
	if h&4 != 0 {
		n += -w
	}
	// n := ((h & 1) != 0) * -u
	// n += ((h & 2) != 0) * -v
	// n += ((h & 4) != 0) * -w
	return n
}

var simplex = [64][4]float64{
	{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 0, 0, 0}, {0, 2, 3, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {1, 2, 3, 0},
	{0, 2, 1, 3}, {0, 0, 0, 0}, {0, 3, 1, 2}, {0, 3, 2, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {1, 3, 2, 0},
	{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0},
	{1, 2, 0, 3}, {0, 0, 0, 0}, {1, 3, 0, 2}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {2, 3, 0, 1}, {2, 3, 1, 0},
	{1, 0, 2, 3}, {1, 0, 3, 2}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {2, 0, 3, 1}, {0, 0, 0, 0}, {2, 1, 3, 0},
	{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0},
	{2, 0, 1, 3}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {3, 0, 1, 2}, {3, 0, 2, 1}, {0, 0, 0, 0}, {3, 1, 2, 0},
	{2, 1, 0, 3}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {3, 1, 0, 2}, {0, 0, 0, 0}, {3, 2, 0, 1}, {3, 2, 1, 0},
}

func noiseSimplex(data *Noise, f []float64) float64 {
	switch data.ndim {
	case 1:
		i0 := int(math.Floor(float64(f[0] * SIMPLEX_SCALE)))
		i1 := i0 + 1
		x0 := f[0]*SIMPLEX_SCALE - float64(i0)
		x1 := x0 - 1.0
		t0 := 1.0 - x0*x0
		t1 := 1.0 - x1*x1
		t0 *= t0
		t1 *= t1
		i0 = int(data.Map[i0&0xFF])
		var n0 float64
		n0 = simplexGradient1D(i0, x0)
		n0 *= t0 * t0
		i1 = int(data.Map[i1&0xFF])
		var n1 float64
		n1 = simplexGradient1D(i1, x1)
		n1 *= t1 * t1
		return clampSignedF(0.25 * (float64(n0) + float64(n1)))
	case 2:
		const F2 float64 = 0.366025403 // 0.5 * (sqrt(3.0)-1.0)
		const G2 float64 = 0.211324865 // (3.0 - sqrt(3.0))/6.0
		s := (f[0] + f[1]) * F2 * SIMPLEX_SCALE
		xs := f[0]*SIMPLEX_SCALE + s
		ys := f[1]*SIMPLEX_SCALE + s
		i := int(math.Floor(float64(xs)))
		j := int(math.Floor(float64(ys)))
		t := float64(i+j) * G2
		xo := float64(i) - t
		yo := float64(j) - t
		x0 := f[0]*SIMPLEX_SCALE - xo
		y0 := f[1]*SIMPLEX_SCALE - yo
		ii := absmod(i, 256)
		jj := absmod(j, 256)
		var i1, j1 int
		if x0 > y0 {
			i1 = 1
			j1 = 0
		} else {
			i1 = 0
			j1 = 1
		}
		x1 := x0 - float64(i1) + G2
		y1 := y0 - float64(j1) + G2
		x2 := x0 - 1.0 + 2.0*G2
		y2 := y0 - 1.0 + 2.0*G2
		t0 := 0.5 - x0*x0 - y0*y0
		var n0 float64
		if t0 >= 0.0 {
			idx := (ii + int(data.Map[jj])) & 0xFF
			t0 *= t0
			idx = int(data.Map[idx])
			n0 = simplexGradient2D(idx, x0, y0)
			n0 *= t0 * t0
		}
		t1 := 0.5 - x1*x1 - y1*y1
		var n1 float64
		if t1 >= 0.0 {
			idx := (ii + i1 + int(data.Map[(jj+j1)&0xFF])) & 0xFF
			t1 *= t1
			idx = int(data.Map[idx])
			n1 = simplexGradient2D(idx, x1, y1)
			n1 *= t1 * t1
		}
		t2 := 0.5 - x2*x2 - y2*y2
		var n2 float64
		if t2 >= 0.0 {
			idx := (ii + 1 + int(data.Map[(jj+1)&0xFF])) & 0xFF
			t2 *= t2
			idx = int(data.Map[idx])
			n2 = simplexGradient2D(idx, x2, y2)
			n2 *= t2 * t2
		}
		return clampSignedF(40.0 * (float64(n0) + float64(n1) + float64(n2)))
	case 3:
		const F3 = 0.333333333
		const G3 = 0.166666667
		s := (f[0] + f[1] + f[2]) * F3 * SIMPLEX_SCALE
		xs := f[0]*SIMPLEX_SCALE + s
		ys := f[1]*SIMPLEX_SCALE + s
		zs := f[2]*SIMPLEX_SCALE + s
		i := int(math.Floor(float64(xs)))
		j := int(math.Floor(float64(ys)))
		k := int(math.Floor(float64(zs)))
		t := float64(i+j+k) * G3
		xo := float64(i) - t
		yo := float64(j) - t
		zo := float64(k) - t
		x0 := f[0]*SIMPLEX_SCALE - xo
		y0 := f[1]*SIMPLEX_SCALE - yo
		z0 := f[2]*SIMPLEX_SCALE - zo
		var i1, j1, k1, i2, j2, k2 int
		if x0 >= y0 {
			if y0 >= z0 {
				i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 1, 0
			} else if x0 >= z0 {
				i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 0, 1
			} else {
				i1, j1, k1, i2, j2, k2 = 0, 0, 1, 1, 0, 1
			}
		} else {
			if y0 < z0 {
				i1, j1, k1, i2, j2, k2 = 0, 0, 1, 0, 1, 1
			} else if x0 < z0 {
				i1, j1, k1, i2, j2, k2 = 0, 1, 0, 0, 1, 1
			} else {
				i1, j1, k1, i2, j2, k2 = 0, 1, 0, 1, 1, 0
			}
		}
		x1 := x0 - float64(i1) + G3
		y1 := y0 - float64(j1) + G3
		z1 := z0 - float64(k1) + G3
		x2 := x0 - float64(i2) + 2.0*G3
		y2 := y0 - float64(j2) + 2.0*G3
		z2 := z0 - float64(k2) + 2.0*G3
		x3 := x0 - 1.0 + 3.0*G3
		y3 := y0 - 1.0 + 3.0*G3
		z3 := z0 - 1.0 + 3.0*G3
		ii := absmod(i, 256)
		jj := absmod(j, 256)
		kk := absmod(k, 256)
		t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
		n0 := float64(0)
		if t0 >= 0 {
			idx := data.Map[(ii+int(data.Map[(jj+int(data.Map[kk]))&0xFF]))&0xFF]
			t0 *= t0
			n0 = simplexGradient3D(int(idx), x0, y0, z0)
			n0 *= t0 * t0
		}
		t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
		n1 := float64(0)
		if t1 >= 0 {
			idx := data.Map[(ii+i1+int(data.Map[(jj+j1+int(data.Map[(kk+k1)&0xFF]))&0xFF]))&0xFF]
			t1 *= t1
			n1 = simplexGradient3D(int(idx), x1, y1, z1)
			n1 *= t1 * t1
		}
		t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
		n2 := float64(0)
		if t2 >= 0 {
			idx := data.Map[(ii+i2+int(data.Map[(jj+j2+int(data.Map[(kk+k2)&0xFF]))&0xFF]))&0xFF]
			t2 *= t2
			n2 = simplexGradient3D(int(idx), x2, y2, z2)
			n2 *= t2 * t2
		}
		t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
		n3 := float64(0)
		if t3 >= 0 {
			idx := data.Map[(ii+1+int(data.Map[(jj+1+int(data.Map[(kk+1)&0xFF]))&0xFF]))&0xFF]
			t3 *= t3
			n3 = simplexGradient3D(int(idx), x3, y3, z3)
			n3 *= t3 * t3
		}
		return clampSignedF(32.0 * (n0 + n1 + n2 + n3))
	case 4:
		const F4 float64 = 0.309016994 // (sqrtf(5.0f)-1.0f)/4.0f
		const G4 float64 = 0.138196601 // (5.0f - sqrtf(5.0f))/20.0f
		s := (f[0] + f[1] + f[2] + f[3]) * F4 * SIMPLEX_SCALE
		xs := f[0]*SIMPLEX_SCALE + s
		ys := f[1]*SIMPLEX_SCALE + s
		zs := f[2]*SIMPLEX_SCALE + s
		ws := f[3]*SIMPLEX_SCALE + s
		i := int(math.Floor(float64(xs)))
		j := int(math.Floor(float64(ys)))
		k := int(math.Floor(float64(zs)))
		l := int(math.Floor(float64(ws)))
		t := float64(i+j+k+l) * G4
		xo := float64(i) - t
		yo := float64(j) - t
		zo := float64(k) - t
		wo := float64(l) - t
		x0 := f[0]*SIMPLEX_SCALE - xo
		y0 := f[1]*SIMPLEX_SCALE - yo
		z0 := f[2]*SIMPLEX_SCALE - zo
		w0 := f[3]*SIMPLEX_SCALE - wo
		var c1, c2, c3, c4, c5, c6 int
		if x0 > y0 {
			c1 = 32
		}
		if x0 > z0 {
			c2 = 16
		}
		if y0 > z0 {
			c3 = 8
		}
		if x0 > w0 {
			c4 = 4
		}
		if y0 > w0 {
			c5 = 2
		}
		if z0 > w0 {
			c6 = 1
		}
		c := c1 + c2 + c3 + c4 + c5 + c6
		i1 := 0
		j1 := 0
		k1 := 0
		l1 := 0
		if simplex[c][0] >= 3 {
			i1 = 1
		}
		if simplex[c][1] >= 3 {
			j1 = 1
		}
		if simplex[c][2] >= 3 {
			k1 = 1
		}
		if simplex[c][3] >= 3 {
			l1 = 1
		}
		i2 := 0
		j2 := 0
		k2 := 0
		l2 := 0
		if simplex[c][0] >= 2 {
			i2 = 1
		}
		if simplex[c][1] >= 2 {
			j2 = 1
		}
		if simplex[c][2] >= 2 {
			k2 = 1
		}
		if simplex[c][3] >= 2 {
			l2 = 1
		}
		i3 := 0
		j3 := 0
		k3 := 0
		l3 := 0
		if simplex[c][0] >= 1 {
			i3 = 1
		}
		if simplex[c][1] >= 1 {
			j3 = 1
		}
		if simplex[c][2] >= 1 {
			k3 = 1
		}
		if simplex[c][3] >= 1 {
			l3 = 1
		}
		x1 := x0 - float64(i1) + G4
		y1 := y0 - float64(j1) + G4
		z1 := z0 - float64(k1) + G4
		w1 := w0 - float64(l1) + G4
		x2 := x0 - float64(i2) + 2.0*G4
		y2 := y0 - float64(j2) + 2.0*G4
		z2 := z0 - float64(k2) + 2.0*G4
		w2 := w0 - float64(l2) + 2.0*G4
		x3 := x0 - float64(i3) + 3.0*G4
		y3 := y0 - float64(j3) + 3.0*G4
		z3 := z0 - float64(k3) + 3.0*G4
		w3 := w0 - float64(l3) + 3.0*G4
		x4 := x0 - 1.0 + 4.0*G4
		y4 := y0 - 1.0 + 4.0*G4
		z4 := z0 - 1.0 + 4.0*G4
		w4 := w0 - 1.0 + 4.0*G4
		ii := absmod(i, 256)
		jj := absmod(j, 256)
		kk := absmod(k, 256)
		ll := absmod(l, 256)
		t0 := 0.6 - x0*x0 - y0*y0 - z0*z0 - w0*w0
		n0 := float64(0)
		if t0 >= 0 {
			idx := data.Map[(ii+int(data.Map[(jj+int(data.Map[(kk+int(data.Map[ll&0xFF]))&0xFF]))&0xFF]))&0xFF]
			t0 *= t0
			n0 = simplexGradient4D(int(idx), x0, y0, z0, w0)
			n0 *= t0 * t0
		}
		t1 := 0.6 - x1*x1 - y1*y1 - z1*z1 - w1*w1
		n1 := float64(0)
		if t1 >= 0 {
			idx := data.Map[(ii+i1+int(data.Map[(jj+j1+int(data.Map[(kk+k1+int(data.Map[(ll+l1)&0xFF]))&0xFF]))&0xFF]))&0xFF]
			t1 *= t1
			n1 = simplexGradient4D(int(idx), x1, y1, z1, w1)
			n1 *= t1 * t1
		}
		t2 := 0.6 - x2*x2 - y2*y2 - z2*z2 - w2*w2
		n2 := float64(0)
		if t2 >= 0 {
			idx := data.Map[(ii+i2+int(data.Map[(jj+j2+int(data.Map[(kk+k2+int(data.Map[(ll+l2)&0xFF]))&0xFF]))&0xFF]))&0xFF]
			t2 *= t2
			n2 = simplexGradient4D(int(idx), x2, y2, z2, w2)
			n2 *= t2 * t2
		}
		t3 := 0.6 - x3*x3 - y3*y3 - z3*z3 - w3*w3
		n3 := float64(0)
		if t3 >= 0 {
			idx := data.Map[(ii+i3+int(data.Map[(jj+j3+int(data.Map[(kk+k3+int(data.Map[(ll+l3)&0xFF]))&0xFF]))&0xFF]))&0xFF]
			t3 *= t3
			n3 = simplexGradient4D(int(idx), x3, y3, z3, w3)
			n3 *= t3 * t3
		}
		t4 := 0.6 - x4*x4 - y4*y4 - z4*z4 - w4*w4
		n4 := float64(0)
		if t4 >= 0 {
			idx := data.Map[(ii+1+int(data.Map[(jj+1+int(data.Map[(kk+1+int(data.Map[(ll+1)&0xFF]))&0xFF]))&0xFF]))&0xFF]
			t4 *= t4
			n4 = simplexGradient4D(int(idx), x4, y4, z4, w4)
			n4 *= t4 * t4
		}
		return clampSignedF(27.0 * (float64(n0) + float64(n1) + float64(n2) + float64(n3) + float64(n4)))
	default:
		return math.NaN()
	}
	return 0.0
}

func noiseFbmInt(noise *Noise, f []float64, octaves float64, fn func(*Noise, []float64) float64) float64 {
	tf := make([]float64, noise.ndim)
	copy(tf, f)
	value := float64(0)
	for i := float64(0); i < octaves; i++ {
		value += fn(noise, tf) * noise.exponent[int(i)]
		for j := 0; j < noise.ndim; j++ {
			tf[j] *= noise.lacunarity
		}
	}
	if octaves-float64(int(octaves)) > DELTA {
		value += (octaves - float64(int(octaves))) * fn(noise, tf) * noise.exponent[int(octaves)]
	}
	return clampSignedF(float64(value))
}

func noiseFbmPerlin(noise *Noise, f []float64, octaves float64) float64 {
	return noiseFbmInt(noise, f, octaves, noisePerlin)
}

func noiseFbmSimplex(noise *Noise, f []float64, octaves float64) float64 {
	return noiseFbmInt(noise, f, octaves, noiseSimplex)
}

func noiseTurbulenceInt(noise *Noise, f []float64, octaves float64, fn func(*Noise, []float64) float64) float64 {
	tf := make([]float64, noise.ndim)
	copy(tf, f)
	value := float64(0)
	for i := float64(0); i < octaves; i++ {
		noiseValue := fn(noise, tf)
		value += math.Abs(noiseValue) * noise.exponent[int(i)]
		for j := 0; j < noise.ndim; j++ {
			tf[j] *= noise.lacunarity
		}
	}
	if octaves-float64(int(octaves)) > DELTA {
		noiseValue := fn(noise, tf)
		value += (octaves - float64(int(octaves))) * math.Abs(noiseValue) * noise.exponent[int(octaves)]
	}
	return clampSignedF(float64(value))
}

func noiseTurbulencePerlin(noise *Noise, f []float64, octaves float64) float64 {
	return noiseTurbulenceInt(noise, f, octaves, noisePerlin)
}

func noiseTurbulenceSimplex(noise *Noise, f []float64, octaves float64) float64 {
	return noiseTurbulenceInt(noise, f, octaves, noiseSimplex)
}

// wavelet noise, adapted from Robert L. Cook and Tony Derose 'Wavelet noise' paper
func noiseWaveletDownsample(from []float64, to []float64, stride int) {
	aCoefficients := []float64{
		0.000334, -0.001528, 0.000410, 0.003545, -0.000938, -0.008233, 0.002172, 0.019120,
		-0.005040, -0.044412, 0.011655, 0.103311, -0.025936, -0.243780, 0.033979, 0.655340,
		0.655340, 0.033979, -0.243780, -0.025936, 0.103311, 0.011655, -0.044412, -0.005040,
		0.019120, 0.002172, -0.008233, -0.000938, 0.003546, 0.000410, -0.001528, 0.000334,
	}
	a := aCoefficients[WAVELET_ARAD:]
	for i := 0; i < WAVELET_TILE_SIZE/2; i++ {
		to[i*stride] = 0
		for k := 2*i - WAVELET_ARAD; k < 2*i+WAVELET_ARAD; k++ {
			to[i*stride] += a[k-2*i] * from[absmod(k, WAVELET_TILE_SIZE)*stride]
		}
	}
}

func noiseWaveletUpsample(from []float64, to []float64, stride int) {
	pCoefficient := []float64{0.25, 0.75, 0.75, 0.25}
	p := pCoefficient[2:]
	for i := 0; i < WAVELET_TILE_SIZE; i++ {
		to[i*stride] = 0
		for k := i / 2; k < i/2+1; k++ {
			to[i*stride] += p[i-2*k] * from[absmod(k, WAVELET_TILE_SIZE/2)*stride]
		}
	}
}

func noiseWaveletInit(data *Noise) {
	sz := WAVELET_TILE_SIZE * WAVELET_TILE_SIZE * WAVELET_TILE_SIZE * int(unsafe.Sizeof(float64(0)))
	temp1 := make([]float64, sz)
	temp2 := make([]float64, sz)
	noise := make([]float64, sz)
	for i := 0; i < WAVELET_TILE_SIZE*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE; i++ {
		noise[i] = TCOD_random_get_float(-1.0, 1.0)
	}
	for iy := 0; iy < WAVELET_TILE_SIZE; iy++ {
		for iz := 0; iz < WAVELET_TILE_SIZE; iz++ {
			i := iy*WAVELET_TILE_SIZE + iz*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE
			noiseWaveletDownsample(noise[i:], temp1[i:], 1)
			noiseWaveletUpsample(temp1[i:], temp2[i:], 1)
		}
	}
	for ix := 0; ix < WAVELET_TILE_SIZE; ix++ {
		for iz := 0; iz < WAVELET_TILE_SIZE; iz++ {
			i := ix + iz*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE
			noiseWaveletDownsample(temp2[i:], temp1[i:], WAVELET_TILE_SIZE)
			noiseWaveletUpsample(temp1[i:], temp2[i:], WAVELET_TILE_SIZE)
		}
	}
	for ix := 0; ix < WAVELET_TILE_SIZE; ix++ {
		for iy := 0; iy < WAVELET_TILE_SIZE; iy++ {
			i := ix + iy*WAVELET_TILE_SIZE
			noiseWaveletDownsample(temp2[i:], temp1[i:], WAVELET_TILE_SIZE*WAVELET_TILE_SIZE)
			noiseWaveletUpsample(temp1[i:], temp2[i:], WAVELET_TILE_SIZE*WAVELET_TILE_SIZE)
		}
	}
	for i := 0; i < WAVELET_TILE_SIZE*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE; i++ {
		noise[i] -= temp2[i]
	}
	offset := WAVELET_TILE_SIZE / 2
	if offset&1 == 0 {
		offset++
	}
	for i, ix := 0, 0; ix < WAVELET_TILE_SIZE; ix++ {
		for iy := 0; iy < WAVELET_TILE_SIZE; iy++ {
			for iz := 0; iz < WAVELET_TILE_SIZE; iz++ {
				temp1[i] = noise[absmod(ix+offset, WAVELET_TILE_SIZE)+absmod(iy+offset, WAVELET_TILE_SIZE)*WAVELET_TILE_SIZE+absmod(iz+offset, WAVELET_TILE_SIZE)*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE]
				i++
			}
		}
	}
	for i := 0; i < WAVELET_TILE_SIZE*WAVELET_TILE_SIZE*WAVELET_TILE_SIZE; i++ {
		noise[i] += temp1[i]
	}
	data.waveletTileData = noise
}

func noiseWavelet(data *Noise, f []float64) float64 {
	const n = WAVELET_TILE_SIZE
	if data.ndim <= 0 || data.ndim > 3 {
		return float64(math.NaN()) // not supported
	}
	if data.waveletTileData == nil {
		noiseWaveletInit(data)
	}
	pf := [3]float64{0, 0, 0}
	for i := 0; i < data.ndim; i++ {
		pf[i] = f[i] * WAVELET_SCALE
	}
	mid := [3]int{}
	w := [3][3]float64{}
	for i := 0; i < 3; i++ {
		mid[i] = int(math.Ceil(float64(pf[i] - 0.5)))
		t := float64(mid[i]) - (pf[i] - 0.5)
		w[i][0] = t * t * 0.5
		w[i][2] = (1.0 - t) * (1.0 - t) * 0.5
		w[i][1] = 1.0 - w[i][0] - w[i][2]
	}
	result := float64(0)
	p := [3]int{}
	c := [3]int{}
	for p[2] = -1; p[2] <= 1; p[2]++ {
		for p[1] = -1; p[1] <= 1; p[1]++ {
			for p[0] = -1; p[0] <= 1; p[0]++ {
				weight := float64(1.0)
				for i := 0; i < 3; i++ {
					c[i] = absmod(mid[i]+p[i], n)
					weight *= w[i][p[i]+1]
				}
				result += weight * data.waveletTileData[c[2]*n*n+c[1]*n+c[0]]
			}
		}
	}
	return clampSignedF(float64(result))
}

func noiseFbmWavelet(noise *Noise, f []float64, octaves float64) float64 {
	return noiseFbmInt(noise, f, octaves, noiseWavelet)
}

func noiseTurbulenceWavelet(noise *Noise, f []float64, octaves float64) float64 {
	return noiseTurbulenceInt(noise, f, octaves, noiseWavelet)
}

func noiseSetType(noise *Noise, noiseType NoiseType) {
	noise.noiseType = noiseType
}

func noiseGetEx(noise *Noise, f []float64, noiseType NoiseType) float64 {
	switch noiseType {
	case NoiseTypePerlin:
		return noisePerlin(noise, f)
	case NoiseTypeSimplex:
		return noiseSimplex(noise, f)
	case NoiseTypeWavelet:
		return noiseWavelet(noise, f)
	default:
		return float64(math.NaN())
	}
}

func noiseGetFbmEx(noise *Noise, f []float64, octaves float64, noiseType NoiseType) float64 {
	switch noiseType {
	case NoiseTypePerlin:
		return noiseFbmPerlin(noise, f, octaves)
	case NoiseTypeSimplex:
		return noiseFbmSimplex(noise, f, octaves)
	case NoiseTypeWavelet:
		return noiseFbmWavelet(noise, f, octaves)
	default:
		return float64(math.NaN())
	}
}

func noiseGetTurbulenceEx(noise *Noise, f []float64, octaves float64, noiseType NoiseType) float64 {
	switch noiseType {
	case NoiseTypePerlin:
		return noiseTurbulencePerlin(noise, f, octaves)
	case NoiseTypeSimplex:
		return noiseTurbulenceSimplex(noise, f, octaves)
	case NoiseTypeWavelet:
		return noiseTurbulenceWavelet(noise, f, octaves)
	default:
		return float64(math.NaN())
	}
}

func noiseGet(noise *Noise, f []float64) float64 {
	return noiseGetEx(noise, f, noise.noiseType)
}

func noiseGetFbm(noise *Noise, f []float64, octaves float64) float64 {
	return noiseGetFbmEx(noise, f, octaves, noise.noiseType)
}

func noiseGetTurbulence(noise *Noise, f []float64, octaves float64) float64 {
	return noiseGetTurbulenceEx(noise, f, octaves, noise.noiseType)
}

/*
func noiseDelete(noise *Noise) {
	if noise != nil && noise.waveletTileData != nil {
		C.free(unsafe.Pointer(noise.waveletTileData))
	}
	C.free(unsafe.Pointer(noise))
}
*/

func noiseGetVectorized(
	noise *Noise,
	noiseType NoiseType,
	n int,
	x []float64,
	y []float64,
	z []float64,
	w []float64,
	out []float64,
) {
	for i := 0; i < n; i++ {
		point := [4]float64{
			x[i],
			0,
			0,
			0,
		}
		if y != nil && noise.ndim >= 2 {
			point[1] = y[i]
		}
		if z != nil && noise.ndim >= 3 {
			point[2] = z[i]
		}
		if w != nil && noise.ndim >= 4 {
			point[3] = w[i]
		}
		switch noiseType {
		case NoiseTypePerlin:
			out[i] = noisePerlin(noise, point[:])
		case NoiseTypeSimplex:
			out[i] = noiseSimplex(noise, point[:])
		case NoiseTypeWavelet:
			out[i] = noiseWavelet(noise, point[:])
		default:
			out[i] = float64(math.NaN())
		}
	}
}
func noiseGetFbmVectorized(
	noise *Noise,
	noiseType NoiseType,
	octaves float64,
	n int,
	x []float64,
	y []float64,
	z []float64,
	w []float64,
	out []float64,
) {
	for i := 0; i < n; i++ {
		point := [4]float64{
			x[i],
			0,
			0,
			0,
		}
		if y != nil && noise.ndim >= 2 {
			point[1] = y[i]
		}
		if z != nil && noise.ndim >= 3 {
			point[2] = z[i]
		}
		if w != nil && noise.ndim >= 4 {
			point[3] = w[i]
		}
		switch noiseType {
		case NoiseTypePerlin:
			out[i] = noiseFbmPerlin(noise, point[:], octaves)
		case NoiseTypeSimplex:
			out[i] = noiseFbmSimplex(noise, point[:], octaves)
		case NoiseTypeWavelet:
			out[i] = noiseFbmWavelet(noise, point[:], octaves)
		default:
			out[i] = float64(math.NaN())
		}
	}
}
func noiseGetTurbulenceVectorized(
	noise *Noise,
	noiseType NoiseType,
	octaves float64,
	n int,
	x []float64,
	y []float64,
	z []float64,
	w []float64,
	out []float64,
) {
	for i := 0; i < n; i++ {
		point := [4]float64{
			x[i],
			0,
			0,
			0,
		}
		if y != nil && noise.ndim >= 2 {
			point[1] = y[i]
		}
		if z != nil && noise.ndim >= 3 {
			point[2] = z[i]
		}
		if w != nil && noise.ndim >= 4 {
			point[3] = w[i]
		}
		switch noiseType {
		case NoiseTypePerlin:
			out[i] = noiseTurbulencePerlin(noise, point[:], octaves)
		case NoiseTypeSimplex, NoiseTypeDefault:
			out[i] = noiseTurbulenceSimplex(noise, point[:], octaves)
		case NoiseTypeWavelet:
			out[i] = noiseTurbulenceWavelet(noise, point[:], octaves)
		default:
			out[i] = float64(math.NaN())
		}
	}
}
