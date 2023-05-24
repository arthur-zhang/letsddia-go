package t_digest

import (
	"lo"
	"math"
	"sort"
)

const delta = 0.01

type TDigestV2 struct {
	centroids    []Centroid
	maxCentroids int
}

func NewTDigestV2(maxCentroids int) TDigestV2 {
	centroids := make([]Centroid, 0)
	return TDigestV2{
		centroids:    centroids,
		maxCentroids: maxCentroids,
	}
}
func (t *TDigestV2) Insert(x float64) {
	centroid := Centroid{
		mean:  x,
		count: 1,
	}
	t.centroids = append(t.centroids, centroid)
	if len(t.centroids) > t.maxCentroids {
		t.compress()
	}
}

func (t *TDigestV2) compress() {
	// sort first
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})
	newCentroids := make([]Centroid, 0)
	totalWeight := float64(0)

	for _, c := range t.centroids {
		var merged bool
		for i := range newCentroids {
			newWeight := newCentroids[i].count + c.count
			newMean := (newCentroids[i].mean*float64(newCentroids[i].count) + c.mean*float64(c.count)) / float64(newWeight)
			if float64(newWeight) <= (totalWeight+float64(newWeight))*delta/float64(2*t.maxCentroids) {
				newCentroids[i].count = newWeight
				newCentroids[i].mean = newMean
				totalWeight += float64(c.count)
				merged = true
				break
			}
		}

		if !merged {
			newCentroids = append(newCentroids, Centroid{mean: c.mean, count: c.count})
			totalWeight += float64(c.count)
		}
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

	totalCount := lo.SumBy(t.centroids, func(c Centroid) int {
		return c.count
	})
	target := q * float64(totalCount)

	totalCount = 0
	for _, c := range t.centroids {
		if float64(totalCount+c.count) >= target {
			return c.mean
		}
		totalCount += c.count
	}
	return t.centroids[len(t.centroids)-1].mean
}
