# HyperLogLog

HyperLogLog is a probabilistic data structure used in the count-distinct problem, approximating the number of distinct elements in a multiset.

## Core implementation

```go
func (hll *HyperLogLog) Cardinality() float64 {
	sum := 0.0
	zeros := 0
	for _, it := range hll.registers {
		sum += 1.0 / float64(int(1)<<it)
		if it == 0 {
			zeros++
		}
	}
	harmonicAvg := float64(hll.m) / sum
	E := hll.alphaM * float64(hll.m) * harmonicAvg // Raw estimate

	// Computes corrected estimate
	if E <= 2.5*float64(hll.m) { // Small range correction
		if zeros != 0 {
			return float64(hll.m) * math.Log(float64(hll.m)/float64(zeros))
		} else {
			return E
		}
	} else if E <= math.Pow(2, 32)/float64(30) { // Intermediate range, no correction
		return E
	} else { // E > 2^32/30, Large range correction
		return -math.Pow(2, 32) * math.Log(1-E/math.Pow(2, 32))
	}
}
```