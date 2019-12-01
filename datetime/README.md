# Go Datetime

I've adapted the datetime format style of js and python. The style is so much different 
in go then I'd like to implement a moment style formatter and parser.

### Format
You can format like following (learn from moment)

``` go
r := Format(time.Now(), "YYYY-MM-DD HH:mm:ss.SSS")
```

### Parse
Moment style parse

``` go
t, err := Parse("2009-11-10 23:05:16.045", "YYYY-MM-DD HH:mm:ss.SSS")
```

### Benchmark
The benchmark result compared with standard library

```
$ go test -bench=. -benchmem
BenchmarkFormat-4        3000000               448 ns/op              32 B/op          1 allocs/op
BenchmarkStdFormat-4     5000000               313 ns/op              32 B/op          1 allocs/op
BenchmarkParse-4        10000000               247 ns/op               0 B/op          0 allocs/op
BenchmarkStdParse-4      5000000               332 ns/op               0 B/op          0 allocs/op
```
