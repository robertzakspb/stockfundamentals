package taghelpers

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_GetEntityTagValues_MissingValue(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}

	_, err := GetEntityTagValues[Foo]("test")
	test.AssertError(t, err)
}

func Test_GetEntityTagValues(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}

	tags, err := GetEntityTagValues[Foo]("sql")
	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(tags))
	test.AssertEqual(t, "a", tags[0])
	test.AssertEqual(t, "b", tags[1])
}
