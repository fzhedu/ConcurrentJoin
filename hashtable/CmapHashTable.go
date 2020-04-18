package hashtable

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
)

type CmapHashTable struct {
	mp     cmap.ConcurrentMap
	length uint64
}

func NewCMHT(length uint64) *CmapHashTable {
	ht := new(CmapHashTable)
	ht.mp = cmap.New()
	ht.length = length
	return ht
}

func (ht *CmapHashTable) GetLen() uint64 {
	return ht.length
}

func (ht *CmapHashTable) ConcurrentPut(hashValue uint64, kv *KVpair) {
	newEntry := new(Entry)
	newEntry.next = nil
	newEntry.KV = *kv
	cb := func(exists bool, valueInMap interface{}, newValue interface{}) interface{} {
		if !exists {
			return newValue
		}
		nv := newValue.(*Entry)
		nv.next = valueInMap.(*Entry)
		return nv
	}
	ht.mp.Upsert(string(hashValue), newEntry, cb)

}
func (ht *CmapHashTable) Count(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key, ht.length)
	pp, _ := ht.mp.Get(string(hashValue))
	entry := pp.(*Entry)
	for entry != nil {
		//fmt.Println(entry)
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}

func (ht *CmapHashTable) Print() {
	ht.mp.IterCb(func(key string, value interface{}) {
		for entry := value.(*Entry); entry != nil; entry = entry.next {
			fmt.Print(entry.KV)
		}
		println()
	})
}
