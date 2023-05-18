# Bloom Filter


A Bloom Filter is a probabilistic data structure used to check whether an element is in a set. It allows for a certain degree of false positives, meaning that it may mistakenly think that an element is in the set at times, but it will never wrongly think that an element is not in the set. The false-positive rate can be adjusted according to actual needs.

This formula is derived from the following parameters:

* M: The size of the Bloom Filter (number of bits)
* n: The number of elements to be inserted into the Bloom Filter
* k: The number of hash functions used
* p: The probability of an element being misjudged as being in the set
* f: The desired false-positive rate

This algorithm is based on the principles and mathematical model of the Bloom Filter. In the optimal case, when the number of hash functions k in the Bloom Filter is at its optimal value, the false-positive rate is minimized. This optimal value k can be obtained by the following formula:

```
k = (M / n) * ln(2)
```
The formula for the false-positive rate p is as follows:

```
p = (1 - e^(-kn/M))^k
```

Here e is the base of the natural logarithm. This formula can be derived through probability theory and mathematical analysis. We can use mathematical methods to solve for the optimal M value so that the false-positive rate reaches the desired value f. According to the formula derivation, we can get:

```
M = -n * ln(f) / (ln(2))^2
```

This is the algorithm we use in this function. It calculates the optimal size M of the Bloom Filter based on the given desired false-positive rate f and the number of elements n to be inserted. By adjusting the values of M and k, we can balance the false-positive rate and memory usage according to actual needs.


## Usage

```go
func TestSingleElement(t *testing.T) {
	bf := NewBloomFilter[Key](10000, 0.10)
	item := Key("Hello")
	bf.Insert(item)
	if !bf.LookUp(item) {
		t.Errorf("Expected item to exist in the bloom filter")
	}
}
```