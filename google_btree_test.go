package motxn

import "testing"

func BenchmarkGoogleBTree(b *testing.B) {
	benchmarkDB[Int, int](func() DB[Int, int] {
		return NewGoogleBTree[Int, int]()
	}, b)
}
