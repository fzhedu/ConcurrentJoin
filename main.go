package main

import (
	"github.com/ParallelBuild/hashtable"
)

var num = uint64(4000)
var concurrency = 10

func main() {

	for
	{
		var w  hashtable.Workload
		w.GenLoad(num)
		//hashtable.BenchamrkUnsafeHT(&w,concurrency)
		w.Reset()
		hashtable.BenchamrkCHT(&w,concurrency)
/*		w.Reset()
		hashtable.BenchamrkLCHT(&w,concurrency)
		w.Reset()
		hashtable.BenchamrkSCHT(&w,concurrency)
		w.Reset()
		hashtable.BenchamrkACHT(&w,concurrency)*/
	}

}
