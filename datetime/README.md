# Go Datetime

I've adapted the datetime format style of js and python. The style is so much different 
in go then I'd like to implement a moment style formatter and parser.

### Format
You can format like following (learn from moment)

``` go
r := Format(time.Now(), "YYYY-MM-DD HH:mm:ss.SSS")
```

### Benchmark
The benchmark result compared with standard library

```
$ go test -bench=. -benchmem
BenchmarkFormat-4        3000000               546 ns/op              32 B/op          1 allocs/op
BenchmarkStdFormat-4     3000000               465 ns/op              32 B/op          1 allocs/op
```

### TODO
- support parser
- performance
