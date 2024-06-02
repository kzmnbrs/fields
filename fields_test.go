package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type X struct {
	Y int
	Z *int
}

type I struct {
	J X
	K *X
	L string
	M uint
}

func TestList(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		exp := []string{"Y", "Z"}
		expType := []string{"int", "*int"}

		fs := List(X{})
		fsptr := List(&X{})

		assert.Equal(t, exp, pathsToKeys(fs))
		assert.Equal(t, exp, pathsToKeys(fsptr))
		assert.Equal(t, expType, pathsToTypes(fs))
		assert.Equal(t, expType, pathsToTypes(fsptr))

	})
	t.Run("nesting", func(t *testing.T) {
		exp := []string{
			"J",
			"J.Y",
			"J.Z",
			"K",
			"K.Y",
			"K.Z",
			"L",
			"M",
		}
		expType := []string{
			"fields.X",
			"int",
			"*int",
			"*fields.X",
			"int",
			"*int",
			"string",
			"uint",
		}

		fs := List(I{})
		fsptr := List(&I{})

		assert.Equal(t, exp, pathsToKeys(fs))
		assert.Equal(t, exp, pathsToKeys(fsptr))
		assert.Equal(t, expType, pathsToTypes(fs))
		assert.Equal(t, expType, pathsToTypes(fsptr))
	})
	t.Run("anon", func(t *testing.T) {
		v := struct {
			O int
			P X
			R struct {
				S int
				T complex128
			}
		}{}

		exp := []string{
			"O",
			"P",
			"P.Y",
			"P.Z",
			"R",
			"R.S",
			"R.T",
		}
		expType := []string{
			"int",
			"fields.X",
			"int",
			"*int",
			"struct { S int; T complex128 }",
			"int",
			"complex128",
		}

		fs := List(v)
		fsptr := List(&v)

		assert.Equal(t, exp, pathsToKeys(fs))
		assert.Equal(t, exp, pathsToKeys(fsptr))
		assert.Equal(t, expType, pathsToTypes(fs))
		assert.Equal(t, expType, pathsToTypes(fsptr))
	})
	t.Run("no-struct", func(t *testing.T) {
		assert.Panics(t, func() {
			List(1)
		})
	})
}

func pathsToKeys(pp []Path) []string {
	ret := make([]string, 0, len(pp))
	for _, p := range pp {
		ret = append(ret, p.String())
	}
	return ret
}

func pathsToTypes(pp []Path) []string {
	ret := make([]string, 0, len(pp))
	for _, p := range pp {
		ret = append(ret, p.Base().Type.String())
	}
	return ret
}
