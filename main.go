package main

import (
	"github.com/ParallelBuild/hashtable"
)

var num = uint64(100)
var concurrency = 4
var Gstep = uint64(2)
var Gmod = uint64(1000)

func main() {
	var w  hashtable.Workload
	w.GenLoad(num,Gstep,Gmod,true)
	w.PrintDis()
	//for
	{

		//hashtable.BenchamrkConcMHT(&w,concurrency,true)
		//hashtable.BenchamrkCMHT(&w,concurrency,true)
		hashtable.BenchamrkUnsafeHT(&w,0,true,nil)
		hashtable.BenchamrkCHT(&w,concurrency,true,nil)
		//hashtable.BenchamrkLCHT(&w,concurrency,true)
		//hashtable.BenchamrkSCHT(&w,concurrency,true)
		//hashtable.BenchamrkACHT(&w,concurrency,false,true)
	}

}
