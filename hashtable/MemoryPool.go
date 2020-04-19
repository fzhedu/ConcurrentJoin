package hashtable

const (
	initialEntrySliceLen = 64
	maxEntrySliceLen     = 8 * 1024
)

type EntryStore struct {
	slices [][]Entry
	cursor int
}

func NewStore() *EntryStore {
	es := new(EntryStore)
	es.slices = [][]Entry{make([]Entry, initialEntrySliceLen, initialEntrySliceLen)}
	es.cursor =0
	return  es
}

func (es *EntryStore) GetAddr() (entry*Entry) {
	sliceIdx := uint32(len(es.slices) - 1)
	slice := es.slices[sliceIdx]
	if es.cursor >= cap(slice) {
		size := cap(slice) * 2
		if size >= maxEntrySliceLen {
			size = maxEntrySliceLen
		}
		slice = make([]Entry, size, size)
		es.slices = append(es.slices, slice)
		sliceIdx++
		es.cursor=0
	}
	entry = &es.slices[sliceIdx][es.cursor]
	es.cursor++
	return
}


