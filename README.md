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
