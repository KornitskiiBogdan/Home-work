//go:build bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func readUsersData(b *testing.B) []byte {
	b.Helper()
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer r.Close()
	f, err := r.File[0].Open()
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	return data
}
func BenchmarkGetDomainStat(b *testing.B) {
	data := readUsersData(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetDomainStat(bytes.NewReader(data), "biz")
		if err != nil {
			b.Fatal(err)
		}
	}
}
