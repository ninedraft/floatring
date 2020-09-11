package floatring_test

import (
	"testing"

	"github.com/ninedraft/floatring"
)

var buf1024 = floatring.NewBuffer(1024)

func BenchmarkReadWrite1M(bench *testing.B) {
	for i := 0; i < 1000000; i++ {
		var x = float64(i)
		buf1024.WriteValue(x)
		if i%2 == 0 {
			buf1024.ReadValue()
		}
	}
}
