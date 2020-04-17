package hashtable

import (
	"sync/atomic"
	"unsafe"
)

type ArrayHashTable struct {
	 writeMap []*Entry
	 length   uint64
	 distribution map[int]int
 }

func NewAHT(length uint64) *ArrayHashTable  {
	ht := new(ArrayHashTable)
	ht.writeMap = make([]*Entry, length)
	ht.length = length
	ht.distribution = make(map[int]int,length)
	return ht
}

func (ht *ArrayHashTable) ConcurrentPut(hashValue uint64, kv *KVpair) {
	newEntry := new(Entry)
	newEntry.next = nil
	newEntry.KV = *kv
	for {
		oldEntry := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ht.writeMap[hashValue])))
		ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ht.writeMap[hashValue])), oldEntry, unsafe.Pointer(newEntry))
		if ok {
			newEntry.next = (*Entry)(oldEntry)
			return
		}
	}
}
func (ht *ArrayHashTable) GetLen() uint64 {
	return ht.length
}
func (ht *ArrayHashTable) Count(kv *KVpair) int {
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

func (ht *ArrayHashTable) Print() {
	for i, entry := range ht.writeMap {
		if entry == nil {
			continue
		}
		println("------------", entry.KV.key, i)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}
}
func (ht *ArrayHashTable) Dis() {
	for _, entry := range ht.writeMap {
		count:=0
		if entry == nil {
			continue
		}
		for entry != nil {
			count++
			entry = entry.next
		}
		ht.distribution[count]++
	}
}
func (ht *ArrayHashTable) PrintDis() {
	for k,v:= range ht.distribution {
		println(k,v)
	}
}