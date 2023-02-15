package motxn

import "testing"

func BenchmarkTidwallBTree(b *testing.B) {
	benchmarkDB[Int, int](func() DB[Int, int] {
		return NewTidwallBTree[Int, int]()
	}, b)
}
