package jsonq

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const TestData = `{
	"foo": 1,
	"bar": 2,
	"test": "Hello, world!",
	"baz": 123.1,
	"numstring": "42",
	"floatstring": "42.1",
	"array": [
		{"foo": 1},
		{"bar": 2},
		{"baz": 3}
	],
	"subobj": {
		"foo": 1,
		"subarray": [1,2,3],
		"subsubobj": {
			"bar": 2,
			"baz": 3,
			"array": ["hello", "world"]
		}
	},
	"bool": true
}`

func tErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
}

func TestQuery(t *testing.T) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(TestData))
	err := dec.Decode(&data)
	tErr(t, err)
	q := NewQuery(data)

	ival, err := q.Int("foo")
	if ival != 1 {
		t.Errorf("Expecting 1, got %v\n", ival)
	}
	tErr(t, err)
	ival, err = q.Int("bar")
	if ival != 2 {
		t.Errorf("Expecting 2, got %v\n", ival)
	}
	tErr(t, err)

	ival, err = q.Int("subobj", "foo")
	if ival != 1 {
		t.Errorf("Expecting 1, got %v\n", ival)
	}
	tErr(t, err)

	// test that strings can get int-ed
	ival, err = q.Int("numstring")
	if ival != 42 {
		t.Errorf("Expecting 42, got %v\n", ival)
	}
	tErr(t, err)

	for i := 0; i < 3; i++ {
		ival, err := q.Int("subobj", "subarray", fmt.Sprintf("%d", i))
		if ival != i+1 {
			t.Errorf("Expecting %d, got %v\n", i+1, ival)
		}
		tErr(t, err)
	}

	fval, err := q.Float("baz")
	if fval != 123.1 {
		t.Errorf("Expecting 123.1, got %f\n", fval)
	}
	tErr(t, err)

	// test that strings can get float-ed
	fval, err = q.Float("floatstring")
	if fval != 42.1 {
		t.Errorf("Expecting 42.1, got %v\n", fval)
	}
	tErr(t, err)

	sval, err := q.String("test")
	if sval != "Hello, world!" {
		t.Errorf("Expecting \"Hello, World!\", got \"%v\"\n", sval)
	}

	sval, err = q.String("subobj", "subsubobj", "array", "0")
	if sval != "hello" {
		t.Errorf("Expecting \"hello\", got \"%s\"\n", sval)
	}
	tErr(t, err)

	bval, err := q.Bool("bool")
	if !bval {
		t.Errorf("Expecting true, got %v\n", bval)
	}
	tErr(t, err)

	obj, err := q.Object("subobj", "subsubobj")
	tErr(t, err)
	q2 := NewQuery(obj)
	sval, err = q2.String("array", "1")
	if sval != "world" {
		t.Errorf("Expecting \"world\", got \"%s\"\n", sval)
	}
	tErr(t, err)

	aobj, err := q.Array("subobj", "subarray")
	tErr(t, err)
	if aobj[0].(float64) != 1 {
		t.Errorf("Expecting 1, got %v\n", aobj[0])
	}
}
