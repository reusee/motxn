package motxn

import "testing"

func benchmarkDB[K, V any](
	newDB func() DB[Int, int],
	b *testing.B,
) {

	b.Run("single transaction set/get single key", func(b *testing.B) {
		db := newDB()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tx := db.NewTransaction()
			tx.Set(1, 42)
			tx.Commit()
			v, ok := tx.Get(1)
			if !ok {
				b.Fatal()
			}
			if v != 42 {
				b.Fatal()
			}
		}
	})

}
