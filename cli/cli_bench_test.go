package cli

import "testing"

//go:generate go-bind-plugin -plugin-package ../internal/test_fixtures/benchmark_plugin -output-package cli -output-path cli_test_plugin.go -output-name BenchmarkPlugin
func BenchmarkCallOverhead(b *testing.B) {
	pl, err := BindBenchmarkPlugin("plugin.so")
	numbers := 100
	sl := prepareLargeSlice(numbers)
	expectedSum := uint64(numbers * (numbers + 1) / 2)

	if err != nil {
		b.Fatalf("Could not setup benchmark: %s", err)
	}

	b.Log("Running benchmarks...")

	b.Run("plugin", func(b *testing.B) {

		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			sum := pl.Sum(sl)

			if sum != expectedSum {
				b.Fatalf("Expected sum: %d, got: %d", expectedSum, sum)
			}
		}

	})

	b.Run("native", func(b *testing.B) {

		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			sum := nativeSum(sl)

			if sum != expectedSum {
				b.Fatalf("Expected sum: %d, got: %d", expectedSum, sum)
			}
		}
	})
}

func prepareLargeSlice(size int) []uint64 {
	var ret = make([]uint64, size+1)

	for i := 0; i <= size; i++ {
		ret[i] = uint64(i)
	}

	return ret
}

func nativeSum(vals []uint64) uint64 {
	var sum uint64

	for _, val := range vals {
		sum += val
	}

	return sum
}
