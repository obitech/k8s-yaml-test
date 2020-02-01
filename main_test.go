package main

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var cases = []struct {
	name     string
	data     string
	expected []*schema.GroupVersionKind
}{
	{
		name: "basic",
		data: `apiVersion: foo
kind: bar
---
apiVersion: test
kind: test
string: |
    ---
---
apiVersion: yup
kind: aha
`,
		expected: []*schema.GroupVersionKind{
			&schema.GroupVersionKind{Group: "", Version: "foo", Kind: "bar"},
			&schema.GroupVersionKind{Group: "", Version: "test", Kind: "test"},
			&schema.GroupVersionKind{Group: "", Version: "yup", Kind: "aha"},
		},
	},
	{
		name: "delimiter in string",
		data: `apiVersion: foo
kind: bar
---
apiVersion: test
kind: test
string: "foo
---
bar"
---
apiVersion: yup
kind: aha`,
		expected: []*schema.GroupVersionKind{
			&schema.GroupVersionKind{Group: "", Version: "foo", Kind: "bar"},
			&schema.GroupVersionKind{Group: "", Version: "test", Kind: "test"},
			&schema.GroupVersionKind{Group: "", Version: "yup", Kind: "aha"},
		},
	},
}

func TestReaders(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("YAMLReader", func(t *testing.T) {
				actual, err := UseYAMLReader(c.data)
				if err != nil {
					t.Fatalf("unexpected error with %q: %v", c.name, err)
				}
				if !reflect.DeepEqual(actual, c.expected) {
					t.Fatalf("expected: %#v\ngot: %#v", c.expected, actual)
				}
			})

			t.Run("YAMLv3Decoder", func(t *testing.T) {
				actual, err := UseYAMLv3Decoder(c.data)
				if err != nil {
					t.Fatalf("unexpected error with %q: %v", c.name, err)
				}
				if !reflect.DeepEqual(actual, c.expected) {
					t.Fatalf("expected: %#v\ngot: %#v", c.expected, actual)
				}
			})
		})
	}
}

var result []*schema.GroupVersionKind

func BenchmarkYAMLReader(b *testing.B) {
	var r []*schema.GroupVersionKind
	for n := 0; n < b.N; n++ {
		r, _ = UseYAMLReader(cases[0].data)
	}
	result = r
}

func BenchmarkYAMLv3Decoder(b *testing.B) {
	var r []*schema.GroupVersionKind
	for n := 0; n < b.N; n++ {
		r, _ = UseYAMLv3Decoder(cases[0].data)
	}
	result = r
}
