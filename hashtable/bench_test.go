package hashtable

import (
	"testing"
)

var num = uint64(4000000)
var concurrency = 4
var w = NewW(num)


func BenchmarkUnsafeHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkUnsafeHT(w,concurrency)
}
func BenchmarkUnsafeHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkUnsafeHT(w,0)
	}
}
func BenchmarkCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkCHT(w,concurrency)
	}
}
func BenchmarkLCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkLCHT(w,concurrency)
	}
}
func BenchmarkSCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkSCHT(w,concurrency)
	}
}
func BenchmarkACHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkACHT(w,concurrency,false)
	}
}
/*func BenchmarkACHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkACHT(w,concurrency,true)
}*/