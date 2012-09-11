// Copyright 2012 Jason Moiron. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

/*
Simplify your golang json usage by extracting fields or items from arrays and
objects with a simple, hierarchical query.

Installing

	go get github.com/jmoiron/jsonq

Usage

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

Missing keys, out of bounds indexes, and type failures will return errors.
For simplicity, integer keys (ie, {"0": "zero"}) are inaccessible by `jsonq` 
as integer strings are assumed to be array indexes.
*/

package jsonq

import (
	"fmt"
	"strconv"
)

type JsonQuery struct {
	blob map[string]interface{}
}

// Create a new JsonQuery obj from a json-decoded interface{}
func NewQuery(data interface{}) *JsonQuery {
	j := new(JsonQuery)
	j.blob = data.(map[string]interface{})
	return j
}

// Extract a Bool from some json
func (j *JsonQuery) Bool(s ...string) (bool, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return false, err
	}
	switch val.(type) {
	case bool:
		return val.(bool), nil
	}
	return false, fmt.Errorf("Expected boolean value for Bool, got \"%v\"\n", val)
}

// Extract a float from some json
func (j *JsonQuery) Float(s ...string) (float64, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return 0.0, err
	}
	switch val.(type) {
	case float64:
		return val.(float64), nil
	}
	return 0.0, fmt.Errorf("Expected numeric value for Float, got \"%v\"\n", val)
}

// Extract an int from some json
func (j *JsonQuery) Int(s ...string) (int, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return 0, err
	}
	switch val.(type) {
	case float64:
		return int(val.(float64)), nil
	}
	return 0, fmt.Errorf("Expected numeric value for Int, got \"%v\"\n", val)
}

// Extract a string from some json
func (j *JsonQuery) String(s ...string) (string, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return "", err
	}
	switch val.(type) {
	case string:
		return val.(string), nil
	}
	return "", fmt.Errorf("Expected string value for String, got \"%v\"\n", val)
}

// Extract an object from some json
func (j *JsonQuery) Object(s ...string) (map[string]interface{}, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return map[string]interface{}{}, err
	}
	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{}), nil
	}
	return map[string]interface{}{}, fmt.Errorf("Expected json object for Object, get \"%v\"\n", val)
}

// Extract an array from some json
func (j *JsonQuery) Array(s ...string) ([]interface{}, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return []interface{}{}, err
	}
	switch val.(type) {
	case []interface{}:
		return val.([]interface{}), nil
	}
	return []interface{}{}, fmt.Errorf("Expected json array for Array, get \"%v\"\n", val)
}

// Recursively query a decoded json blob
func rquery(blob interface{}, s ...string) (interface{}, error) {
	var (
		val interface{}
		err error
	)
	val = blob
	for _, q := range s {
		val, err = query(val, q)
		if err != nil {
			return nil, err
		}
	}
	switch val.(type) {
	case nil:
		return nil, fmt.Errorf("Nil value found at %s\n", s[len(s)-1])
	}
	return val, nil
}

// Query a json blob for a single field or index.  If query is a string, then
// the blob is treated as a json object (map[string]interface{}).  If query is
// an integer, the blob is treated as a json array ([]interface{}).  Any kind
// of key or index error will result in a nil return value with an error set.
func query(blob interface{}, query string) (interface{}, error) {
	index, err := strconv.Atoi(query)
	// if it's an integer, then we treat the current interface as an array
	if err == nil {
		switch blob.(type) {
		case []interface{}:
		default:
			return nil, fmt.Errorf("Array index on non-array %v\n", blob)
		}
		if len(blob.([]interface{})) > index {
			return blob.([]interface{})[index], nil
		}
		return nil, fmt.Errorf("Array index %d on array %v out of bounds\n", index, blob)
	}

	// blob is likely an object, but verify first
	switch blob.(type) {
	case map[string]interface{}:
	default:
		return nil, fmt.Errorf("Object lookup \"%s\" on non-object %v\n", query, blob)
	}

	val, ok := blob.(map[string]interface{})[query]
	if !ok {
		return nil, fmt.Errorf("Object %v does not contain field %s\n", blob, query)
	}
	return val, nil
}
