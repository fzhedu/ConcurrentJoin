package hashtable

import (
	"math/rand"
	"os"
	"sync"
)

func GenLoad(num uint64) (kv []KVpair) {
	for i := uint64(0); i < num; i++ {
		kv = append(kv, KVpair{i+11, rand.Int63()})
	}
	return
}

func BenchamrkSTHT(kvLoad *[]KVpair, len uint64, time int) {
	ht := NewHt( len * uint64(time))

	for t := 0; t < time; t++ {
		for i := uint64(0); i < len; i++ {
			ht.UnsafePut(getHashValue((*kvLoad)[i].key), &(*kvLoad)[i])
		}
	}
	ht.UnsafePrint()
	ok := UnsafeCheck(kvLoad, len, ht, time)
	if !ok {
		println("ERROR occor")
	} else {
		println("ok")
	}
}
func BenchamrkCHT(kvLoad *[]KVpair, len uint64, time int) {
	var ht BaseHashTable
	ht = NewHt( len * uint64(time))
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(kvLoad, len, &ht)
		}()
	}
	wg.Wait()
	ht.Print()
	ok := Check(kvLoad, len, &ht, time)
	if !ok {
		println("ERROR occor")
	} else {
		println("ok")
	}
}
func BenchamrkSCHT(kvLoad *[]KVpair, len uint64, time int) {
	var ht BaseHashTable
	ht = NewSHT(len * uint64(time))
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(kvLoad, len, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
	//	ht.Print()
		ok := Check(kvLoad, len, &ht, time)
		if !ok {
			println("ERROR occor")
			os.Exit(1)
		} else {
			//os.Exit(0)
			println("ok")
		}
}

func BenchamrkACHT(kvLoad *[]KVpair, len uint64, time int) {
	var ht BaseHashTable
	ht = NewAHT(len * uint64(time))
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(kvLoad, len, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
	ht.Print()
	ok := Check(kvLoad, len, &ht, time)
	if !ok {
		println("ERROR occor")
		os.Exit(1)
	} else {
		//os.Exit(0)
		println("ok")
	}
}
func BenchamrkLCHT(kvLoad *[]KVpair, len uint64, time int) {
	var ht BaseHashTable
	ht = NewLHT(len * uint64(time))
	wg := &sync.WaitGroup{}
	wg.Add(time)

	for t := 0; t < time; t++ {
		go func() {
			defer wg.Done()
			putLoad(kvLoad, len, &ht)
		}()
	}
	wg.Wait()
//	println("all thread Done")
/*	ht.Print()
	ok := Check(kvLoad, len, &ht, time)
	if !ok {
		println("ERROR occor")
	} else {
		println("ok")
	}*/
}

func putLoad(kvLoad *[]KVpair, len uint64, ht *BaseHashTable) {
	for i := uint64(0); i < len; i++ {
		(*ht).ConcurrentPut(getHashValue((*kvLoad)[i].key), &(*kvLoad)[i])
	}
	//println("succeed input")
}
func UnsafeCheck(kv *[]KVpair, len uint64, ht *HashTable, time int) bool {
	for i := uint64(0); i < len; i++ {
		if c:=ht.UnsafeCount(&(*kv)[i]) ; c!= time {
			println("ERROR: time error ", (*kv)[i].key, " actural time = ",c," time = ", time)
			return false
		}
	}
	return true
}
func Check(kv *[]KVpair, len uint64, ht *BaseHashTable, time int) bool {
	for i := uint64(0); i < len; i++ {
		if c:=(*ht).Count(&(*kv)[i]) ; c!= time {
			println("ERROR: time error ", (*kv)[i].key, " actural time = ",c," time = ", time)
			return false
		}
	}
	return true
}
