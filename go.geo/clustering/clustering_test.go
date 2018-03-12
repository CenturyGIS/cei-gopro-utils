package clustering

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"

	"cei-gopro-utils/go.geo"
)

func TestClusteringClusterClusters(t *testing.T) {
	preclusters, pointers := loadPrefilteredTestClusters(t)
	bound := geo.NewBoundFromPoints(pointers[0].Point(), pointers[1].Point())
	for _, p := range pointers {
		bound.Extend(p.Point())
	}
	bound.GeoPad(1) // for round off

	clusters := ClusterClusters(preclusters, CentroidGeoDistance{}, 30)

	if l := len(clusters); l != 27 {
		t.Errorf("incorrect number of clusters, got %d", l)
	}

	total := 0
	for _, c := range clusters {
		total += len(c.Pointers)
	}
	if total != len(pointers) {
		t.Errorf("missing pointers, got %d", total)
	}

	for i, c := range clusters {
		if c == nil {
			t.Errorf("cluster %d nil", i)
		}

		if c.Centroid == nil {
			t.Errorf("cluster %d center nil", i)
		}

		if !bound.Contains(c.Centroid) {
			t.Errorf("centroid must at least be within original bound, got %v", c.Centroid)
		}

		if len(c.Pointers) == 0 {
			t.Errorf("no pointers in cluster %d", i)
		}

		for _, pointer := range c.Pointers {
			if !bound.Contains(pointer.Point()) {
				t.Errorf("pointer must at least be within original bound, got %v", pointer.Point())
			}
		}
	}
}

func TestClusterGeoClusters(t *testing.T) {
	preclusters, pointers := loadPrefilteredTestClusters(t)
	bound := geo.NewBoundFromPoints(pointers[0].Point(), pointers[1].Point())
	for _, p := range pointers {
		bound.Extend(p.Point())
	}
	bound.GeoPad(1) // for projection loop round off

	clusters := ClusterGeoClusters(preclusters, 30)

	if l := len(clusters); l != 27 {
		t.Errorf("incorrect number of clusters, got %d", l)
	}

	total := 0
	for _, c := range clusters {
		total += len(c.Pointers)
	}
	if total != len(pointers) {
		t.Errorf("missing pointers, got %d", total)
	}

	for i, c := range clusters {
		if c == nil {
			t.Errorf("cluster %d nil", i)
		}

		if c.Centroid == nil {
			t.Errorf("clusters %d center nil", i)
		}

		if !bound.Contains(c.Centroid) {
			t.Errorf("centroid must at least be within original bound, got %v %v", c.Centroid, bound)
		}

		if len(c.Pointers) == 0 {
			t.Errorf("no pointers in cluster %d", i)
		}

		for _, pointer := range c.Pointers {
			if !bound.Contains(pointer.Point()) {
				t.Errorf("pointer must at least be within original bound, got %v", pointer.Point())
			}
		}
	}

	// shouldn't harm original data
	for _, c := range preclusters {
		if !bound.Contains(c.Centroid) {
			t.Errorf("centroid must at least be within original bound, got %v", c.Centroid)
		}

		for _, pointer := range c.Pointers {
			if !bound.Contains(pointer.Point()) {
				t.Errorf("pointer must at least be within original bound, got %v", pointer.Point())
			}
		}
	}
}

// > go test -c && ./clustering.test -test.bench=ClusterClusters -test.cpuprofile=cpu.out -test.benchtime=10s
// > go tool pprof clustering.test cpu.out
func BenchmarkClusterClusters(b *testing.B) {
	clusters, _ := loadPrefilteredTestClusters(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs := ClusterClusters(clusters, CentroidGeoDistance{}, 30)
		if len(cs) != 27 {
			b.Fatalf("incorrect number of clusters, got %v", len(cs))
		}
	}
}

// > go test -c && ./clustering.test -test.bench=ClusterPointers -test.cpuprofile=cpu.out -test.benchtime=10s
// > go tool pprof clustering.test cpu.out
func BenchmarkClusterPointers(b *testing.B) {
	_, pointers := loadPrefilteredTestClusters(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs := ClusterPointers(pointers, CentroidGeoDistance{}, 30)
		if len(cs) != 26 {
			b.Fatalf("incorrect number of clusters, got %v", len(cs))
		}
	}
}

// > go test -c && ./clustering.test -test.bench=PointClusteringGeoProjected -test.cpuprofile=cpu.out -test.benchtime=10s
// > go tool pprof clustering.test cpu.out
func BenchmarkClusterGeoClusters(b *testing.B) {
	clusters, _ := loadPrefilteredTestClusters(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs := ClusterGeoClusters(clusters, 30)

		if len(cs) != 27 {
			b.Fatalf("incorrect number of clusters, got %v", len(cs))
		}
	}
}

func BenchmarkInitClusterDistances(b *testing.B) {
	clusters, _ := loadPrefilteredTestClusters(b)

	bound := geo.NewBoundFromPoints(clusters[0].Centroid, clusters[1].Centroid)
	for _, cluster := range clusters {
		bound.Extend(cluster.Centroid)
		geo.Mercator.Project(cluster.Centroid)
	}
	factor := geo.MercatorScaleFactor(bound.Center().Lat())
	threshold := 30 * 30 * factor * factor

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		initClusterDistances(clusters, CentroidSquaredDistance{}, threshold)
	}
}

func loadPrefilteredTestClusters(tb testing.TB) ([]*Cluster, []geo.Pointer) {
	f, err := os.Open("testdata/prefiltered.json.gz")
	if err != nil {
		tb.Fatalf("unable to open test file %v", err)
	}
	defer f.Close()

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		tb.Fatalf("unable to create gz reader: %v", err)
	}
	defer gzReader.Close()

	var sets [][]*geo.Point
	err = json.NewDecoder(gzReader).Decode(&sets)
	if err != nil {
		tb.Fatalf("could not unmarshal data: %v", err)
	}

	var clusters []*Cluster
	for _, s := range sets {
		var pointers []geo.Pointer
		for _, p := range s {
			pointers = append(pointers, &event{Location: p})
		}

		clusters = append(clusters, NewCluster(pointers...))
	}

	var pointers []geo.Pointer
	for _, c := range clusters {
		pointers = append(pointers, c.Pointers...)
	}

	return clusters, pointers
}

type event struct {
	Location *geo.Point
}

func (e *event) Point() *geo.Point {
	return e.Location
}
