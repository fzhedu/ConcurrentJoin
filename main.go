package main

import (
	"github.com/ParallelBuild/hashtable"
)

var num = uint64(10000)
var concurrency = 10

func main() {
	var w  hashtable.Workload
	w.GenLoad(num)
	w.PrintDis()

	//for
	{
		hashtable.BenchamrkUnsafeHT(&w,concurrency)
		hashtable.BenchamrkCHT(&w,concurrency)
		//hashtable.BenchamrkLCHT(&w,concurrency)
		//hashtable.BenchamrkSCHT(&w,concurrency)
		//hashtable.BenchamrkACHT(&w,concurrency)
	}

}
