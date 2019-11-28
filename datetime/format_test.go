package datetime

import (
	"fmt"
	"testing"
	"time"
)

func ExampleFormat() {
	t := time.Date(2009, time.November, 10, 23, 5, 16, 45000000, time.UTC)
	r := Format(t, "lll YYYY-MM-DD NNN HH:mm:ss.SSS ===")
	fmt.Println(r)
	// Output:
	// lll 2009-11-10 NNN 23:05:16.045 ===
}

var result string

func BenchmarkFormat(b *testing.B) {
	t := time.Date(2009, time.November, 10, 23, 5, 16, 45000000, time.UTC)
	var r string
	for i := 0; i < b.N; i++ {
		r = Format(t, "YYYY-MM-DD HH:mm:ss.SSS")
	}
	result = r
}

func BenchmarkStdFormat(b *testing.B) {
	t := time.Date(2009, time.November, 10, 23, 5, 16, 45000000, time.UTC)
	var r string
	for i := 0; i < b.N; i++ {
		r = t.Format("2006-01-02 15:04:05.000")
	}
	result = r
}
