package hashtable

import (
	"sync/atomic"
	"unsafe"
)

func (e *Entry) CASinsert(newEntry *Entry) {
	for {
		oldEntry := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&e.next)))
		ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&e.next)), oldEntry, unsafe.Pointer(newEntry))
		if ok {
			newEntry.next = (*Entry)(oldEntry)
			return
		}
	}
}


func (ht *HashTable) ConcurrentPut(entry *Entry) {
	hashValue :=getHashValue(entry.KV.key,ht.length)
	read, _ := ht.read.Load().(readOnly)
	if v, ok := read.m[hashValue]; ok {
		// update an Entry: insert a key into read only map
		v.CASinsert(entry)
	} else {
		ht.mu.Lock()
		// recheck
		read, _ = ht.read.Load().(readOnly)
		if v, ok := read.m[hashValue]; ok {
			ht.mu.Unlock()
			// update an Entry: insert a key into read only map
			v.CASinsert(entry)
		} else if v, ok := ht.writeMap[hashValue]; ok {
			// only exist in writeMap map
			ht.writeMap[hashValue] = entry
			entry.next = v
			// a miss occurs
			ht.missLocked()
			ht.mu.Unlock()
		} else {
			// not exist in both writeMap and readonly map
			if !read.amended {
				ht.dirtyLocked()
				ht.read.Store(readOnly{read.m, true})
			}
			ht.writeMap[hashValue] = entry
			ht.mu.Unlock()
		}
	}
}

func (ht *HashTable) missLocked() {
	ht.misses++
	if ht.misses < len(ht.writeMap) {
		return
	}
	ht.read.Store(readOnly{ht.writeMap, false})
	ht.writeMap = nil
	ht.misses = 0
}

func (ht *HashTable) dirtyLocked() {
	if ht.writeMap != nil {
		return
	}
	read, _ := ht.read.Load().(readOnly)
	// TODO: avoid this fully copy
	ht.writeMap = make(map[uint64]*Entry, len(read.m))
	for k, e := range read.m {
		ht.writeMap[k] = e
	}
}

func (ht *HashTable) Synchronize() {
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		// m.writeMap contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the writeMap copy
		// immediately!
		ht.mu.Lock()
		read, _ = ht.read.Load().(readOnly)
		if read.amended {
			read = readOnly{ ht.writeMap,false}
			ht.read.Store(read)
			ht.writeMap = nil
			ht.misses = 0
		}
		ht.mu.Unlock()
	}
}

// TODO: returned value is copied or referenced?
func (ht *HashTable) GetABucket(hashValue uint64) (kvPtr []*Entry) {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	entry := (*Entry)(read.m[hashValue])
	for entry != nil {
		kvPtr = append(kvPtr, entry)
		entry = entry.next
	}
	return
}

func (ht *HashTable) Find(kv *KVpair) bool {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	hashValue := getHashValue(kv.key,ht.length)
	entry := (*Entry)(read.m[hashValue])
	for entry != nil {
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			return true
		}
		entry = entry.next
	}
	return false
}

func (ht *HashTable) Count(kv *KVpair) int {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	count := 0
	hashValue := getHashValue(kv.key,ht.length)
	entry := (*Entry)(read.m[hashValue])
	for entry != nil {
		if entry.KV.key == kv.key && entry.KV.value == kv.value {
			count = count + 1
		}
		entry = entry.next
	}
	return count
}


func (ht *HashTable) Print() {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	for _, entry := range read.m {
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
		println()
	}
}
