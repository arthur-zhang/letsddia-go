package t_digest

import (
	"lo"
	"math"
	"sort"
)

// δ
const Delta = 10.0

type TDigestV2 struct {
	centroids []Centroid
	count     int
}

func mergeCentroid(a, b Centroid) Centroid {
	weight := a.weight + b.weight
	mean := (a.mean*float64(a.weight) + b.mean*float64(b.weight)) / float64(weight)
	return Centroid{
		mean:   mean,
		weight: weight,
	}
}
func NewTDigestV2() TDigestV2 {
	centroids := make([]Centroid, 0)
	return TDigestV2{
		centroids: centroids,
		count:     0,
	}
}
func (t *TDigestV2) Insert(x float64) {
	centroid := Centroid{
		mean:   x,
		weight: 1,
	}
	t.centroids = append(t.centroids, centroid)
	t.compress()
	t.count += 1
}

func (t *TDigestV2) compress() {
	if len(t.centroids) == 0 {
		return
	}
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})
	totalWeight := lo.SumBy(t.centroids, func(c Centroid) int {
		return c.weight
	})
	newCentroids := make([]Centroid, 0)
	newCentroids = append(newCentroids, t.centroids[0])
	minK := getK(0.0) // -2.5
	weightSoFar := t.centroids[0].weight
	for i := 1; i < len(t.centroids); i++ {
		x := t.centroids[i]
		nextQ := float64(weightSoFar+x.weight) / float64(totalWeight)
		//Two adjacent centroids can be merged if their combined potential difference does not violate constraint
		// σ(b) = π(h) − π(ℓ), σ(b) ≤ 1
		if getK(nextQ)-minK <= 1 {
			newCentroids[len(newCentroids)-1] = mergeCentroid(newCentroids[len(newCentroids)-1], x)
		} else {
			newCentroids = append(newCentroids, x)
			minK = getK(float64(weightSoFar) / float64(totalWeight))
		}
		weightSoFar += x.weight
	}

	t.centroids = newCentroids
}
func (t *TDigestV2) Quantile(q float64) float64 {
	if len(t.centroids) == 0 {
		return math.NaN()
	}
	if len(t.centroids) == 1 {
		return t.centroids[0].mean
	}
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})

	totalWeight := lo.SumBy(t.centroids, func(c Centroid) int { return c.weight })
	idx := q * float64(totalWeight)

	maxIdx := float64(t.centroids[0].weight / 2)

	// idx is on the first half of the first bin
	if idx < maxIdx {
		return t.centroids[0].mean
	}

	for i := 0; i < len(t.centroids)-1; i++ {
		c := t.centroids[i]
		nextC := t.centroids[i+1]
		intervalWeight := float64(c.weight+nextC.weight) / 2
		if idx < maxIdx+intervalWeight {
			lambda := (idx - maxIdx) / intervalWeight
			// q(θ) = (1 − λ)b(i) + λb(i + 1)
			return (1-lambda)*c.mean + lambda*nextC.mean
		}
		maxIdx += intervalWeight
	}
	// idx is on second half the last bin
	return t.centroids[len(t.centroids)-1].mean
}

// k(q) = δ/2π * arcsin(2q − 1)
// q=1, k(q)=δ/2π * arcsin(1) = δ/2π * π/2 = δ/4
// q=0, k(q)=δ/2π * arcsin(-1) = δ/2π * -π/2 = -δ/4
// when δ=10, k(q) ∈ [-2.5, 2.5]
func getK(q float64) float64 {
	return Delta * (math.Asin(2*q-1) / (2 * math.Pi))
}
