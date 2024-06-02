package fields

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Field struct{ reflect.StructField }
	Path  []Field
)

// List collects the Walk results.
func List(x any) []Path {
	var ret []Path
	Walk(x)(func(path Path) bool {
		ret = append(ret, path)
		return true
	})
	return ret
}

// Walk traverses the struct fields in descending order.
// Panics if x is neither a struct nor a struct pointer.
func Walk(x any) func(yield func(path Path) bool) {
	var t reflect.Type
	switch x := x.(type) {
	case reflect.Type:
		t = x
	case reflect.Value:
		t = x.Type()
	default:
		t = reflect.TypeOf(x)
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		msg := fmt.Sprintf("want struct, have %T", x)
		panic(msg)
	}

	return func(yield func(path Path) bool) {
		walk(nil, t)(func(dir []Field, base Field) bool {
			return yield(append(dir, base))
		})
	}
}

func walk(dir []Field, t reflect.Type) func(yield func(dir []Field, base Field) bool) {
	return func(yield func([]Field, Field) bool) {
		n := t.NumField()
		for i := 0; i < n; i++ {
			base := t.Field(i)
			if !yield(dir, Field{base}) {
				return
			}

			t := base.Type
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}

			if t.Kind() == reflect.Struct {
				walk(append(dir, Field{base}), t)(func(dir []Field, base Field) bool {
					return yield(dir, base)
				})
			}
		}
	}
}

// Zero returns the type zero value.
func (f Field) Zero() any {
	return reflect.New(f.Type).Elem()
}

// Base returns the last path component, assuming that the path is not nil.
func (p Path) Base() Field {
	return p[len(p)-1]
}

// String return the dot-joined field path.
func (p Path) String() string {
	if len(p) == 0 {
		return ""
	}

	w := strings.Builder{}
	for i := range p {
		w.WriteString(p[i].Name)
		if i != len(p)-1 {
			w.WriteByte('.')
		}
	}
	return w.String()
}
