package hashtable

import (
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
)

type Workload struct {
	KV     []KVpair
	cursor uint64
	Length uint64
	Step   uint64
	times []int
}

func NewW(num uint64) *Workload  {
	w:=new(Workload)
	w.GenLoad(num)
	return w
}

func (w *Workload) Reset () {
	w.cursor = 0
}

func (w *Workload)GenLoad(num uint64) {
	w.cursor =0;
	w.Length = num
	w.Step = 10000
	w.times = make([]int,num)
	for i := uint64(0); i < num; i++ {
		w.KV = append(w.KV, KVpair{rand.Uint64()% num , rand.Int63()})
	}
	for i := uint64(0); i < num; i++ {
		w.times[i] = 0
		for j:= uint64(0); j < num; j++ {
			if w.KV[i].key == w.KV[j].key && w.KV[i].value == w.KV[j].value {
				w.times[i]++
			}
		}
	}
	return
}
func (w *Workload) Read() (uint64, uint64)  {
	old := atomic.AddUint64(&w.cursor,w.Step)
	if old < w.Length {
		return  old - w.Step, old
	} else if old > w.Length  && old < w.Length + w.Step{
		return  old - w.Step, w.Length
	} else {
		return 0, 0
	}
}

func BenchamrkUnsafeHT(w *Workload, time int) {
	w.Reset()
	ht := NewHt( w.Length)
	for i := uint64(0); i < w.Length; i++ {
		ht.UnsafePut(getHashValue((*w).KV[i].key,ht.length), &(*w).KV[i])
	}
//	ht.UnsafePrint()
	ok := UnsafeCheck(w, ht)
	if !ok {
		println("ERROR occor")
		os.Exit(1)
	} else {
		//println("ok")
	}
}
func BenchamrkCHT(w *Workload, time int) {
	w.Reset()
	var ht BaseHashTable
	ht = NewHt(w.Length)
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, &ht)
		}()
	}
	wg.Wait()
	//ht.Print()
	ok := Check(w, &ht)
	if !ok {
		println("ERROR occor")
		os.Exit(1)
	} else {
		//println("ok")
	}
}
func BenchamrkSCHT(w *Workload, time int) {
	w.Reset()
	var ht BaseHashTable
	ht = NewSHT(w.Length)
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
	//	ht.Print()
	ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
		} else {
			//os.Exit(0)
			//println("ok")
		}
}

func BenchamrkACHT(w *Workload, time int) {
	w.Reset()
	var ht BaseHashTable
	ht = NewAHT(w.Length)
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
//	ht.Print()
	ok := Check(w, &ht)
	if !ok {
		println("ERROR occor")
		os.Exit(1)
	} else {
		//os.Exit(0)
//		println("ok")
	}
}
func BenchamrkLCHT(w *Workload, time int) {
	w.Reset()
	var ht BaseHashTable
	ht = NewLHT(w.Length)
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
//	ht.Print()
	ok := Check(w, &ht)
	if !ok {
		println("ERROR occor")
		os.Exit(1)
	} else {
	//	println("ok")
	}
}

func putLoad(w *Workload, ht *BaseHashTable) {
	for i, end := w.Read(); i < end; i++ {
		(*ht).ConcurrentPut(getHashValue((*w).KV[i].key,(*ht).GetLen()), &(*w).KV[i])
	}
	//println("succeed input")
}
func UnsafeCheck(w *Workload, ht *HashTable) bool {
	for i := uint64(0); i < w.Length; i++ {
		if c:=ht.UnsafeCount(&(*w).KV[i]) ; c!= w.times[i] {
			println("ERROR: time error ", (*w).KV[i].key, " actural time = ",c," time = ", w.times[i])
			return false
		}
	}
	return true
}
func Check(w *Workload,ht *BaseHashTable) bool {
	for i := uint64(0); i < w.Length; i++ {
		if c:=(*ht).Count(&(*w).KV[i]) ; c!= w.times[i] {
			println("ERROR: time error ", (*w).KV[i].key, " actural time = ",c," time = ", w.times[i])
			return false
		}
	}
	return true
}
