package floatring_test

import (
	"fmt"
	"testing"

	"github.com/ninedraft/floatring"
)

func TestBufferEmpty(test *testing.T) {
	const N = 100
	var empty = floatring.New(0)
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
	for n := 0; n < 100; n++ {
		for size := 0; size < 100; size++ {
			testBufferConsistency(test, n, size)
		}
	}
}

func testBufferConsistency(test *testing.T, n, size int) {
	var name = fmt.Sprintf("values=%d  size=%d", n, size)
	test.Run(name, func(test *testing.T) {
		var values = make([]float64, n)
		var ring = floatring.New(size)
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

func TestBufferLen(test *testing.T) {
	var t = func(size, expected, nWrite, nRead int) {
		var name = fmt.Sprintf("size=%d  expected=%d  nWrite=%d  nRead=%d", size, expected, nWrite, nRead)
		test.Run(name, func(test *testing.T) {
			var buf = floatring.New(size)
			for v := 0; v < nWrite; v++ {
				buf.WriteValue(float64(v))
			}
			for v := 0; v < nRead; v++ {
				buf.ReadValue()
			}
			if buf.Len() != expected {
				test.Errorf("expected len=%d, got %d", expected, buf.Len())
			}
		})
	}

	t(0, 0, 100, 0)
	t(1, 1, 100, 0)
	t(100, 100, 100, 0)
	t(11, 11, 31, 0)
	t(10, 5, 5, 0)
	t(10, 6, 15, 4)
}

func TestBatchMethods(test *testing.T) {
	var buf = floatring.New(16)
	var values = make([]float64, 40)
	for i := range values {
		values[i] = float64(i)
	}
	buf.Write(values)

	var short = make([]float64, 4)
	var expectedShort = []float64{24, 25, 26, 27}
	buf.Read(short)
	if !eq(short, expectedShort) {
		test.Errorf("expected %#v, got %#v", expectedShort, short)
	}

	var long = make([]float64, 12)
	var expectedLong = []float64{28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	buf.Read(long)
	if !eq(short, expectedShort) {
		test.Errorf("expected %#v, got %#v", expectedLong, long)
	}
}

func TestForEach(test *testing.T) {
	var t = func(size, input int) {
		var name = fmt.Sprintf("size=%d  input=%d", size, input)
		test.Run(name, func(t *testing.T) {
			var buf = floatring.New(size)
			var values = testValues(input)
			buf.Write(values)
			var got []float64
			buf.ForEach(func(x float64) bool {
				got = append(got, x)
				return true
			})
			if size < input {
				values = values[input-size:]
			}
			if !eq(values, got) {
				test.Errorf("expected %+v got %+v", values, got)
			}
		})
	}

	for size := 0; size < 100; size++ {
		for input := 0; input < 100; input++ {
			t(size, input)
		}
	}
}

func testValues(n int) []float64 {
	var values = make([]float64, n)
	for i := range values {
		values[i] = float64(i)
	}
	return values
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
