package hashtable

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type SyncHashTable struct {
	 writeMap sync.Map
	 length   uint64
 }

func NewSHT(length uint64) *SyncHashTable  {
	ht := new(SyncHashTable)
	ht.length =0
	return ht
}


func (ht *SyncHashTable) ConcurrentPut(hashValue uint64, kv *KVpair) {
	newEntry := new(Entry)
	newEntry.next = nil
	newEntry.KV = *kv
	_, existed :=ht.writeMap.LoadOrStore(hashValue, newEntry)
	if existed {
		for {
			oldp, ok := ht.writeMap.Load(kv.key)
			oldv := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&(oldp.(*Entry).next))))
			if ok && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(oldp.(*Entry).next))), unsafe.Pointer(oldv), unsafe.Pointer(newEntry)) {
				newEntry.next = (*Entry)(oldv)
				break
			}
		}
	}

}

func (ht *SyncHashTable) Count(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key)
	pp,_ :=ht.writeMap.Load(hashValue)
	entry := pp.(*Entry)
	for entry != nil {
		//fmt.Println(entry)
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}else {
			break
		}
		entry = entry.next
	}
	return count
}

func (ht *SyncHashTable) Print() {
	/*for k, entry := range ht.writeMap {
		println("------------", k)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}*/
}
