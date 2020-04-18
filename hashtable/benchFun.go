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
	distribution map[int]int
}


func NewW(num, Gstep,Gmod uint64) *Workload  {
	w:=new(Workload)
	w.GenLoad(num,Gstep,Gmod,false)
	w.PrintDis()
	return w
}

func (w *Workload) Reset () {
	w.cursor = 0
}
func (w *Workload)PrintDis()  {
	for k,v := range w.distribution {
		println(k,"  ",v)
	}
}
func (w *Workload)GenLoad(num uint64, Gstep, Gmod uint64, check bool) {
	w.cursor =0;
	w.Length = num
	w.Step = Gstep
	w.times = make([]int,num)
	w.distribution = make(map[int]int, num)
	for i := uint64(0); i < num; i++ {
		w.KV = append(w.KV, KVpair{rand.Uint64()% Gmod , rand.Int63()})
	}
	if check {
		for i := uint64(0); i < num; i++ {
			w.times[i] = 0
			for j:= uint64(0); j < num; j++ {
			if w.KV[i].key == w.KV[j].key && w.KV[i].value == w.KV[j].value {
				w.times[i]++
			}
		}
		}
		for i := uint64(0); i < num; i++ {
			w.distribution[w.times[i]]++
		}
	}
	return
}
func (w *Workload) Read() (uint64, uint64)  {
	old := atomic.AddUint64(&w.cursor,w.Step)
	if old <= w.Length {
		return  old - w.Step, old
	} else if old > w.Length  && old - w.Step < w.Length{
		return  old - w.Step, w.Length
	} else {
		return w.Length, w.Length
	}
}

func BenchamrkUnsafeHT(w *Workload, time int, check bool) {
	w.Reset()
	ht := NewHt( w.Length)
	for i := uint64(0); i < w.Length; i++ {
		ht.UnsafePut(getHashValue((*w).KV[i].key,ht.length), &(*w).KV[i])
	}
	if time > 0 {
		ht.UnsafeDis()
		ht.PrintDis()
	}
	if check {
		ht.UnsafePrint()
		ok := UnsafeCheck(w, ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//println("ok")
		} else {
		}
	}
}
func BenchamrkCHT(w *Workload, time int, check bool) {
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
	if check {
		println("all thread Done")
		ht.Print()
		ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//os.Exit(0)
		} else {
			println("ok")
		}
	}
}
func BenchamrkSCHT(w *Workload, time int, check bool) {
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
	if check {
		println("all thread Done")
		ht.Print()
		ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//os.Exit(0)
		} else {
			println("ok")
		}
	}
}
func BenchamrkCMHT(w *Workload, time int, check bool) {
	w.Reset()
	var ht BaseHashTable
	ht = NewCMHT(w.Length)
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, &ht)
		}()
	}
	wg.Wait()
	if check {
		println("all thread Done")
	//	ht.Print()
		ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//os.Exit(0)
		} else {
			println("ok")
		}
	}
}
func BenchamrkACHT(w *Workload, time int, dis,check bool) {
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

	if dis {
		ht.(*ArrayHashTable).Dis()
		ht.(*ArrayHashTable).PrintDis()
	}
	if check {
		println("all thread Done")
		ht.Print()
		ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//os.Exit(0)
		} else {
			println("ok")
		}
	}
}
func BenchamrkLCHT(w *Workload, time int, check bool) {
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
	if check {
		println("all thread Done")
		ht.Print()
		ok := Check(w, &ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
			//os.Exit(0)
		} else {
			println("ok")
		}
	}
}

func putLoad(w *Workload, ht *BaseHashTable) {

	for begin,end := w.Read(); begin < w.Length; begin,end = w.Read(){
		for i:= begin; i < end ; i++ {
			(*ht).ConcurrentPut(getHashValue(w.KV[i].key,(*ht).GetLen()), &(w.KV[i]))
		}
	}
	//println("succeed input")
}
func UnsafeCheck(w *Workload, ht *HashTable) bool {
	return true
	for i := uint64(0); i < w.Length; i++ {
		if c:=ht.UnsafeCount(&(*w).KV[i]) ; c!= w.times[i] {
			println("ERROR: time error ", (*w).KV[i].key, " actural time = ",c," time = ", w.times[i])
			return false
		}
	}
	return true
}
func Check(w *Workload,ht *BaseHashTable) bool {
	//return  true
	for i := uint64(0); i < w.Length; i++ {
		if c:=(*ht).Count(&w.KV[i]) ; c!= w.times[i] {
			println("ERROR: time error ", w.KV[i].key, " actural time = ",c," time = ", w.times[i])
			return false
		}
	}
	return true
}
