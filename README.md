We set even data distribution, where the length of buckets is not too long, because we usually choose a good hash function to hash distribute data. The distribution are show below, each pair means the length, and the occurence time of the leght.

From the results below, we see `CmapChangedHashTable` is the best if we use the map to record the buckets. It is slower than `ArrayHashTable` due to lock conflicts. But `ArrayHashTable` need to map a larger range to its length, this may worse the data distribution in the hash table so that to increase the probing time.

However, in [Impala](https://github.com/apache/impala/blob/04fd9ae268d89b07e2a692a916bf2ddcfb2e351b/be/src/exec/hash-table.h), and [Spark](https://github.com/apache/spark/blob/master/sql/core/src/main/scala/org/apache/spark/sql/execution/joins/HashedRelation.scala), they take the array to record the buckets, and adopt linear and quadratic probing to solve hash conflicts. 
JAVA [hash table](https://www.jianshu.com/p/6c95f8216950) also use the array, and use the linked list to solve hash conflicts. So here we pefer the `ArrayHashTable`, and use the linked list to solve hash conflicts, similar to JAVA hash table, But we also implement the `CmapChangedHashTable` to avoid the worst cases.

```
go test -v -bench=. -benchmem -benchtime 10s
 
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
2 22411
```
SHARD_COUNT=320, there are 320 locks for `concurrentMap`
```
BenchmarkUnsafeHT-12        	      10	1115802586 ns/op	809876892 B/op	     880 allocs/op
// inspried by go.sync.map, but just use a mutex to solve write conflicts on a map, it is slower than unsafeHT
BenchmarkCHT-12             	       8	1332759384 ns/op	873115653 B/op	   20112 allocs/op
// use a mutex and a map
BenchmarkLCHT-12            	       9	1256560841 ns/op	489313306 B/op	     923 allocs/op
// use go.sync.Map
BenchmarkSCHT-12            	       5	2264735399 ns/op	403812873 B/op	 9050212 allocs/op
// use [ConcurrentMap], multi map instances and mutli mutex
BenchmarkCMHT-12            	      13	 899988807 ns/op	358876289 B/op	 7034643 allocs/op
// inspired by [ConcurrentMap], multi map instances and mutli mutex
BenchmarkConcMHT-12         	      27	 467489548 ns/op	224376689 B/op	   29855 allocs/op
BenchmarkACHT-12            	     100	 116403336 ns/op	545231995 B/op	     923 allocs/op
BenchmarkUnsafeHTRead-12    	       2	6789193614 ns/op	       0 B/op	       0 allocs/op
BenchmarkCHTRead-12         	       6	1853222645 ns/op	     210 B/op	       1 allocs/op
BenchmarkLCHTRead-12        	       6	1808264814 ns/op	      16 B/op	       1 allocs/op
BenchmarkSCHTRead-12        	       5	2345573828 ns/op	      57 B/op	       1 allocs/op
BenchmarkCMHTRead-12        	       4	3046690736 ns/op	      16 B/op	       1 allocs/op
BenchmarkConcHTRead-12      	       6	2006695104 ns/op	      16 B/op	       1 allocs/op
BenchmarkAHTRead-12         	       6	1743778684 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/ParallelBuild/hashtable	262.071s
```
SHARD_COUNT=1000
```
BenchmarkUnsafeHT-12        	      10	1141535048 ns/op	809876892 B/op	     880 allocs/op
BenchmarkCHT-12             	       8	1281362889 ns/op	873108256 B/op	   20059 allocs/op
BenchmarkLCHT-12            	       9	1217738020 ns/op	489291600 B/op	     924 allocs/op
BenchmarkSCHT-12            	       5	2293363611 ns/op	403861398 B/op	 9050257 allocs/op
BenchmarkCMHT-12            	      13	 895834882 ns/op	358840019 B/op	 7034517 allocs/op
BenchmarkConcMHT-12         	      27	 429976364 ns/op	253890519 B/op	   37411 allocs/op
BenchmarkACHT-12            	     100	 104002024 ns/op	545259494 B/op	     923 allocs/op
BenchmarkUnsafeHTRead-12    	       2	6789266753 ns/op	       0 B/op	       0 allocs/op
BenchmarkCHTRead-12         	       6	1867191092 ns/op	     210 B/op	       1 allocs/op
BenchmarkLCHTRead-12        	       6	1801903654 ns/op	     272 B/op	       1 allocs/op
BenchmarkSCHTRead-12        	       5	2363689679 ns/op	      38 B/op	       1 allocs/op
BenchmarkCMHTRead-12        	       4	3125928058 ns/op	     208 B/op	       1 allocs/op
BenchmarkConcHTRead-12      	       5	2042129845 ns/op	      92 B/op	       1 allocs/op
BenchmarkAHTRead-12         	       6	1729392486 ns/op	     272 B/op	       1 allocs/op
```
[ConcurrentMap](https://github.com/orcaman/concurrent-map)