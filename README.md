# jsonq

Simplify your golang json usage by extracting fields or items from arrays and
objects with a simple, hierarchical query.

# installing

```
go get github.com/jmoiron/jsonq
```

# usage

Given some json data like:

```javascript
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
```

Decode it into a `map[string]interrface{}`:

```go
import (
	"strings"
	"encoding/json"
	"github.com/jmoiron/jsonq"
)

data := map[string]interface{}{}
dec := json.NewDecoder(strings.NewReader(jsonstring))
dec.Decode(&data)
jq := jsonq.NewQuery(data)
```

From here, you can query along different keys and indexes:

```go
// data["foo"] -> 1
jq.Int("foo")

// data["subobj"]["subarray"][1] -> 2
jq.Int("subobj", "subarray", "1")

// data["subobj"]["subarray"]["array"][0] -> "hello"
jq.String("subobj", "subsubobj", "array", "0")
```

Missing keys, out of bounds indexes, and type failures are all returned in the
error channels. For simplicity, integer keys (ie, {"0": "zero"}) are inaccessible
by `jsonq` as integer strings are assumed to be array indexes.

Suggestions/comments please tweet [@jmoiron](http://twitter.com/jmoiron)

