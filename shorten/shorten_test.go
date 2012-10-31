package shorten

import (
	"fmt"
	"testing"
)

func TestShorten(t *testing.T) {
	shortid := Shorten()
	fmt.Printf("shortid='%s'\n", shortid)
	if len(shortid) != 6 {
		t.FailNow()
	}
}

func BenchmarkShorten(b *testing.B) {
        for i := 0; i < b.N; i++ {
                shortid := Shorten()
                if len(shortid) != 6 {
                        b.FailNow()
                }
        }
}
