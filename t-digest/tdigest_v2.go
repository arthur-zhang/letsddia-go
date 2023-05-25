package t_digest

import (
	"lo"
	"math"
	"sort"
)

const delta = 10.0

type TDigestV2 struct {
	centroids    []Centroid
	maxCentroids int
	count        int
}

func mergeCentroid(a, b Centroid) Centroid {

	weight := a.weight + b.weight
	mean := (a.mean*float64(a.weight) + b.mean*float64(b.weight)) / float64(weight)
	return Centroid{
		mean:   mean,
		weight: weight,
	}
}
func NewTDigestV2(maxCentroids int) TDigestV2 {
	centroids := make([]Centroid, 0)
	return TDigestV2{
		centroids:    centroids,
		maxCentroids: maxCentroids,
		count:        0,
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
	ys := make([]Centroid, 0)
	ys = append(ys, t.centroids[0])
	minPotential := getPotential(0.0)
	total := t.centroids[0].weight
	for i := 1; i < len(t.centroids); i++ {
		x := t.centroids[i]
		nextQid := float64(total+x.weight) / float64(totalWeight)
		if getPotential(nextQid)-minPotential <= 1 {
			ys[len(ys)-1] = mergeCentroid(ys[len(ys)-1], x)
		} else {
			ys = append(ys, x)
			minPotential = getPotential(float64(total) / float64(totalWeight))
		}
		total += x.weight
	}

	t.centroids = ys
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

	totalCount := lo.SumBy(t.centroids, func(c Centroid) int {
		return c.weight
	})
	idx := q * float64(totalCount)

	maxIdx := float64(t.centroids[0].weight / 2)

	if idx < maxIdx {
		return t.centroids[0].mean
	}

	for i := 0; i < len(t.centroids)-1; i++ {
		c := t.centroids[i]
		nextC := t.centroids[i+1]
		intervalLength := float64(c.weight+nextC.weight) / 2
		if idx < maxIdx+intervalLength {
			return c.mean*(1-(idx-maxIdx)/intervalLength) + nextC.mean*((idx-maxIdx)/intervalLength)
		}
		maxIdx += intervalLength
	}
	return t.centroids[len(t.centroids)-1].mean
}

func getPotential(q float64) float64 {
	return delta * (math.Asin(2*q-1) / (2 * math.Pi))
}
