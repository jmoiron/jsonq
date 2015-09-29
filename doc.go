/*
Package jsonq simplify your json usage with a simple hierarchical query.

Given some json data like:

	{
		"foo": 1,
		"bar": 2,
		"test": "Hello, world!",
		"baz": 123.1,
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
	}

Decode it into a map[string]interrface{}:

	import (
		"strings"
		"encoding/json"
		"github.com/jmoiron/jsonq"
	)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

From here, you can query along different keys and indexes:

	// data["foo"] -> 1
	jq.Int("foo")

	// data["subobj"]["subarray"][1] -> 2
	jq.Int("subobj", "subarray", "1")

	// data["subobj"]["subarray"]["array"][0] -> "hello"
	jq.String("subobj", "subsubobj", "array", "0")

	// data["subobj"] -> map[string]interface{}{"subobj": ...}
	obj, err := jq.Object("subobj")

Also you can query using a *path* specification as show in the examples below:

	```go
	// data["subobj"]["subarray"][1] -> 2
	jq.Int("subobj.subarray[1]")

	// data["subobj"]["subarray"]["array"][0] -> "hello"
	jq.String("subobj.subsubobj.array[0]")
	```

Finally, As functions have been included so that if you are *sure* the call will succeed you can inline the values. If these calls encounter an error they will panic:

	```go
	// data["subobj"]["subarray"][1] -> 2
	fmt.Printf("%d\n", jq.AsInt("subobj", "subarray", "1"))
	fmt.Printf("%d\n", jq.AsInt("subobj.subarray[1]"))

	// data["subobj"]["subarray"]["array"][0] -> "hello"
	fmt.Printf("%s\n", jq.AsString("subobj", "subsubobj", "array", "0"))
	fmt.Printf("%s\n", jq.AsString("subobj.subsubobj.array[0]"))
	```


	Notes:

Missing keys, out of bounds indexes, and type failures will return errors.
For simplicity, integer keys (ie, {"0": "zero"}) are inaccessible by `jsonq`
as integer strings are assumed to be array indexes.

*/
package jsonq
