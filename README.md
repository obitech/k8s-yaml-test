# Evaluation Multi-Doc Parsing

Benchmark:

```
$ go test -benchmem -bench=.
goos: darwin
goarch: amd64
pkg: github.com/obitech/yaml_test
BenchmarkYAMLReader-4              21572             54142 ns/op           28852 B/op        272 allocs/op
BenchmarkYAMLv3Decoder-4           36081             31269 ns/op           15496 B/op        151 allocs/op
PASS
ok      github.com/obitech/yaml_test    3.281s
```
