package hashtable

import (
	"sync"
)

type LockedHashTable struct {
	 writeMap map[uint64]*Entry
	 length   uint64
	 Lock    sync.Mutex
 }

func NewLHT(length uint64) *LockedHashTable  {
	ht := new(LockedHashTable)
	ht.writeMap = make(map[uint64]*Entry, length)
	ht.length =0
	return ht
}

func (ht *LockedHashTable) ConcurrentPut(hashValue uint64, kv *KVpair) {
	newEntry := new(Entry)
	newEntry.next = nil
	newEntry.KV = *kv
	ht.Lock.Lock()
	oldEntry := ht.writeMap[newEntry.KV.key]
	ht.writeMap[newEntry.KV.key]=newEntry
	newEntry.next = oldEntry
	ht.Lock.Unlock()
}

func (ht *LockedHashTable) Count(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key)
	entry := (*Entry)(ht.writeMap[hashValue])
	for entry != nil {
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}

func (ht *LockedHashTable) Print() {
	for k, entry := range ht.writeMap {
		println("------------", k)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}
}
