package dfstyle

import (
	"math"
	"math/rand"
)

type Heightmap struct {
	w, h   int
	values []float64
}

func MIN(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MAX(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns true if `x`,`y` are valid coordinates for this heightmap.
func in_bounds(hm *Heightmap, x, y int) bool {
	if hm == nil {
		return false
	} // No valid coordinates on a NULL pointer.
	if x < 0 || x >= hm.w {
		return false
	}
	if y < 0 || y >= hm.h {
		return false
	}
	return true
}

// Returns true if these heightmaps have the same shape and are non-NULL.
func is_same_size(hm1, hm2 *Heightmap) bool {
	if hm1 == nil || hm2 == nil {
		return false
	}
	if hm1.w != hm2.w || hm1.h != hm2.h {
		return false
	}
	return true
}

func HeightmapNew(w, h int) *Heightmap {
	hm := &Heightmap{w: w, h: h}
	hm.values = make([]float64, w*h)
	return hm
}

func (hm *Heightmap) Clear() {
	if hm == nil {
		return
	}
	for i := 0; i != hm.w*hm.h; i++ {
		hm.values[i] = 0
	}
}

func (hm *Heightmap) GetValue(x, y int) float64 {
	if in_bounds(hm, x, y) {
		return hm.values[x+y*hm.w]
	}
	return 0.0
}

func (hm *Heightmap) SetValue(x, y int, value float64) {
	if in_bounds(hm, x, y) {
		hm.values[x+y*hm.w] = value
	}
}

func (hm *Heightmap) GetMinMax(min, max *float64) {
	if !in_bounds(hm, 0, 0) {
		*min = 0
		*max = 0
		return
	}
	if min != nil {
		*min = hm.values[0]
	}
	if max != nil {
		*max = hm.values[0]
	}
	for i := 0; i != hm.h*hm.w; i++ {
		value := hm.values[i]
		if min != nil {
			*min = math.Min(*min, value)
		}
		if max != nil {
			*max = math.Max(*max, value)
		}
	}
}

func (hm *Heightmap) Normalize(min, max float64) {
	if hm == nil {
		return
	}
	var current_min, current_max float64
	hm.GetMinMax(&current_min, &current_max)
	if current_max-current_min < math.SmallestNonzeroFloat64 {
		for i := 0; i != hm.w*hm.h; i++ {
			hm.values[i] = min
		}
	} else {
		normalize_scale := (max - min) / (current_max - current_min)
		for i := 0; i != hm.w*hm.h; i++ {
			hm.values[i] = min + (hm.values[i]-current_min)*normalize_scale
		}
	}
}

func (hm *Heightmap) AddHill(hx, hy int, h_radius, h_height float64) {
	if hm == nil {
		return
	}
	h_radius2 := h_radius * h_radius
	coef := h_height / h_radius2
	minx := int(math.Max(float64(hx)-h_radius, 0))
	miny := int(math.Max(float64(hy)-h_radius, 0))
	maxx := int(math.Min(math.Ceil(float64(hx)+h_radius), float64(hm.w)))
	maxy := int(math.Min(math.Ceil(float64(hy)+h_radius), float64(hm.h)))
	for y := miny; y < maxy; y++ {
		y_dist := (float64(y) - float64(hy)) * (float64(y) - float64(hy))
		for x := minx; x < maxx; x++ {
			x_dist := (float64(x) - float64(hx)) * (float64(x) - float64(hx))
			z := h_radius2 - x_dist - y_dist
			if z > 0 {
				hm.values[y*hm.w+x] += z * coef
			}
		}
	}
}

func (hm *Heightmap) DigHill(hx, hy, h_radius, h_height float64) {
	if hm == nil {
		return
	}
	h_radius2 := h_radius * h_radius
	coef := h_height / h_radius2
	minx := int(math.Max(float64(hx-h_radius), 0))
	miny := int(math.Max(float64(hy-h_radius), 0))
	maxx := int(math.Min(math.Ceil(float64(hx+h_radius)), float64(hm.w)))
	maxy := int(math.Min(math.Ceil(float64(hy+h_radius)), float64(hm.h)))
	for y := miny; y < maxy; y++ {
		for x := minx; x < maxx; x++ {
			x_dist := (float64(x) - hx) * (float64(x) - hx)
			y_dist := (float64(y) - hy) * (float64(y) - hy)
			dist := x_dist + y_dist
			if dist < h_radius2 {
				z := (h_radius2 - dist) * coef
				if h_height > 0 {
					if hm.GetValue(x, y) < z {
						hm.SetValue(x, y, z)
					}
				} else {
					if hm.GetValue(x, y) > z {
						hm.SetValue(x, y, z)
					}
				}
			}
		}
	}
}

func CopyHeightmap(hm_source, hm_dest *Heightmap) {
	if !is_same_size(hm_source, hm_dest) {
		return
	}
	copy(hm_dest.values, hm_source.values)
}

func (hm *Heightmap) AddFBM(noise *Noise, mul_x, mul_y, add_x, add_y, octaves, delta, scale float64) {
	if hm == nil {
		return
	}
	x_coefficient := mul_x / float64(hm.w)
	y_coefficient := mul_y / float64(hm.h)
	for y := 0; y < hm.h; y++ {
		for x := 0; x < hm.w; x++ {
			f := [2]float64{(float64(x) + add_x) * x_coefficient, (float64(y) + add_y) * y_coefficient}
			hm.values[y*hm.w+x] += delta + noiseGetFbm(noise, f[:], octaves)*scale
		}
	}
}

func (hm *Heightmap) ScaleFBM(noise *Noise, mul_x, mul_y, add_x, add_y, octaves, delta, scale float64) {
	if hm == nil {
		return
	}
	x_coefficient := mul_x / float64(hm.w)
	y_coefficient := mul_y / float64(hm.h)
	for y := 0; y < hm.h; y++ {
		for x := 0; x < hm.w; x++ {
			f := [2]float64{(float64(x) + add_x) * x_coefficient, (float64(y) + add_y) * y_coefficient}
			hm.values[y*hm.w+x] *= delta + noiseGetFbm(noise, f[:], octaves)*scale
		}
	}
}

func (hm *Heightmap) GetInterpolatedValue(x, y float64) float64 {
	if hm == nil {
		return 0.0
	}
	x = clamp(0.0, float64(hm.w-1), x)
	y = clamp(0.0, float64(hm.h-1), y)
	var fix, fiy float64
	fx := math.Mod(float64(x), fix)
	fy := math.Mod(float64(y), fiy)
	ix := int(fix)
	iy := int(fiy)

	if ix >= hm.w-1 {
		ix = hm.w - 2
		fx = 1.0
	}
	if iy >= hm.h-1 {
		iy = hm.h - 2
		fy = 1.0
	}
	c1 := hm.GetValue(ix, iy)
	c2 := hm.GetValue(ix+1, iy)
	c3 := hm.GetValue(ix, iy+1)
	c4 := hm.GetValue(ix+1, iy+1)
	top := lerp(c1, c2, fx)
	bottom := lerp(c3, c4, fx)
	return lerp(top, bottom, fy)
}

func (hm *Heightmap) GetNormal(x, y, waterLevel float64, n *[3]float64) {
	if hm == nil {
		return
	}
	var h0, hx, hy, invlen float64
	n[0] = 0.0
	n[1] = 0.0
	n[2] = 1.0
	if x >= float64(hm.w-1) || y >= float64(hm.h-1) {
		return
	}
	h0 = hm.GetInterpolatedValue(x, y)
	if h0 < waterLevel {
		h0 = waterLevel
	}
	hx = hm.GetInterpolatedValue(x+1, y)
	if hx < waterLevel {
		hx = waterLevel
	}
	hy = hm.GetInterpolatedValue(x, y+1)
	if hy < waterLevel {
		hy = waterLevel
	}
	/* vx = 1       vy = 0 */
	/*      0            1 */
	/*      hx-h0        hy-h0 */
	/* vz = vx cross vy */
	n[0] = 255 * (h0 - hx)
	n[1] = 255 * (h0 - hy)
	n[2] = 16.0
	/* normalize */
	invlen = 1.0 / float64(math.Sqrt(float64(n[0]*n[0]+n[1]*n[1]+n[2]*n[2])))
	n[0] *= invlen
	n[1] *= invlen
	n[2] *= invlen
}

func (hm *Heightmap) DigBezier(px [4]int, py [4]int, startRadius, startDepth, endRadius, endDepth float64) {
	if hm == nil {
		return
	}
	xFrom := px[0]
	yFrom := py[0]
	for i := 0; i <= 1000; i++ {
		t := float64(i) / 1000.0
		it := 1.0 - t
		xTo := int(float64(px[0])*it*it*it + 3*float64(px[1])*t*it*it + 3*float64(px[2])*t*t*it + float64(px[3])*t*t*t)
		yTo := int(float64(py[0])*it*it*it + 3*float64(py[1])*t*it*it + 3*float64(py[2])*t*t*it + float64(py[3])*t*t*t)
		if xTo != xFrom || yTo != yFrom {
			radius := startRadius + (endRadius-startRadius)*t
			depth := startDepth + (endDepth-startDepth)*t
			hm.DigHill(float64(xTo), float64(yTo), radius, depth)
			xFrom = xTo
			yFrom = yTo
		}
	}
}

func (hm *Heightmap) HasLandOnBorder(waterLevel float64) bool {
	if hm == nil {
		return false
	}
	for x := 0; x < hm.w; x++ {
		if hm.GetValue(x, 0) > waterLevel || hm.GetValue(x, hm.h-1) > waterLevel {
			return true
		}
	}
	for y := 0; y < hm.h; y++ {
		if hm.GetValue(0, y) > waterLevel || hm.GetValue(hm.w-1, y) > waterLevel {
			return true
		}
	}
	return false
}

func (hm *Heightmap) Add(value float64) {
	for i := range hm.values {
		hm.values[i] += value
	}
}

func (hm *Heightmap) AddValue(value float64) {
	if hm == nil {
		return
	}
	for i := 0; i < hm.w*hm.h; i++ {
		hm.values[i] += value
	}
}

func (hm *Heightmap) CountCells(min, max float64) int {
	if hm == nil {
		return 0
	}
	count := 0
	for i := 0; i < hm.w*hm.h; i++ {
		if hm.values[i] >= min && hm.values[i] <= max {
			count++
		}
	}
	return count
}

func (hm *Heightmap) Scale(value float64) {
	if hm == nil {
		return
	}
	for i := 0; i < hm.w*hm.h; i++ {
		hm.values[i] *= value
	}
}

func (hm *Heightmap) Clamp(min, max float64) {
	if hm == nil {
		return
	}
	for i := 0; i < hm.w*hm.h; i++ {
		hm.values[i] = clamp(min, max, hm.values[i])
	}
}

func HeightmapLerpHm(hm1 *Heightmap, hm2 *Heightmap, hm_out *Heightmap, coef float64) {
	if !is_same_size(hm1, hm2) || !is_same_size(hm1, hm_out) {
		return
	}
	for i := 0; i < hm1.w*hm1.h; i++ {
		hm_out.values[i] = lerp(hm1.values[i], hm2.values[i], coef)
	}
}

func HeightmapAddHm(hm1 *Heightmap, hm2 *Heightmap, hm_out *Heightmap) {
	if !is_same_size(hm1, hm2) || !is_same_size(hm1, hm_out) {
		return
	}
	for i := 0; i < hm1.w*hm1.h; i++ {
		hm_out.values[i] = hm1.values[i] + hm2.values[i]
	}
}

func HeightmapMultiplyHm(hm1 *Heightmap, hm2 *Heightmap, hm_out *Heightmap) {
	if !is_same_size(hm1, hm2) || !is_same_size(hm1, hm_out) {
		return
	}
	for i := 0; i < hm1.w*hm1.h; i++ {
		hm_out.values[i] = hm1.values[i] * hm2.values[i]
	}
}

func (hm *Heightmap) GetSlope(x, y int) float64 {
	dix := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}
	diy := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
	min_dy := float64(0.0)
	max_dy := float64(0.0)
	if !in_bounds(hm, x, y) {
		return 0
	}
	v := hm.GetValue(x, y)
	for i := 0; i < 8; i++ {
		nx := x + dix[i]
		ny := y + diy[i]
		if in_bounds(hm, nx, ny) {
			n_slope := hm.GetValue(nx, ny) - v
			min_dy = math.Min(min_dy, n_slope)
			max_dy = math.Max(max_dy, n_slope)
		}
	}
	return float64(math.Atan2(float64(max_dy+min_dy), 1.0))
}

func (hm *Heightmap) RainErosion(nbDrops int, erosionCoef, aggregationCoef float64) {
	if hm == nil {
		return
	}
	for nbDrops > 0 {
		curx := randInt(0, hm.w-1)
		cury := randInt(0, hm.h-1)
		dix := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}
		diy := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
		sediment := float64(0.0)
		for {
			next_x := 0
			next_y := 0
			v := hm.GetValue(curx, cury)
			slope := float64(math.Inf(-1))
			for i := 0; i < 8; i++ {
				nx := curx + dix[i]
				ny := cury + diy[i]
				if !in_bounds(hm, nx, ny) {
					continue
				}
				n_slope := v - hm.GetValue(nx, ny)
				if n_slope > slope {
					slope = n_slope
					next_x = nx
					next_y = ny
				}
			}
			if slope > 0.0 {
				hm.SetValue(curx, cury, hm.GetValue(curx, cury)-erosionCoef*slope)
				curx = next_x
				cury = next_y
				sediment += slope
			} else {
				hm.SetValue(curx, cury, hm.GetValue(curx, cury)+aggregationCoef*sediment)
				break
			}
		}
		nbDrops--
	}
}

func (hm *Heightmap) HeatErosion(nbPass int, minSlope, erosionCoef, aggregationCoef float64) {
	if hm == nil {
		return
	}
	for nbPass > 0 {
		for y := 0; y < hm.h; y++ {
			for x := 0; x < hm.w; x++ {
				dix := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}
				diy := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
				next_x := 0
				next_y := 0
				v := hm.GetValue(x, y)
				slope := float64(0.0)
				for i := 0; i < 8; i++ {
					nx := x + dix[i]
					ny := y + diy[i]
					if in_bounds(hm, nx, ny) {
						n_slope := v - hm.GetValue(nx, ny)
						if n_slope > slope {
							slope = n_slope
							next_x = nx
							next_y = ny
						}
					}
				}
				if slope > minSlope {
					hm.SetValue(x, y, hm.GetValue(x, y)-erosionCoef*(slope-minSlope))
					hm.SetValue(next_x, next_y, hm.GetValue(next_x, next_y)+aggregationCoef*(slope-minSlope))
				}
			}
		}
		nbPass--
	}
}

func (hm *Heightmap) KernelTransform(kernel_size int, dx []int, dy []int, weight []float64, minLevel float64, maxLevel float64) {
	if hm == nil {
		return
	}
	for y := 0; y < hm.h; y++ {
		for x := 0; x < hm.w; x++ {
			if hm.GetValue(x, y) >= minLevel && hm.GetValue(x, y) <= maxLevel {
				val := float64(0.0)
				totalWeight := float64(0.0)
				for i := 0; i < kernel_size; i++ {
					nx := x + dx[i]
					ny := y + dy[i]
					if in_bounds(hm, nx, ny) {
						val += weight[i] * hm.GetValue(nx, ny)
						totalWeight += weight[i]
					}
				}
				hm.SetValue(x, y, val/totalWeight)
			}
		}
	}
}

func (hm *Heightmap) AddVoronoi(nbPoints, nbCoef int, coef []float64) {
	if hm == nil {
		return
	}
	type point_t struct {
		x, y int
		dist float64
	}
	if nbPoints <= 0 {
		return
	}
	pt := make([]point_t, nbPoints)
	nbCoef = MIN(nbCoef, nbPoints)
	for i := 0; i < nbPoints; i++ {
		pt[i].x = randInt(0, hm.w-1)
		pt[i].y = randInt(0, hm.h-1)
	}
	for y := 0; y < hm.h; y++ {
		for x := 0; x < hm.w; x++ {
			/* calculate distance to voronoi points */
			for i := 0; i < nbPoints; i++ {
				dx := pt[i].x - x
				dy := pt[i].y - y
				pt[i].dist = float64(dx*dx + dy*dy)
			}
			for i := 0; i < nbCoef; i++ {
				/* get closest point */
				minDist := float64(1e8)
				idx := -1
				for j := 0; j < nbPoints; j++ {
					if pt[j].dist < minDist {
						idx = j
						minDist = pt[j].dist
					}
				}
				if idx == -1 {
					break
				}
				hm.SetValue(x, y, hm.GetValue(x, y)+coef[i]*pt[idx].dist)
				pt[idx].dist = 1e8
			}
		}
	}
}

func (hm *Heightmap) MidPointDisplacement(roughness float64) {
	if hm == nil {
		return
	}
	step := 1
	offset := float64(1.0)
	initsz := MIN(hm.w, hm.h) - 1
	sz := initsz
	hm.SetValue(0, 0, randFloat(0.0, 1.0))
	hm.SetValue(sz-1, 0, randFloat(0.0, 1.0))
	hm.SetValue(0, sz-1, randFloat(0.0, 1.0))
	hm.SetValue(sz-1, sz-1, randFloat(0.0, 1.0))
	for sz > 0 {
		// diamond step
		for y := 0; y < step; y++ {
			for x := 0; x < step; x++ {
				diamond_x := sz/2 + x*sz
				diamond_y := sz/2 + y*sz
				z := hm.GetValue(x*sz, y*sz)
				z += hm.GetValue((x+1)*sz, y*sz)
				z += hm.GetValue((x+1)*sz, (y+1)*sz)
				z += hm.GetValue(x*sz, (y+1)*sz)
				z *= 0.25
				hm.setMPDHeight(diamond_x, diamond_y, z, offset)
			}
		}
		offset *= roughness
		// square step
		for y := 0; y < step; y++ {
			for x := 0; x < step; x++ {
				diamond_x := sz/2 + x*sz
				diamond_y := sz/2 + y*sz
				// north
				setMDPHeightSquare(hm, diamond_x, diamond_y-sz/2, initsz, sz/2, offset)
				// south
				setMDPHeightSquare(hm, diamond_x, diamond_y+sz/2, initsz, sz/2, offset)
				// west
				setMDPHeightSquare(hm, diamond_x-sz/2, diamond_y, initsz, sz/2, offset)
				// east
				setMDPHeightSquare(hm, diamond_x+sz/2, diamond_y, initsz, sz/2, offset)
			}
		}
		sz /= 2
		step *= 2
	}
}

/* private stuff */
func (hm *Heightmap) setMPDHeight(x, y int, z, offset float64) {
	z += randFloat(-offset, offset)
	hm.SetValue(x, y, z)
}

func setMDPHeightSquare(hm *Heightmap, x, y, initsz, sz int, offset float64) {
	z := float64(0.0)
	count := 0
	if y >= sz {
		z += hm.GetValue(x, y-sz)
		count++
	}
	if x >= sz {
		z += hm.GetValue(x-sz, y)
		count++
	}
	if y+sz < initsz {
		z += hm.GetValue(x, y+sz)
		count++
	}
	if x+sz < initsz {
		z += hm.GetValue(x+sz, y)
		count++
	}
	z /= float64(count)
	hm.setMPDHeight(x, y, z, offset)
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func clamp(min, max, v float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lerp(a, b, x float64) float64 {
	return a + (b-a)*x
}

func floor(a float64) int {
	if a > 0 {
		return int(a)
	}
	return int(a) - 1
}

func cubic(a float64) float64 {
	return a * a * (3 - 2*a)
}

func genericSwap(x, y interface{}) {
	x, y = y, x
}

// Return a floating point value clamped between -1.0f and 1.0f exclusively.
// The return value excludes -1.0f and 1.0f to avoid rounding issues.

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
