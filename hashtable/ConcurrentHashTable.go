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

func (ht *HashTable) ConcurrentPut(hashValue uint64, kv *KVpair) {
	newEntry := new(Entry)
	newEntry.next = nil
	newEntry.KV = *kv
	read, _ := ht.read.Load().(readOnly)
	if v, ok := read.m[newEntry.KV.key]; ok {
		// update an Entry: insert a key into read only map
		v.CASinsert(newEntry)
	} else {
		ht.mu.Lock()
		// recheck
		read, _ = ht.read.Load().(readOnly)
		if v, ok := read.m[newEntry.KV.key]; ok {
			ht.mu.Unlock()
			// update an Entry: insert a key into read only map
			v.CASinsert(newEntry)
		} else if v, ok := ht.dirty[newEntry.KV.key]; ok {
			// only exist in dirty map
			ht.dirty[newEntry.KV.key] = newEntry
			newEntry.next = v
			// a miss occurs
			ht.missLocked()
			ht.mu.Unlock()
		} else {
			// not exist in both dirty and readonly map
			if !read.amended {
				ht.dirtyLocked()
				ht.read.Store(readOnly{read.m, true})
			}
			ht.dirty[newEntry.KV.key] = newEntry
			ht.mu.Unlock()
		}
	}
}

func (ht *HashTable) missLocked() {
	ht.misses++
	if ht.misses < len(ht.dirty) {
		return
	}
	ht.read.Store(readOnly{ht.dirty, false})
	ht.dirty = nil
	ht.misses = 0
}

func (ht *HashTable) dirtyLocked() {
	if ht.dirty != nil {
		return
	}
	read, _ := ht.read.Load().(readOnly)
	// TODO: avoid this fully copy
	ht.dirty = make(map[uint64]*Entry, len(read.m))
	for k, e := range read.m {
		ht.dirty[k] = e
	}
}

// TODO: returned value is copied or referenced?
func (ht *HashTable) GetABucket(hashValue uint64) (kvPtr []*Entry) {
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		ht.mu.Lock()
		read, _ = ht.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: ht.dirty}
			ht.read.Store(read)
			ht.dirty = nil
			ht.misses = 0
		}
		ht.mu.Unlock()
	}
	entry := (*Entry)(read.m[hashValue])
	for entry != nil {
		kvPtr = append(kvPtr, entry)
		entry = entry.next
	}
	return
}

func (ht *HashTable) Find(kv *KVpair) bool {
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		ht.mu.Lock()
		read, _ = ht.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: ht.dirty}
			ht.read.Store(read)
			ht.dirty = nil
			ht.misses = 0
		}
		ht.mu.Unlock()
	}
	hashValue := getHashValue(kv.key)
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
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		ht.mu.Lock()
		read, _ = ht.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: ht.dirty}
			ht.read.Store(read)
			ht.dirty = nil
			ht.misses = 0
		}
		ht.mu.Unlock()
	}
	count := 0
	hashValue := getHashValue(kv.key)
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
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		ht.mu.Lock()
		read, _ = ht.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: ht.dirty}
			ht.read.Store(read)
			ht.dirty = nil
			ht.misses = 0
		}
		ht.mu.Unlock()
	}
	for k, entry := range read.m {
		println("------------", k)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}
}
