package motxn

import "testing"

func benchmarkDB[K, V any](
	newDB func() DB[Int, int],
	b *testing.B,
) {

	b.Run("single transaction set and get single key", func(b *testing.B) {
		db := newDB()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tx := db.NewTransaction()
			tx.Set(1, 42)
			v, ok := tx.Get(1)
			if !ok {
				b.Fatal()
			}
			if v != 42 {
				b.Fatal()
			}
			tx.Commit()
		}
	})

	b.Run("single transaction get single key", func(b *testing.B) {
		db := newDB()
		tx := db.NewTransaction()
		tx.Set(1, 42)
		tx.Commit()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tx := db.NewReadOnlyTransaction()
			v, ok := tx.Get(1)
			if !ok {
				b.Fatal()
			}
			if v != 42 {
				b.Fatal()
			}
		}
	})

	b.Run("single transaction set and get multi key", func(b *testing.B) {
		db := newDB()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tx := db.NewTransaction()
			key := Int(i)
			tx.Set(key, i)
			v, ok := tx.Get(key)
			if !ok {
				b.Fatal()
			}
			if v != i {
				b.Fatal()
			}
			tx.Commit()
		}
	})

	b.Run("parallel transaction set and get multi key", func(b *testing.B) {
		db := newDB()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			tx := db.NewTransaction()
			defer tx.Commit()
			for i := 0; pb.Next(); i++ {
				key := Int(i)
				tx.Set(key, i)
				v, ok := tx.Get(key)
				if !ok {
					b.Fatal()
				}
				if v != i {
					b.Fatal()
				}
			}
		})
	})

	b.Run("parallel transaction get single key", func(b *testing.B) {
		db := newDB()
		tx := db.NewTransaction()
		tx.Set(1, 42)
		tx.Commit()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			tx := db.NewReadOnlyTransaction()
			for i := 0; pb.Next(); i++ {
				v, ok := tx.Get(1)
				if !ok {
					b.Fatal()
				}
				if v != 42 {
					b.Fatal()
				}
			}
		})
	})

}
