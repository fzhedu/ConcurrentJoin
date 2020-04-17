package hashtable

import (
	"testing"
)

var num = uint64(4000)
var concurrency = 10
var w = NewW(num)

func BenchmarkUnsafeHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkUnsafeHT(w,concurrency)
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
		BenchamrkACHT(w,concurrency)
	}
}
