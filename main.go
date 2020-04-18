package main

import (
	"github.com/ParallelBuild/hashtable"
)

var num = uint64(4000)
var concurrency = 4
var Gstep = uint64(2)
var Gmod = uint64(1000)

func main() {
	var w  hashtable.Workload
	w.GenLoad(num,Gstep,Gmod,true)
	w.PrintDis()

	for
	{
		hashtable.BenchamrkConcMHT(&w,concurrency,true)
		//hashtable.BenchamrkCMHT(&w,concurrency,true)
		//hashtable.BenchamrkUnsafeHT(&w,concurrency,false)
		//hashtable.BenchamrkCHT(&w,concurrency,false)
		//hashtable.BenchamrkLCHT(&w,concurrency,false)
		//hashtable.BenchamrkSCHT(&w,concurrency,false)
		//hashtable.BenchamrkACHT(&w,concurrency,false)
	}

}
