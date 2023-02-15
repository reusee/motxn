package motxn

type Ordered[T any] interface {
	Less(than T) bool
}

type DB[K Ordered[K], V any] interface {
	NewTransaction() Transaction[K, V]
	NewReadOnlyTransaction() ReadOnlyTransaction[K, V]
}

type ReadOnlyTransaction[K Ordered[K], V any] interface {
	Get(K) (V, bool)
}

type Transaction[K Ordered[K], V any] interface {
	ReadOnlyTransaction[K, V]

	Set(K, V)
	Commit()
	Abort()
}
