package hashtable

import (
	"testing"
)

var num = uint64(7000000)
var concurrency = 4
var Gstep = uint64(10000)
var Gmod = uint64(1000000)
var w = NewW(num,Gstep,Gmod)

var check =false
func BenchmarkUnsafeHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkUnsafeHT(w,concurrency,check)
}
func BenchmarkUnsafeHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkUnsafeHT(w,0,check)
	}
}

func BenchmarkCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkCHT( w,concurrency,check)
	}
}

func BenchmarkLCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkLCHT(w,concurrency,check)
	}
}
func BenchmarkSCHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkSCHT(w,concurrency,check)
	}
}
func BenchmarkCMHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkCMHT(w,concurrency,check)
	}
}
func BenchmarkConcMHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkConcMHT(w,concurrency,check)
	}
}
func BenchmarkACHT(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkACHT(w,concurrency,false,check)
	}
}
/*func BenchmarkACHTDis(b *testing.B) {
	b.ResetTimer()
	BenchamrkACHT(w,concurrency,true)
}*/
func BenchmarkUnsafeHTRead(b *testing.B) {
	ht :=BenchamrkUnsafeHT(w,0,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnsafeCheck(w,ht)
	}
}
func BenchmarkCHTRead(b *testing.B) {
	ht := BenchamrkCHT( w,concurrency,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}
func BenchmarkLCHTRead(b *testing.B) {
	ht := BenchamrkLCHT( w,concurrency,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}
func BenchmarkSCHTRead(b *testing.B) {
	ht := BenchamrkSCHT( w,concurrency,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}
func BenchmarkCMHTRead(b *testing.B) {
	ht := BenchamrkCMHT( w,concurrency,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}
func BenchmarkConcHTRead(b *testing.B) {
	ht := BenchamrkConcMHT( w,concurrency,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}
func BenchmarkAHTRead(b *testing.B) {
	ht := BenchamrkACHT( w,concurrency,true,check)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchamrkRead(ht,w,concurrency)
	}
}