package asndb

import (
	"testing"

	"github.com/JustinAzoff/asnlookup/asndb"
)

func doTestLookup(t *testing.T, b *asndb.AsnBackend) {
	rec, err := b.Lookup("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Owner != "GOOGLE - Google Inc." {
		t.Fatalf("Expected 'GOOGLE - Google Inc.', got %q", rec.Owner)
	}

}

func TestBackend(t *testing.T) {
	b, err := asndb.NewAsnBackend("../asn.db", "../asnames.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("test=Lookup", func(t *testing.T) {
		doTestLookup(t, b)
	})
}

var result *AsnBackend

func BenchmarkNewBackend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		be, err := NewAsnBackend("../asn.db", "../asnames.json")
		if err != nil {
			b.Fatal(err)
		}
		result = be
	}

}

var record Record

func BenchmarkLookup(b *testing.B) {
	be, err := NewAsnBackend("../asn.db", "../asnames.json")
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
