package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	yamlv3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	sigsyaml "sigs.k8s.io/yaml"
)

var input = `apiVersion: foo
kind: bar
---
apiVersion: test
kind: test
string: |
   ---
---
apiVersion: yup
kind: aha
`

func gvkFromTypeMeta(tm runtime.TypeMeta) (*schema.GroupVersionKind, error) {
	gv, err := schema.ParseGroupVersion(tm.APIVersion)
	if err != nil {
		return nil, err
	}
	gvk := &schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    tm.Kind,
	}
	return gvk, nil
}

// UseYAMLReader uses sigs.k8s.io/yaml.YAMLReader to unmarshal a multi-doc
// into a slice of GVK. This is the way kubeadm does it in SplitYAMLDocuments:
// https://github.com/kubernetes/kubernetes/blob/v1.17.1/cmd/kubeadm/app/util/marshal.go#L76
func UseYAMLReader(data string) ([]*schema.GroupVersionKind, error) {
	// NewYAMLReader needs a *bufio.Reader which needs an io.Reader
	fr := yaml.NewYAMLReader(bufio.NewReader(bytes.NewReader([]byte(data))))
	out := []*schema.GroupVersionKind{}
	for {
		// Read from `data` until we reach EOF
		b, err := fr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Unmarshal and add to `out`. Uses gopkg.in/yaml.v2 under the hood.
		var tm runtime.TypeMeta
		if err := sigsyaml.Unmarshal(b, &tm); err != nil {
			return nil, err
		}
		gvk, err := gvkFromTypeMeta(tm)
		if err != nil {
			return nil, err
		}
		out = append(out, gvk)
	}
	return out, nil
}

func UseYAMLv3Decoder(data string) ([]*schema.GroupVersionKind, error) {
	dec := yamlv3.NewDecoder(bytes.NewReader([]byte(data)))
	out := []*schema.GroupVersionKind{}
	for {
		var tm runtime.TypeMeta
		err := dec.Decode(&tm)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		gvk, err := gvkFromTypeMeta(tm)
		if err != nil {
			return nil, err
		}
		out = append(out, gvk)
	}

	return out, nil
}

type ClosingBuffer struct {
	*bytes.Reader
}

func (cb ClosingBuffer) Close() error {
	return nil
}
func UseFrameReader(data string) ([]*schema.GroupVersionKind, error) {
	cb := ClosingBuffer{
		bytes.NewReader([]byte(data)),
	}
	framer := json.YAMLFramer.NewFrameReader(cb)
	framer.Read([]byte{})

	// serializer := json.NewSerializerWithOptions(json.SimpleMetaFactory{}, nil, nil, json.SerializerOptions{Yaml: true, Pretty: false, Strict: options.Strict})
	// sd := streaming.NewDecoder(json.YAMLFramer.NewFrameReader(cb), serializer)
	// gvk := 	&schema.GroupVersionKind{
	// 	Group:   "",
	// 	Version: "foo",
	// 	Kind:    "bar",
	// }
	// sd.Decode(gvk, into runtime.Object)
	return nil, nil
}

func main() {
	fmt.Println("-- UseYAMLReader")
	if _, err := UseYAMLReader(input); err != nil {
		panic(err)
	}
	fmt.Println("-- UseYAMLv3Decoder")
	if _, err := UseYAMLv3Decoder(input); err != nil {
		panic(err)
	}
	fmt.Println("-- UseFrameReader")
	if _, err := UseFrameReader(input); err != nil {
		panic(err)
	}
}
