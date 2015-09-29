package jsonq

import (
	"fmt"
	"strconv"
	"strings"
)

// JsonQuery is an object that enables querying of a Go map with a simple
// positional query language.
type JsonQuery struct {
	blob                    map[string]interface{}
	SingleValuePanicOnError bool
}

// stringFromInterface converts an interface{} to a string and returns an error if types don't match.
func stringFromInterface(val interface{}) (string, error) {
	switch val.(type) {
	case string:
		return val.(string), nil
	}
	return "", fmt.Errorf("Expected string value for String, got \"%v\"\n", val)
}

// boolFromInterface converts an interface{} to a bool and returns an error if types don't match.
func boolFromInterface(val interface{}) (bool, error) {
	switch val.(type) {
	case bool:
		return val.(bool), nil
	}
	return false, fmt.Errorf("Expected boolean value for Bool, got \"%v\"\n", val)
}

// floatFromInterface converts an interface{} to a float64 and returns an error if types don't match.
func floatFromInterface(val interface{}) (float64, error) {
	switch val.(type) {
	case float64:
		return val.(float64), nil
	case int:
		return float64(val.(int)), nil
	case string:
		fval, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return fval, nil
		}
	}
	return 0.0, fmt.Errorf("Expected numeric value for Float, got \"%v\"\n", val)
}

// intFromInterface converts an interface{} to an int and returns an error if types don't match.
func intFromInterface(val interface{}) (int, error) {
	switch val.(type) {
	case float64:
		return int(val.(float64)), nil
	case string:
		ival, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return int(ival), nil
		}
	case int:
		return val.(int), nil
	}
	return 0, fmt.Errorf("Expected numeric value for Int, got \"%v\"\n", val)
}

// objectFromInterface converts an interface{} to a map[string]interface{} and returns an error if types don't match.
func objectFromInterface(val interface{}) (map[string]interface{}, error) {
	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{}), nil
	}
	return map[string]interface{}{}, fmt.Errorf("Expected json object for Object, got \"%v\"\n", val)
}

// arrayFromInterface converts an interface{} to an []interface{} and returns an error if types don't match.
func arrayFromInterface(val interface{}) ([]interface{}, error) {
	switch val.(type) {
	case []interface{}:
		return val.([]interface{}), nil
	}
	return []interface{}{}, fmt.Errorf("Expected json array for Array, got \"%v\"\n", val)
}

// NewQuery creates a new JsonQuery obj from an interface{}.
func NewQuery(data interface{}) *JsonQuery {
	j := new(JsonQuery)
	j.blob = data.(map[string]interface{})
	j.SingleValuePanicOnError = true
	return j
}

// Bool extracts a bool the JsonQuery
func (j *JsonQuery) Bool(s ...string) (bool, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return false, err
	}
	return boolFromInterface(val)
}

// AsBool extracts a bool the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsBool(s ...string) bool {
	val, err := j.Bool(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
		return false
	}
	return val
}

// Float extracts a float from the JsonQuery
func (j *JsonQuery) Float(s ...string) (float64, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return 0.0, err
	}
	return floatFromInterface(val)
}

// AsFloat extracts a float from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsFloat(s ...string) float64 {
	val, err := j.Float(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
		return 0.0
	}
	return val
}

// Int extracts an int from the JsonQuery
func (j *JsonQuery) Int(s ...string) (int, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return 0, err
	}
	return intFromInterface(val)
}

// AsInt extracts an int from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsInt(s ...string) int {
	val, err := j.Int(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// String extracts a string from the JsonQuery
func (j *JsonQuery) String(s ...string) (string, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return "", err
	}
	return stringFromInterface(val)
}

// AsString extracts a string from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsString(s ...string) string {
	val, err := j.String(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// Exists return true if node can be accessed, else false
func (j *JsonQuery) Exists(s ...string) bool {
	_, err := j.Interface(s...)
	if err != nil {
		return false
	}
	return true
}

// Object extracts a json object from the JsonQuery
func (j *JsonQuery) Object(s ...string) (map[string]interface{}, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return objectFromInterface(val)
}

// AsObject extracts a json object from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsObject(s ...string) map[string]interface{} {
	val, err := j.Object(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// Array extracts a []interface{} from the JsonQuery
func (j *JsonQuery) Array(s ...string) ([]interface{}, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return []interface{}{}, err
	}
	return arrayFromInterface(val)
}

// AsArray extracts a []interface{} from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsArray(s ...string) []interface{} {
	val, err := j.Array(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// Interface extracts an interface{} from the JsonQuery
func (j *JsonQuery) Interface(s ...string) (interface{}, error) {
	val, err := rquery(j.blob, s...)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// AsInterface extracts an interface{} from the JsonQuery, but panics on error so it can be used inline
func (j *JsonQuery) AsInterface(s ...string) interface{} {
	val, err := j.Interface(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfStrings extracts an array of strings from some json
func (j *JsonQuery) ArrayOfStrings(s ...string) ([]string, error) {
	array, err := j.Array(s...)
	if err != nil {
		return []string{}, err
	}
	toReturn := make([]string, len(array))
	for index, val := range array {
		toReturn[index], err = stringFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfStrings extracts an array of strings from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfStrings(s ...string) []string {
	val, err := j.ArrayOfStrings(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfInts extracts an array of ints from some json
func (j *JsonQuery) ArrayOfInts(s ...string) ([]int, error) {
	array, err := j.Array(s...)
	if err != nil {
		return []int{}, err
	}
	toReturn := make([]int, len(array))
	for index, val := range array {
		toReturn[index], err = intFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfInts extracts an array of ints from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfInts(s ...string) []int {
	val, err := j.ArrayOfInts(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfFloats extracts an array of float64s from some json
func (j *JsonQuery) ArrayOfFloats(s ...string) ([]float64, error) {
	array, err := j.Array(s...)
	if err != nil {
		return []float64{}, err
	}
	toReturn := make([]float64, len(array))
	for index, val := range array {
		toReturn[index], err = floatFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfFloats extracts an array of float64s from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfFloats(s ...string) []float64 {
	val, err := j.ArrayOfFloats(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfBools extracts an array of bools from some json
func (j *JsonQuery) ArrayOfBools(s ...string) ([]bool, error) {
	array, err := j.Array(s...)
	if err != nil {
		return []bool{}, err
	}
	toReturn := make([]bool, len(array))
	for index, val := range array {
		toReturn[index], err = boolFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfBools extracts an array of bools from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfBools(s ...string) []bool {
	val, err := j.ArrayOfBools(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfObjects extracts an array of map[string]interface{} (objects) from some json
func (j *JsonQuery) ArrayOfObjects(s ...string) ([]map[string]interface{}, error) {
	array, err := j.Array(s...)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	toReturn := make([]map[string]interface{}, len(array))
	for index, val := range array {
		toReturn[index], err = objectFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfObjects extracts an array of map[string]interface{} (objects) from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfObjects(s ...string) []map[string]interface{} {
	val, err := j.ArrayOfObjects(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// ArrayOfArrays extracts an array of []interface{} (arrays) from some json
func (j *JsonQuery) ArrayOfArrays(s ...string) ([][]interface{}, error) {
	array, err := j.Array(s...)
	if err != nil {
		return [][]interface{}{}, err
	}
	toReturn := make([][]interface{}, len(array))
	for index, val := range array {
		toReturn[index], err = arrayFromInterface(val)
		if err != nil {
			return toReturn, err
		}
	}
	return toReturn, nil
}

// AsArrayOfArrays extracts an array of []interface{} (arrays) from some json, but panics on error so it can be used inline
func (j *JsonQuery) AsArrayOfArrays(s ...string) [][]interface{} {
	val, err := j.ArrayOfArrays(s...)
	if err != nil {
		if j.SingleValuePanicOnError {
			panic(err)
		}
	}
	return val
}

// Matrix2D is an alias for ArrayOfArrays
func (j *JsonQuery) Matrix2D(s ...string) ([][]interface{}, error) {
	return j.ArrayOfArrays(s...)
}

// AsMatrix2D is an alias for ArrayOfArrays
func (j *JsonQuery) AsMatrix2D(s ...string) [][]interface{} {
	return j.AsArrayOfArrays(s...)
}

// Recursively query a decoded json blob
func rquery(blob interface{}, s ...string) (interface{}, error) {
	var (
		val interface{}
		err error
	)

	// If there is only a single string argument and if that single string argument has either a "." or a "[" in it
	// the assume it is a path specification and disagregate it into an array of indexes.
	terms := s
	if len(s) == 1 && strings.IndexAny(s[0], ".[]") != -1 {
		terms = strings.FieldsFunc(s[0], func(c rune) bool {
			return c == '.' || c == '[' || c == ']'
		})
	}
	val = blob
	for _, q := range terms {
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

// query a json blob for a single field or index.  If query is a string, then
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
