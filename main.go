package main

import (
	"fmt"
	"github.com/ParallelBuild/hashtable"
	"sync"
	"sync/atomic"
	"unsafe"
)

var num = uint64(400)
var concurrency = (20)


type e struct{ val int }

func test() {
	var m sync.Map

	// make a map head
	var head unsafe.Pointer
	// point head to a val
	val := &e{val: 1}
	atomic.StorePointer(&head, unsafe.Pointer(val))
	// store head into map
	h, ok := m.LoadOrStore("1", &head)
	fmt.Println(ok)
	// load and check
	p := atomic.LoadPointer(h.(*unsafe.Pointer))
	val2 := (*e)(p)
	fmt.Println(val2)

	// load and cas
	h, ok = m.Load("1")
	fmt.Println(ok)
	new := unsafe.Pointer(&e{val: 2})
	for {
		if atomic.CompareAndSwapPointer(h.(*unsafe.Pointer), p, new) {
			break
		}
	}
	// check updated value
	h, ok = m.Load("1")
	fmt.Println(ok)
	fmt.Println((*e)(atomic.LoadPointer(h.(*unsafe.Pointer))))
}

func main() {

	for
	{
		kvLoad:=hashtable.GenLoad(num)
		//hashtable.BenchamrkSTHT(&kvLoad, num,concurrency)
		//hashtable.BenchamrkCHT(&kvLoad, num,concurrency)
		//hashtable.BenchamrkCHT(&kvLoad, num,concurrency)
		hashtable.BenchamrkSCHT(&kvLoad, num,concurrency)
		//hashtable.BenchamrkACHT(&kvLoad, num,concurrency)

	}

}
