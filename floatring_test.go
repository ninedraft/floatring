package floatring_test

import (
	"fmt"
	"testing"

	"github.com/ninedraft/floatring"
)

func TestBufferEmpty(test *testing.T) {
	const N = 100
	var empty = floatring.NewBuffer(0)
	for x := 0.0; x < N; x++ {
		empty.WriteValue(x)
	}
	for i := 0; i < N; i++ {
		var _, ok = empty.ReadValue()
		if ok {
			test.Fatalf("read: expected no values in the buffer")
		}
	}
	empty.ForEach(func(x float64) bool {
		test.Fatalf("ForEach: expected no values in the buffer")
		return true
	})
	var b = empty.DumpTo(nil)
	if len(b) != 0 {
		test.Fatalf("WriteTo: expected no values in the buffer")
	}
	if empty.Len() != 0 {
		test.Fatalf("Len: expected no values in the buffer")
	}
	if empty.Cap() != 0 {
		test.Fatalf("Cap: expected no values in the buffer")
	}
}

func TestBufferConsistency(test *testing.T) {
	testBufferConsistency(test, 1, 1)
	testBufferConsistency(test, 1, 50)
	testBufferConsistency(test, 100, 50)
	testBufferConsistency(test, 100, 51)
	testBufferConsistency(test, 100, 200)
	testBufferConsistency(test, 100, 0)
}

func testBufferConsistency(test *testing.T, n, size int) {
	var name = fmt.Sprintf("values=%d  size=%d", n, size)
	test.Run(name, func(test *testing.T) {
		var values = make([]float64, n)
		var ring = floatring.NewBuffer(size)
		for i := 0; i < n; i++ {
			var x = float64(i)
			values[i] = x
			ring.WriteValue(x)
		}
		if !ring.IsFull() && !ring.IsEmpty() && n >= size {
			test.Errorf("buffer is not full: len=%d", ring.Len())
			return
		}
		var got = ring.DumpTo(nil)
		var tail = n - size
		switch {
		case size > n:
			tail = 0
		case size == 0:
			tail = n
		}
		var expected = values[tail:]
		if !eq(expected, got) {
			test.Errorf("expected %d values %+v, got %d: %+v", len(expected), expected, len(got), got)
			return
		}
	})
}

func eq(aa, bb []float64) bool {
	if len(aa) != len(bb) {
		return false
	}
	for i, a := range aa {
		if bb[i] != a {
			return false
		}
	}
	return true
}
