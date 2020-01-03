package hostdb

import (
	"testing"
)

func doTestLookup(t *testing.T, b *HostBackend) {
	rec, err := b.Lookup("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Host != "dns" {
		t.Fatalf("Expected 'dns', got %q", rec.Host)
	}

}

func TestBackend(t *testing.T) {
	b, err := NewHostBackend("shrunken.csv.gz")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("test=Lookup", func(t *testing.T) {
		doTestLookup(t, b)
	})
}

var result *HostBackend

func BenchmarkNewBackend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		be, err := NewHostBackend("shrunken.csv.gz")
		if err != nil {
			b.Fatal(err)
		}
		result = be
	}

}

var record Record

func BenchmarkLookup(b *testing.B) {
	be, err := NewHostBackend("shrunken.csv.gz")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec, err := be.Lookup("8.8.8.8")
		if err != nil {
			b.Fatal(err)
		}
		record = rec
	}

}
