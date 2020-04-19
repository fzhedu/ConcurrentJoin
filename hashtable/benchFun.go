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
	times map[KVpair]int
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
	println("---------------workload---------------")
}
func (w *Workload)GenLoad(num uint64, Gstep, Gmod uint64, check bool) {
	w.cursor =0;
	w.Length = num
	w.Step = Gstep
	w.times = make(map[KVpair]int,num)
	w.distribution = make(map[int]int, num)
	for i := uint64(0); i < num; i++ {
		kv :=KVpair{rand.Uint64()% Gmod , rand.Int63()}
		w.KV = append(w.KV, kv)
		w.times[kv]++
	}
	if check {
		for _,time := range w.times{
			w.distribution[time]++
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

func BenchamrkUnsafeHT(w *Workload, time int, check bool) *HashTable {
	w.Reset()
	ht := NewHt( w.Length)
	es :=NewStore()
	for i := uint64(0); i < w.Length; i++ {
		tmp:=es.GetAddr()
		tmp.KV=w.KV[i]
		tmp.next=nil
		ht.UnsafePut(tmp)
	}
	if time > 0 {
		ht.UnsafeDis()
		ht.PrintDis()
	}
	if check {
		//ht.UnsafePrint()
		ok := UnsafeCheck(w, ht)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
		} else {
			//println("ok")
		}
	}
	return ht
}
func benchamrk(ht *BaseHashTable,w *Workload, time int, check bool) {
	w.Reset()
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(w, ht)
		}()
	}
	wg.Wait()
	
	if check {
		wg.Add(time)
		w.Reset()
		//println("all thread Done")
		//ht.Print()
		for t := 0; t < time; t++ {
			go func() {
				defer wg.Done()
				ok := Check(w, ht)
				if !ok {
					println("ERROR occor")
					os.Exit(1)
					//os.Exit(0)
				} else {
					//println("ok")
				}
			}()
		}
	}
	wg.Wait()
}

func BenchamrkCHT(w *Workload, time int, check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewHt(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht
}
func BenchamrkRead(ht *BaseHashTable,w *Workload, time int) {
	wg := &sync.WaitGroup{}
	wg.Add(time)
	w.Reset()
	//println("all thread Done")
	//ht.Print()
	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			ok := ConcCheck(w, ht)
			if !ok {
				println("ERROR occor")
				os.Exit(1)
				//os.Exit(0)
			} else {
				//println("ok")
			}
		}()
	}
	wg.Wait()
}
func BenchamrkSCHT(w *Workload, time int, check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewSHT(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht
}
func BenchamrkCMHT(w *Workload, time int, check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewCMHT(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht

}
func BenchamrkConcMHT(w *Workload, time int, check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewConcMHT(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht

}
func BenchamrkACHT(w *Workload, time int, dis,check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewAHT(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht

}
func BenchamrkLCHT(w *Workload, time int, check bool) *BaseHashTable{
	w.Reset()
	var ht BaseHashTable
	ht = NewLHT(w.Length)
	benchamrk(&ht,w,time,check)
	return &ht

}

func putLoad(w *Workload, ht *BaseHashTable) {
	es := NewStore()
	for begin,end := w.Read(); begin < w.Length; begin,end = w.Read(){
		for i:= begin; i < end ; i++ {
			tmp :=es.GetAddr()
			tmp.KV = w.KV[i]
			tmp.next = nil
			(*ht).ConcurrentPut(tmp)
		}
	}
	//println("succeed input")
}
func UnsafeCheck(w *Workload, ht *HashTable) bool {
	//return true
	for i := uint64(0); i < w.Length; i++ {
		if c:=ht.UnsafeCount(&(*w).KV[i]) ; c!= w.times[w.KV[i]] {
			println("ERROR: time error ", (*w).KV[i].key, " actural time = ",c," time = ", w.times[w.KV[i]])
			return false
		}
	}
	return true
}
func ConcCheck(w *Workload,ht *BaseHashTable) bool {
	//return  true
		for begin,end := w.Read(); begin < w.Length; begin,end = w.Read(){
		for i:= begin; i < end ; i++ {
			if c:=(*ht).Count(&w.KV[i]) ; c!= w.times[w.KV[i]] {
				println("ERROR: time error ", w.KV[i].key, " actural time = ",c," time = ", w.times[w.KV[i]])
				return false
			}
		}
	}
		return true
}

func Check(w *Workload,ht *BaseHashTable) bool {
	for i := uint64(0); i < w.Length; i++ {
		if c:=(*ht).Count(&w.KV[i]) ;c!= w.times[w.KV[i]] {
			println("ERROR: time error ", (*w).KV[i].key, " actural time = ",c," time = ", w.times[w.KV[i]])
			return false
		}
	}
	return true
}
