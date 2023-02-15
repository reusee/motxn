package motxn

import (
	"sync"
	"sync/atomic"

	"github.com/google/btree"
)

type googleBTree[K Ordered[K], V any] struct {
	writeLock sync.Mutex
	tree      atomic.Pointer[btree.BTreeG[googleBTreeEntry[K, V]]]
}

type googleBTreeEntry[K Ordered[K], V any] struct {
	Key   K
	Value V
}

func comparegoogleBTreeEntry[K Ordered[K], V any](a, b googleBTreeEntry[K, V]) bool {
	return a.Key.Less(b.Key)
}

func NewGoogleBTree[K Ordered[K], V any]() *googleBTree[K, V] {
	t := btree.NewG(2, comparegoogleBTreeEntry[K, V])
	ret := &googleBTree[K, V]{}
	ret.tree.Store(t)
	return ret
}

var _ DB[Int, int] = new(googleBTree[Int, int])

type googleBTreeTransaction[K Ordered[K], V any] struct {
	tree     *googleBTree[K, V]
	snapshot *btree.BTreeG[googleBTreeEntry[K, V]]
	readOnly bool
}

func (t *googleBTree[K, V]) NewTransaction() Transaction[K, V] {
	t.writeLock.Lock()
	return &googleBTreeTransaction[K, V]{
		tree:     t,
		snapshot: t.tree.Load().Clone(),
	}
}

func (t *googleBTree[K, V]) NewReadOnlyTransaction() ReadOnlyTransaction[K, V] {
	return &googleBTreeTransaction[K, V]{
		snapshot: t.tree.Load(),
	}
}

func (t *googleBTreeTransaction[K, V]) Commit() {
	t.tree.tree.Store(t.snapshot)
	t.tree.writeLock.Unlock()
}

func (t *googleBTreeTransaction[K, V]) Abort() {
	t.tree.writeLock.Unlock()
}

func (t *googleBTreeTransaction[K, V]) Get(key K) (value V, ok bool) {
	entry, ok := t.snapshot.Get(googleBTreeEntry[K, V]{
		Key: key,
	})
	if !ok {
		return
	}
	value = entry.Value
	return
}

func (t *googleBTreeTransaction[K, V]) Set(key K, value V) {
	t.snapshot.ReplaceOrInsert(googleBTreeEntry[K, V]{
		Key:   key,
		Value: value,
	})
}
