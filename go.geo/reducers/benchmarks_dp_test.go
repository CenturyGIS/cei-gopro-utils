package reducers

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"

	"cei-gopro-utils/go.geo"
)

func TestDouglasPeuckerBenchmarkData(t *testing.T) {
	type reduceTest struct {
		Threshold float64
		Length    int
	}

	tests := []reduceTest{
		{0.1, 1118},
		{0.5, 257},
		{1.0, 144},
		{1.5, 95},
		{2.0, 71},
		{3.0, 46},
		{4.0, 39},
		{5.0, 33},
	}
	path := benchmarkData()
	for i := range tests {
		p := DouglasPeucker(path, tests[i].Threshold)
		if p.Length() != tests[i].Length {
			t.Errorf("douglas peucker benchmark data reduced poorly, got %d, expected %d", p.Length(), tests[i].Length)
		}
	}
}

func BenchmarkDouglasPeucker(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DouglasPeucker(path, 0.1)
	}
}

func BenchmarkDouglasPeuckerIndexMap(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DouglasPeuckerIndexMap(path, 0.1)
	}
}

func benchmarkData() *geo.Path {
	// Data taken from the simplify-js example at http://mourner.github.io/simplify-js/
	f, err := os.Open("lisbon2portugal.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// decompress and decode the json
	var points []float64
	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	json.NewDecoder(gzipReader).Decode(&points)

	// create the geo path
	path := geo.NewPathPreallocate(0, len(points)/2)
	for i := 0; i < len(points); i += 2 {
		path.Push(geo.NewPoint(points[i], points[i+1]))
	}

	return path
}
