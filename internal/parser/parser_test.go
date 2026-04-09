package parser

import (
	"testing"
)

const eachOutput = `goos: linux
goarch: amd64
pkg: github.com/bold-minds/each
cpu: Intel(R) Core(TM) Ultra 9 275HX
BenchmarkFind_Hit-24           	 1492314	       764.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkFind_Miss-24          	  707406	      1549 ns/op	       0 B/op	       0 allocs/op
BenchmarkFilter_Half-24        	   98682	     11607 ns/op	   49152 B/op	       1 allocs/op
BenchmarkGroupBy-24            	   30714	     39332 ns/op	  154273 B/op	      20 allocs/op
BenchmarkKeyBy-24              	   35157	     33012 ns/op	  131153 B/op	       5 allocs/op
BenchmarkPartition-24          	   68347	     16801 ns/op	   98304 B/op	       2 allocs/op
BenchmarkCount-24              	  767194	      1513 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvery_AllTrue-24      	  803811	      1504 ns/op	       0 B/op	       0 allocs/op
BenchmarkEvery_EarlyFail-24    	765236575	         1.473 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/bold-minds/each	13.488s
`

const multiPkgOutput = `goos: linux
goarch: amd64
pkg: github.com/example/repo/internal/foo
BenchmarkFoo-8   	 1000000	      1234 ns/op	     128 B/op	       2 allocs/op
PASS
ok  	github.com/example/repo/internal/foo	1.5s
pkg: github.com/example/repo/internal/bar
BenchmarkBar-8   	  500000	      5678 ns/op	     256 B/op	       4 allocs/op
BenchmarkBaz-8   	  250000	      9999 ns/op	     512 B/op	       8 allocs/op
PASS
ok  	github.com/example/repo/internal/bar	3.2s
`

const noBopOutput = `goos: linux
goarch: amd64
pkg: github.com/example/repo
BenchmarkSimple-8   	 1000000	       500.0 ns/op
PASS
ok  	github.com/example/repo	0.5s
`

func TestParse_SinglePackage(t *testing.T) {
	results, err := Parse(eachOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 package, got %d", len(results))
	}

	pkg := results[0]
	if pkg.PackagePath != "github.com/bold-minds/each" {
		t.Errorf("expected package github.com/bold-minds/each, got %q", pkg.PackagePath)
	}
	if len(pkg.Results) != 9 {
		t.Fatalf("expected 9 benchmarks, got %d", len(pkg.Results))
	}

	// Spot check first and last
	first := pkg.Results[0]
	if first.Name != "BenchmarkFind_Hit-24" {
		t.Errorf("expected BenchmarkFind_Hit-24, got %q", first.Name)
	}
	if first.NsOp != 764.5 {
		t.Errorf("expected 764.5 ns/op, got %f", first.NsOp)
	}
	if first.BOp != 0 {
		t.Errorf("expected 0 B/op, got %d", first.BOp)
	}

	last := pkg.Results[8]
	if last.Name != "BenchmarkEvery_EarlyFail-24" {
		t.Errorf("expected BenchmarkEvery_EarlyFail-24, got %q", last.Name)
	}
	if last.NsOp != 1.473 {
		t.Errorf("expected 1.473 ns/op, got %f", last.NsOp)
	}

	// Check a benchmark with non-zero allocations
	groupBy := pkg.Results[3]
	if groupBy.BOp != 154273 {
		t.Errorf("expected 154273 B/op, got %d", groupBy.BOp)
	}
	if groupBy.AllocsOp != 20 {
		t.Errorf("expected 20 allocs/op, got %d", groupBy.AllocsOp)
	}
}

func TestParse_MultiPackage(t *testing.T) {
	results, err := Parse(multiPkgOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(results))
	}

	if results[0].PackagePath != "github.com/example/repo/internal/foo" {
		t.Errorf("wrong first package: %q", results[0].PackagePath)
	}
	if len(results[0].Results) != 1 {
		t.Errorf("expected 1 result in foo, got %d", len(results[0].Results))
	}

	if results[1].PackagePath != "github.com/example/repo/internal/bar" {
		t.Errorf("wrong second package: %q", results[1].PackagePath)
	}
	if len(results[1].Results) != 2 {
		t.Errorf("expected 2 results in bar, got %d", len(results[1].Results))
	}
}

func TestParse_NoBopOutput(t *testing.T) {
	results, err := Parse(noBopOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 package, got %d", len(results))
	}
	if len(results[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results[0].Results))
	}

	r := results[0].Results[0]
	if r.NsOp != 500.0 {
		t.Errorf("expected 500.0 ns/op, got %f", r.NsOp)
	}
	if r.BOp != 0 {
		t.Errorf("expected 0 B/op for missing field, got %d", r.BOp)
	}
	if r.AllocsOp != 0 {
		t.Errorf("expected 0 allocs/op for missing field, got %d", r.AllocsOp)
	}
}

func TestParse_EmptyOutput(t *testing.T) {
	results, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 packages, got %d", len(results))
	}
}

func TestParse_OrderPreserved(t *testing.T) {
	results, err := Parse(multiPkgOutput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// foo should come before bar (order of appearance)
	if results[0].PackagePath != "github.com/example/repo/internal/foo" {
		t.Errorf("expected foo first, got %q", results[0].PackagePath)
	}
	if results[1].PackagePath != "github.com/example/repo/internal/bar" {
		t.Errorf("expected bar second, got %q", results[1].PackagePath)
	}
}
