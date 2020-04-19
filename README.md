We set even data distribution, where the length of buckets is not too long, because we usually choose a good hash function to hash distribute data. The distribution are show below, each pair means the length, and the occurence timeof the leght.

From the results below, we see `CmapChangedHashTable` is the best if we use the map to record the buckets. It is slower than `ArrayHashTable` due to lock conflicts. But `ArrayHashTable` need to map a larger range to its length, this may worse the data distribution in the hash table so that to increase the probing time.

```
mac:hashtable fzh$ go test -v -bench=. 
---------------workload---------------
11 45025
5 128110
20 37
19 96
22 7
15 3358
1 6401
12 26328
10 70951
7 149149
13 14269
8 130040
18 226
21 12
3 52154
9 101380
14 7253
4 91006
17 608
6 148821
16 1408
```
SHARD_COUNT=320, there are 320 locks for `concurrentMap`
```
2 22411
goos: darwin
goarch: amd64
pkg: github.com/ParallelBuild/hashtable
BenchmarkUnsafeHTDis-12     	       1	2293166737 ns/op
BenchmarkUnsafeHT-12        	       1	1124882396 ns/op
BenchmarkCHT-12             	       1	1271832027 ns/op
BenchmarkLCHT-12            	       1	1219375592 ns/op
BenchmarkSCHT-12            	       1	2275332512 ns/op
BenchmarkCMHT-12            	       2	 822315852 ns/op
BenchmarkConcMHT-12         	       3	 438444883 ns/op
BenchmarkACHT-12            	       9	 132368679 ns/op
BenchmarkUnsafeHTRead-12    	       1	6762562949 ns/op
BenchmarkCHTRead-12         	       1	1849035497 ns/op
BenchmarkLCHTRead-12        	       1	1820317635 ns/op
BenchmarkSCHTRead-12        	       1	2289176365 ns/op
BenchmarkCMHTRead-12        	       1	2992674152 ns/op
BenchmarkConcHTRead-12      	       1	1969843381 ns/op
BenchmarkAHTRead-12         	       1	1707496822 ns/op
PASS
ok  	github.com/ParallelBuild/hashtable	43.563s
```
SHARD_COUNT=1000
```
goos: darwin
goarch: amd64
pkg: github.com/ParallelBuild/hashtable
BenchmarkUnsafeHTDis-12     	       1	2367000802 ns/op
BenchmarkUnsafeHT-12        	       1	1082716750 ns/op
BenchmarkCHT-12             	       1	1274581951 ns/op
BenchmarkLCHT-12            	       1	1207945528 ns/op
BenchmarkSCHT-12            	       1	2266861326 ns/op
BenchmarkCMHT-12            	       2	 903871798 ns/op
BenchmarkConcMHT-12         	       3	 429208573 ns/op
BenchmarkACHT-12            	      10	 137949060 ns/op
BenchmarkUnsafeHTRead-12    	       1	6796573730 ns/op
BenchmarkCHTRead-12         	       1	1857026478 ns/op
BenchmarkLCHTRead-12        	       1	1838800899 ns/op
BenchmarkSCHTRead-12        	       1	2398534357 ns/op
BenchmarkCMHTRead-12        	       1	3011738603 ns/op
BenchmarkConcHTRead-12      	       1	1993160606 ns/op
BenchmarkAHTRead-12         	       1	1732397450 ns/op
PASS
ok  	github.com/ParallelBuild/hashtable	44.376s
```
