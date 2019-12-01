package datetime

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	t1 := time.Date(2009, time.November, 10, 23, 5, 16, 45000000, time.UTC)
	t2, err := Parse("2009-11-10 23:05:16.045", "YYYY-MM-DD HH:mm:ss.SSS")
	if err != nil {
		t.Fatal(err)
	}
	if !t1.Equal(t2) {
		t.Fatal("not equal")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("2009-11-10 23:05:16.045", "YYYY-MM-DD HH:mm:ss.SSS")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := time.Parse("2006-01-02 15:04:05.000", "2009-11-10 23:05:16.045")
		if err != nil {
			b.Fatal(err)
		}
	}
}
