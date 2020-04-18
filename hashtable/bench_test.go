package hashtable

import (
	"testing"
)

var num = uint64(4000000)
var concurrency = 4
var Gstep = uint64(100000)
var Gmod = uint64(10000000)
var w = NewW(num,Gstep,Gmod)


func BenchmarkUnsafeHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkUnsafeHT(w,concurrency,false)
}
func BenchmarkUnsafeHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkUnsafeHT(w,0,false)
	}
}
func BenchmarkCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkCHT(w,concurrency,false)
	}
}
func BenchmarkLCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkLCHT(w,concurrency,false)
	}
}
func BenchmarkSCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkSCHT(w,concurrency,false)
	}
}
func BenchmarkCMHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkCMHT(w,concurrency,false)
	}
}
func BenchmarkConcMHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkConcMHT(w,concurrency,false)
	}
}
func BenchmarkACHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkACHT(w,concurrency,false,false)
	}
}
/*func BenchmarkACHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkACHT(w,concurrency,true)
}*/