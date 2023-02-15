package motxn

import (
	"sync"
	"sync/atomic"

	"github.com/tidwall/btree"
)

type tidwallBTree[K Ordered[K], V any] struct {
	writeLock sync.Mutex
	tree      atomic.Pointer[btree.BTreeG[tidwallBTreeEntry[K, V]]]
}

type tidwallBTreeEntry[K Ordered[K], V any] struct {
	Key   K
	Value V
}

func compareTidwallBTreeEntry[K Ordered[K], V any](a, b tidwallBTreeEntry[K, V]) bool {
	return a.Key.Less(b.Key)
}

func NewTidwallBTree[K Ordered[K], V any]() *tidwallBTree[K, V] {
	t := btree.NewBTreeG(compareTidwallBTreeEntry[K, V])
	ret := &tidwallBTree[K, V]{}
	ret.tree.Store(t)
	return ret
}

var _ DB[Int, int] = new(tidwallBTree[Int, int])

type tidwallBTreeTransaction[K Ordered[K], V any] struct {
	tree     *tidwallBTree[K, V]
	snapshot *btree.BTreeG[tidwallBTreeEntry[K, V]]
	readOnly bool
}

func (t *tidwallBTree[K, V]) NewTransaction() Transaction[K, V] {
	t.writeLock.Lock()
	return &tidwallBTreeTransaction[K, V]{
		tree:     t,
		snapshot: t.tree.Load().Copy(),
	}
}

func (t *tidwallBTree[K, V]) NewReadOnlyTransaction() ReadOnlyTransaction[K, V] {
	return &tidwallBTreeTransaction[K, V]{
		snapshot: t.tree.Load(),
	}
}

func (t *tidwallBTreeTransaction[K, V]) Commit() {
	t.tree.tree.Store(t.snapshot)
	t.tree.writeLock.Unlock()
}

func (t *tidwallBTreeTransaction[K, V]) Abort() {
	t.tree.writeLock.Unlock()
}

func (t *tidwallBTreeTransaction[K, V]) Get(key K) (value V, ok bool) {
	entry, ok := t.snapshot.Get(tidwallBTreeEntry[K, V]{
		Key: key,
	})
	if !ok {
		return
	}
	value = entry.Value
	return
}

func (t *tidwallBTreeTransaction[K, V]) Set(key K, value V) {
	t.snapshot.Set(tidwallBTreeEntry[K, V]{
		Key:   key,
		Value: value,
	})
}
