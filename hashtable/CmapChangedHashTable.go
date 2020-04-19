package hashtable

import (
	"fmt"
)

type ConcMapHashTable struct {
	mp     ConcurrentMap
	length uint64
}

func NewConcMHT(length uint64) *ConcMapHashTable {
	ht := new(ConcMapHashTable)
	ht.mp = NewConcMap()
	ht.length = length
	return ht
}

func (ht *ConcMapHashTable) GetLen() uint64 {
	return ht.length
}

func (ht *ConcMapHashTable) ConcurrentPut(entry *Entry) {
	hashValue:=getHashValue(entry.KV.key,ht.length)
	cb := func(exists bool, valueInMap, newValue *Entry) *Entry {
		if !exists {
			return newValue
		}
		nv := newValue
		nv.next = valueInMap
		return nv
	}
	ht.mp.Upsert(hashValue, entry, cb)

}
func (ht *ConcMapHashTable) Count(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key, ht.length)
	pp, _ := ht.mp.Get((hashValue))
	entry := pp
	for entry != nil {
		//fmt.Println(entry)
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}

func (ht *ConcMapHashTable) Print() {
	ht.mp.IterCb(func(key uint64, value *Entry) {
		for entry := value; entry != nil; entry = entry.next {
			fmt.Print(entry.KV)
		}
		println()
	})
}
