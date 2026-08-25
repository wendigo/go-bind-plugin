package main

// Sum add all values on the slice
func Sum(vals []uint64) uint64 {
	var sum uint64

	for _, val := range vals {
		sum += val
	}

	return sum
}

// main is never called; it exists so the package also compiles with the
// default buildmode (go build ./...), not only with -buildmode=plugin.
func main() {}
