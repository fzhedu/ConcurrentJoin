package hashtable

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)
// use go.sycn.Map

type SyncHashTable struct {
	 writeMap sync.Map
	 length   uint64
 }

func NewSHT(length uint64) *SyncHashTable  {
	ht := new(SyncHashTable)
	ht.length =length
	return ht
}

func (ht *SyncHashTable) GetLen() uint64 {
	return ht.length
}
func (ht *SyncHashTable) ConcurrentPut(entry *Entry) {
	hashValue:=getHashValue(entry.KV.key,ht.length)
	oldp, existed :=ht.writeMap.LoadOrStore(hashValue, entry)
	if existed {
		for {
			//oldp, ok := ht.writeMap.Load(hashValue)
			oldv := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&(oldp.(*Entry).next))))
			if  atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(oldp.(*Entry).next))), unsafe.Pointer(oldv), unsafe.Pointer(entry)) {
				entry.next = (*Entry)(oldv)
				break
			}
		}
	}
}

func (ht *SyncHashTable) Count(kv *KVpair) int {
	count := 0
	hashValue := getHashValue(kv.key,ht.length)
	pp,_ :=ht.writeMap.Load(hashValue)
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

func (ht *SyncHashTable) Print() {
	ht.writeMap.Range(func(key, value interface{}) bool {
		fmt.Println(value.(*Entry).KV)
		return true
	})
}
