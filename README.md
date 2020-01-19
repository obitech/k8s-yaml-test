# Evaluating Multi-Doc Parsing in Kubernetes

Evaluating options to perform multi doc parsing in Kubernetes.

Currently there seem to be two ways on how to parse multi-doc YAML files in
Kubernetes: via [YAMLReader][yamlreader] and [FrameReader][framereader]. The
first is used by `kubeadm` whereas  the latter seems to be for parsing YAML
files from streams, as used by `kubectl`.

A third option would be to use [gopkg.in/yaml.v3][yamlv3] which seems to be
both cleaner and faster, but afaik it's not vendored into k/k.

Benchmark:

```
$ go test -benchmem -bench=.
goos: darwin
goarch: amd64
pkg: github.com/obitech/k8s-yaml-test
BenchmarkYAMLReader-4              18910             57976 ns/op           28852 B/op        272 allocs/op
BenchmarkYAMLv3Decoder-4           36643             36277 ns/op           15496 B/op        151 allocs/op
PASS
ok      github.com/obitech/k8s-yaml-test        3.776s
```

[yamlreader]: github.com/kubernetes/kubernetes/blob/1e40f93d34802f8a41cb916446f660e226c832ee/staging/src/k8s.io/apimachinery/pkg/util/yaml/decoder.go#L256
[framereader]: https://github.com/kubernetes/kubernetes/blob/470dfbfc4848cee4897a1b176d20611668820492/staging/src/k8s.io/apimachinery/pkg/runtime/serializer/json/json.go#L372
[yamlv3]: https://godoc.org/gopkg.in/yaml.v3
