package hashtable

import (
	"sync"
	"sync/atomic"
)

var num = uint64(1000007)

func getHashValue(key uint64) uint64 {
	return key
	return key % uint64(num*1.0/3)
}

type KVpair struct {
	key   uint64
	value int64
}
type Entry struct {
	KV   KVpair
	next *Entry
}

type HashTable struct {
	dirty  map[uint64]*Entry
	length uint64
	read   atomic.Value // readOnly
	mu     sync.Mutex
	misses int
}

// readOnly is an immutable struct stored atomically in the HashTable.read field.
type readOnly struct {
	m       map[uint64]*Entry
	amended bool // true if the dirty map contains some key not in m.
}

func NewHt(length uint64) *HashTable {
	ht := new(HashTable)
	ht.dirty = make(map[uint64]*Entry, length)
	return ht
}

func (ht *HashTable) UnsafePut(hashValue uint64, kv *KVpair) {
	oldEntry := ht.dirty[hashValue]
	newEntry := new(Entry)
	newEntry.KV = *kv
	newEntry.next = (*Entry)(oldEntry)
	ht.length++
	ht.dirty[hashValue] = newEntry
}

func (ht *HashTable) UnsafeCount(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key)
	entry := (*Entry)(ht.dirty[hashValue])
	for entry != nil {
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}

func (ht *HashTable) UnsafePrint() {
	for k, entry := range ht.dirty {
		println("------------", k)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}
}
