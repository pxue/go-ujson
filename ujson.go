package ujson

import (
	"errors"
	"log"
)

type JSON struct {
	Root     interface{}
	mapItems []*MapItem
}

func NewFromBytes(data []byte, pool *MapPool) (*JSON, error) {
	j := &JSON{}
	if len(data) < 2 { // Need at least "{}"
		return nil, errors.New("no data passed in")
	}
	dec := NewDecoder(simpleStore{}, pool, data)
	root, items, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	j.Root = root
	j.mapItems = items
	return j, nil
}

func (j *JSON) Done() error {
	var err error
	for _, item := range j.mapItems {
		if err = item.Close(); err != nil {
			return err
		}
	}
	j.mapItems = nil
	return nil
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *JSON) Get(key string) *JSON {
	m, err := j.MaybeMap()
	if err == nil {
		if val, ok := m[key]; ok {
			return &JSON{val, j.mapItems}
		}
	}
	return &JSON{nil, j.mapItems}
}

// Map guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").Map() {
//			fmt.Println(k, v)
//		}
func (j *JSON) Map(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Map() received too many arguments %d", len(args))
	}

	a, err := j.MaybeMap()
	if err == nil {
		return a
	}

	return def
}

// MaybeMap type asserts to `map`
func (j *JSON) MaybeMap() (map[string]interface{}, error) {
	if m, ok := (j.Root).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// String guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").String(), js.Get("optional_param").String("my_default"))
func (j *JSON) String(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("String() received too many arguments %d", len(args))
	}

	s, err := j.MaybeString()
	if err == nil {
		return s
	}

	return def
}

// MaybeString type asserts to `string`
func (j *JSON) MaybeString() (string, error) {
	if s, ok := (j.Root).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Int64 guarantees the return of an `int64` (with optional default)
//
// useful when you explicitly want an `int64` in a single value return context:
//     myFunc(js.Get("param1").Int64(), js.Get("optional_param").Int64(5150))
func (j *JSON) Int64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Int64() received too many arguments %d", len(args))
	}

	i, err := j.MaybeInt64()
	if err == nil {
		return i
	}

	return def
}

// MaybeInt64 type asserts and parses an `int64`
func (j *JSON) MaybeInt64() (int64, error) {
	if n, ok := (j.Root).(numeric); ok {
		return n.Int64()
	}
	return -1, errors.New("type assertion to numeric failed")
}

func (j *JSON) Array(args ...*[]interface{}) *[]interface{} {
	var def *[]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Array() received too many arguments %d", len(args))
	}

	a, err := j.MaybeArray()
	if err == nil {
		return a
	}

	return def
}

func (j *JSON) MaybeArray() (*[]interface{}, error) {
	if a, ok := (j.Root).(*[]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

func (j *JSON) Float64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("Float64() received too many arguments %d", len(args))
	}

	i, err := j.MaybeFloat64()
	if err == nil {
		return i
	}

	return def
}

// MaybeInt64 type asserts and parses an `int64`
func (j *JSON) MaybeFloat64() (float64, error) {
	if n, ok := (j.Root).(numeric); ok {
		return n.Float64()
	}
	return -1, errors.New("type assertion to numeric failed")
}

func (j *JSON) Bool() bool {
	if b, ok := (j.Root).(bool); ok {
		return b
	}

	return false
}
