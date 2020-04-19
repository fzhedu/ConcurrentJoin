package hashtable
//TODO: a failed try, it cannot just use two map and exchange them, because, the older one may be reading by other threads
type DCHashTable struct {
	HashTable
	delta []*Entry
	cursor int
}

func NewCHT(length uint64) *DCHashTable {
	ht := new(DCHashTable)
	ht.writeMap = make(map[uint64]*Entry, length)
	ht.cursor=0
	ht.delta = make([]*Entry, length, length)
	return ht
}

func (ht *DCHashTable) ConcurrentPut(newEntry *Entry) {
	hashValue:=getHashValue(newEntry.KV.key,ht.length)
	read, _ := ht.read.Load().(readOnly)
	if v, ok := read.m[hashValue]; ok {
		// update an Entry: insert a key into read only map
		v.CASinsert(newEntry)
	} else {
		ht.mu.Lock()
		// recheck
		read, _ = ht.read.Load().(readOnly)
		if v, ok := read.m[hashValue]; ok {
			ht.mu.Unlock()
			// update an Entry: insert a key into read only map
			v.CASinsert(newEntry)
		} else if v, ok := ht.writeMap[hashValue]; ok {
			// only exist in writeMap map
			// update the next of the first node v, i.e., insert a new node after the first node
			oldEntry := v.next
			v.next = newEntry
			newEntry.next = oldEntry
			// a miss occurs
			ht.exchangeAndUpdate(false)
			ht.mu.Unlock()
		} else {
			// not exist in both writeMap and readonly map
			if !read.amended {
				ht.read.Store(readOnly{read.m, true})
			}
			ht.delta[ht.cursor] = newEntry
			ht.cursor++
			ht.writeMap[hashValue] = newEntry
			ht.exchangeAndUpdate(false)
			ht.mu.Unlock()
		}
	}
}

func (ht *DCHashTable) exchangeAndUpdate(forced bool) {
	ht.misses++
	if ht.misses < len(ht.writeMap) && ht.cursor < cap(ht.delta) && !forced {
		return
	}
	if ht.cursor > 0 {
		read, _ := ht.read.Load().(readOnly)
		ht.read.Store(readOnly{ht.writeMap, false})
		// exchange the read and writeMap
		// insert the delta into the new writeMap
		ht.writeMap = nil
		if len(read.m) == 0 {
			ht.writeMap = make(map[uint64]*Entry, cap(ht.delta))
		} else {
			ht.writeMap = read.m
		}
		for i:= 0; i < ht.cursor; i++ {
			ht.writeMap[ht.delta[i].KV.key]=ht.delta[i]
		}
		ht.cursor = 0
		ht.misses = 0
	}
}

func (ht *DCHashTable) Synchronize() {
	// atomically read
	read, _ := ht.read.Load().(readOnly)
	if read.amended {
		ht.mu.Lock()
		ht.exchangeAndUpdate(true)
		ht.mu.Unlock()
	}
}

// TODO: returned value is copied or referenced?
func (ht *DCHashTable) GetABucket(hashValue uint64) (kvPtr []*Entry) {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	entry := (*Entry)(read.m[hashValue])
	for entry != nil {
		kvPtr = append(kvPtr, entry)
		entry = entry.next
	}
	return
}

func (ht *DCHashTable) Find(kv *KVpair) bool {
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

func (ht *DCHashTable) Count(kv *KVpair) int {
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


func (ht *DCHashTable) Print() {
	ht.Synchronize()
	read, _ := ht.read.Load().(readOnly)
	for k, entry := range read.m {
		println("------------", k)
		for entry != nil {
			print(" (", entry.KV.key, "  ", entry.KV.value, ") ")
			entry = entry.next
		}
	}
}
