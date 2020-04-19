package hashtable

import (
	"sync"
	"sync/atomic"
)


func getHashValue(key uint64, seed uint64) uint64 {
	return key % seed
}

type KVpair struct {
	key   uint64
	value int64
}
type Entry struct {
	KV   KVpair
	next *Entry
}

type BaseHashTable interface {
	ConcurrentPut(entry *Entry)
	Count(kv *KVpair) int
	Print()
	GetLen() uint64
}

type HashTable struct {
	writeMap map[uint64]*Entry
	length   uint64
	read     atomic.Value // readOnly
	mu       sync.Mutex
	misses   int
	distribution map[int]int
}

// readOnly is an immutable struct stored atomically in the HashTable.read field.
type readOnly struct {
	m       map[uint64]*Entry
	amended bool // true if the writeMap map contains some key not in m.
}

func NewHt(length uint64) *HashTable {
	ht := new(HashTable)
	ht.writeMap = make(map[uint64]*Entry, length)
	ht.length = length
	ht.distribution = make(map[int]int,length)
	return ht
}

func (ht *HashTable) UnsafePut(entry *Entry) {
	hashValue:= getHashValue(entry.KV.key,ht.length)
	oldEntry := ht.writeMap[hashValue]
	entry.next = (*Entry)(oldEntry)
	ht.writeMap[hashValue] = entry
}

func (ht *HashTable) GetLen() uint64  {
	return ht.length
}

func (ht *HashTable) UnsafeCount(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key,ht.length)
	entry := (*Entry)(ht.writeMap[hashValue])
	for entry != nil {
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}

func (ht *HashTable) UnsafePrint() {
	for _, entry := range ht.writeMap {
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
		println()
	}
	println("---------")
}
func (ht *HashTable) UnsafeDis() {
	for _, entry := range ht.writeMap {
		count:=0
		for entry != nil {
			count++
			entry = entry.next
		}
		if count != 0 {
			ht.distribution[count]++
		}
	}
}
func (ht *HashTable) PrintDis() {
	for key, value := range ht.distribution {
		println(key,value)
	}
}
