//go:build benchmark
// +build benchmark

package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// Использование:
//
//	Запуск до изменений
//	go test -bench=BenchmarkGetDomainStat -benchmem -count=10 -tags benchmark ./hw10_program_optimization > before.txt
//	Запуск после изменений
//	go test -bench=BenchmarkGetDomainStat -benchmem -count=10 -tags benchmark ./hw10_program_optimization > after.txt
//	benchstat before.txt after.txt
func BenchmarkGetDomainStat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, _ := r.File[0].Open()
		GetDomainStat(data, "biz")
		data.Close()
		r.Close()
	}
}
