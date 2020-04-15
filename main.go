package main

import (
	"github.com/ParallelBuild/hashtable"
)

var num = uint64(20)

func main()  {
	concurrency :=(7)
	//for
	{
		kvLoad:=hashtable.GenLoad(num)
		hashtable.BenchamrkSTHT(&kvLoad, num,concurrency)
		hashtable.BenchamrkCHT(&kvLoad, num,concurrency)
	}

}
