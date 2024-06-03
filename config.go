package hocon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Type of an hocon Value
type Type int

// Type constants
const (
	ObjectType Type = iota
	StringType
	ArrayType
	NumberType
	BooleanType
	NullType
	SubstitutionType
	ConcatenationType
	valueWithAlternativeType
)

// Config stores the root of the configuration tree
// and provides an API to retrieve configuration values with the path expressions
type Config struct {
	root Value
}

// String method returns the string representation of the Config object
func (c *Config) String() string { return c.root.String() }

func (c *Config) Json() string {
	js := make(map[string]interface{})

	err := json.Unmarshal([]byte(c.root.Json()), &js)
	if err != nil {
		panic(fmt.Sprintf("Error on parsing json: %s", err))
	}

	return jsonMarshal(js)
}

// GetRoot method returns the root value of the configuration
func (c *Config) GetRoot() Value {
	return c.root
}

// GetObject method finds the value at the given path and returns it as an Object, returns nil if the value is not found
func (c *Config) GetObject(path string) (Object, error) {
	value := c.Get(path)
	if value == nil {
		return nil, fmt.Errorf("config value not found at path: %s", path)
	}

	val, ok := value.(Object)
	if !ok {
		return nil, fmt.Errorf("config value at path: %s is not an object", path)
	}

	return val, nil
}

// GetConfig method finds the value at the given path and returns it as a Config, returns nil if the value is not found
func (c *Config) GetConfig(path string) (*Config, error) {
	value, err := c.GetObject(path)
	if err != nil {
		return nil, err
	}

	return value.ToConfig(), nil
}

// GetStringMap method finds the value at the given path and returns it as a map[string]Value
// returns nil if the value is not found
func (c *Config) GetStringMap(path string) (map[string]Value, error) {
	return c.GetObject(path)
}

// GetStringMapString method finds the value at the given path and returns it as a map[string]string
// returns nil if the value is not found
func (c *Config) GetStringMapString(path string) (map[string]string, error) {
	value := c.Get(path)
	if value == nil {
		return nil, fmt.Errorf("config value not found at path: %s", path)
	}

	object, ok := value.(Object)
	if !ok {
		return nil, fmt.Errorf("config value at path: %s is not an object", path)
	}

	var m = make(map[string]string, len(object))
	for k, v := range object {
		m[k] = v.String()
	}

	return m, nil
}

// GetArray method finds the value at the given path and returns it as an Array, returns nil if the value is not found
func (c *Config) GetArray(path string) (Array, error) {
	value := c.Get(path)
	if value == nil {
		return nil, fmt.Errorf("config value not found at path: %s", path)
	}

	val, ok := value.(Array)
	if !ok {
		return val, fmt.Errorf("config value at path: %s is not an array", path)
	}

	return val, nil
}

// GetIntSlice method finds the value at the given path and returns it as []int, returns nil if the value is not found
func (c *Config) GetIntSlice(path string) ([]int, error) {
	value := c.Get(path)
	if value == nil {
		return nil, fmt.Errorf("config value not found at path: %s", path)
	}

	arr, ok := value.(Array)
	if !ok {
		return nil, fmt.Errorf("config value at path: %s is not an array of integers", path)
	}

	slice := make([]int, 0, len(arr))
	for _, v := range arr {
		i, ok := v.(Int)
		if !ok {
			return nil, fmt.Errorf("config value at path: %s is not an array of integers", path)
		}

		slice = append(slice, int(i))
	}

	return slice, nil
}

// GetStringSlice method finds the value at the given path and returns it as []string
// returns nil if the value is not found
func (c *Config) GetStringSlice(path string) ([]string, error) {
	value := c.Get(path)
	if value == nil {
		return nil, fmt.Errorf("config value not found at path: %s", path)
	}

	arr, ok := value.(Array)
	if !ok {
		return nil, fmt.Errorf("config value at path: %s is not an array", path)
	}

	slice := make([]string, 0, len(arr))

	for _, v := range arr {
		slice = append(slice, v.String())
	}

	return slice, nil
}

// GetString method finds the value at the given path and returns it as a String
// returns empty string if the value is not found
func (c *Config) GetString(path string) (string, error) {
	value := c.Get(path)
	if value == nil {
		return "", fmt.Errorf("config value not found at path: %s", path)
	}

	return value.String(), nil
}

// GetInt method finds the value at the given path and returns it as an Int, returns zero if the value is not found
func (c *Config) GetInt(path string) (int, error) {
	value := c.Get(path)
	if value == nil {
		return 0, fmt.Errorf("config value not found at path: %s", path)
	}

	switch val := value.(type) {
	case Int:
		return int(val), nil
	case String:
		i, err := strconv.Atoi(string(val))
		if err != nil {
			return 0, fmt.Errorf("cannot parse value: %s to int", path)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("cannot parse value: %s to int", path)
	}
}

// GetFloat32 method finds the value at the given path and returns it as a Float32
// returns float32(0.0) if the value is not found
func (c *Config) GetFloat32(path string) (float32, error) {
	value := c.Get(path)
	if value == nil {
		return float32(0.0), fmt.Errorf("config value not found at path: %s", path)
	}

	switch val := value.(type) {
	case Float32:
		return float32(val), nil
	case Float64:
		return float32(val), nil
	case String:
		floatValue, err := strconv.ParseFloat(string(val), 32)
		if err != nil {
			return float32(0.0), fmt.Errorf("cannot parse value: %s to float32", path)
		}
		return float32(floatValue), nil
	default:
		return float32(0.0), fmt.Errorf("cannot parse value: %s to float32", path)
	}
}

// GetFloat64 method finds the value at the given path and returns it as a Float64
// returns 0.0 if the value is not found
func (c *Config) GetFloat64(path string) (float64, error) {
	value := c.Get(path)
	if value == nil {
		return 0.0, fmt.Errorf("config value not found at path: %s", path)
	}

	switch val := value.(type) {
	case Float64:
		return float64(val), nil
	case Float32:
		return float64(val), nil
	case String:
		floatValue, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return 0.0, fmt.Errorf("cannot parse value: %s to float64", path)
		}
		return floatValue, nil
	default:

		return 0.0, fmt.Errorf("cannot parse value: %s to float64", path)
	}
}

// GetBoolean method finds the value at the given path and returns it as a Boolean
// returns false if the value is not found
func (c *Config) GetBoolean(path string) (bool, error) {
	value := c.Get(path)
	if value == nil {
		return false, fmt.Errorf("config value not found at path: %s", path)
	}

	switch val := value.(type) {
	case Boolean:
		return bool(val), nil
	case String:
		switch val {
		case "true", "yes", "on":
			return true, nil
		case "false", "no", "off":
			return false, nil
		default:
			return false, fmt.Errorf("cannot parse value: %s to boolean", path)
		}
	default:
		return false, fmt.Errorf("cannot parse value: %s to boolean", path)
	}
}

// GetDuration method finds the value at the given path and returns it as a time.Duration
// returns 0 if the value is not found
func (c *Config) GetDuration(path string) (time.Duration, error) {
	value := c.Get(path)
	if value == nil {
		return 0, fmt.Errorf("config value not found at path: %s", path)
	}

	dur, ok := value.(Duration)
	if !ok {
		return 0, fmt.Errorf("cannot parse value: %s to Duration", path)
	}

	return time.Duration(dur), nil
}

// Get method finds the value at the given path and returns it without casting to any type
// returns nil if the value is not found
func (c *Config) Get(path string) Value {
	if c.root.Type() != ObjectType {
		return nil
	}

	return c.root.(Object).find(path)
}

// WithFallback method returns a new *Config (or the current config, if the given fallback doesn't get used)
// 1. merges the values of the current and fallback *Configs, if the root of both of them are of type Object
// for the same keys current values overrides the fallback values
// 2. if any of the *Configs has non-object root then returns the current *Config ignoring the fallback parameter
func (c *Config) WithFallback(fallback *Config) *Config {
	if current, ok := c.root.(Object); ok {
		if fallbackObject, ok := fallback.root.(Object); ok {
			resultConfig := fallbackObject.copy()
			mergeObjects(resultConfig, current)

			return resultConfig.ToConfig()
		}
	}

	return c
}

// Value interface represents a value in the configuration tree, all the value types implements this interface
type Value interface {
	Type() Type
	String() string
	Json() string
	isConcatenable() bool
}

// String represents a string value
type String string

// Type String
func (s String) Type() Type { return StringType }

func (s String) String() string {
	return string(s)
}

func (s String) Json() string {
	return jsonMarshal(s)
}

func (s String) isConcatenable() bool { return true }

// valueWithAlternative represents a value with Substitution which might override the original value
type valueWithAlternative struct {
	value       Value
	alternative *Substitution
}

func (s *valueWithAlternative) Type() Type { return valueWithAlternativeType }

func (s *valueWithAlternative) String() string {
	return fmt.Sprintf("(%s | %s)", s.value, s.alternative.String())
}

func (s *valueWithAlternative) Json() string {
	return jsonMarshal(s.String())
}

func (s *valueWithAlternative) isConcatenable() bool { return false }

// Object represents an object node in the configuration tree
type Object map[string]Value

// Type Object
func (o Object) Type() Type           { return ObjectType }
func (o Object) isConcatenable() bool { return false }

// String method returns the string representation of the Object
func (o Object) String() string {
	return o.Json()
}

func (o Object) Json() string {
	var builder strings.Builder

	itemsSize := len(o)
	i := 1

	builder.WriteString(objectStartToken)

	for key, value := range o {
		builder.WriteString(jsonMarshal(key))
		builder.WriteString(colonToken)

		if value != nil {
			builder.WriteString(value.Json())
		} else {
			builder.WriteString(string(null))
		}

		if i < itemsSize {
			builder.WriteString(", ")
		}
		i++
	}

	builder.WriteString(objectEndToken)

	return builder.String()
}

// ToConfig method converts object to *Config
func (o Object) ToConfig() *Config {
	return &Config{o}
}

func (o Object) find(path string) Value {
	keys := strings.Split(path, dotToken)
	size := len(keys)
	lastKey := keys[size-1]
	keysWithoutLast := keys[:size-1]
	object := o

	for _, key := range keysWithoutLast {
		value, ok := object[key]
		if !ok {
			return nil
		}

		object = value.(Object)
	}

	return object[lastKey]
}

func (o Object) copy() Object {
	result := Object{}

	for k, v := range o {
		subObject, ok := v.(Object)
		if ok {
			result[k] = subObject.copy()
		} else {
			result[k] = v
		}
	}

	return result
}

// Array represents an array node in the configuration tree
type Array []Value

// Type Array
func (a Array) Type() Type           { return ArrayType }
func (a Array) isConcatenable() bool { return false }

// String method returns the string representation of the Array
func (a Array) String() string {
	return a.Json()
}

func (a Array) Json() string {
	if len(a) == 0 {
		return "[]"
	}

	var builder strings.Builder

	builder.WriteString(arrayStartToken)
	builder.WriteString(a[0].Json())

	for _, value := range a[1:] {
		builder.WriteString(commaToken)
		builder.WriteString(value.Json())
	}

	builder.WriteString(arrayEndToken)

	return builder.String()
}

// Int represents an Integer value
type Int int

// Type Number
func (i Int) Type() Type           { return NumberType }
func (i Int) String() string       { return strconv.Itoa(int(i)) }
func (i Int) Json() string         { return i.String() }
func (i Int) isConcatenable() bool { return true }

// Float32 represents a Float32 value
type Float32 float32

// Type Number
func (f Float32) Type() Type           { return NumberType }
func (f Float32) String() string       { return strconv.FormatFloat(float64(f), 'e', -1, 32) }
func (f Float32) Json() string         { return fmt.Sprintf(`"%s"`, f.String()) }
func (f Float32) isConcatenable() bool { return false }

// Float64 represents a Float64 value
type Float64 float64

// Type Number
func (f Float64) Type() Type           { return NumberType }
func (f Float64) String() string       { return strconv.FormatFloat(float64(f), 'e', -1, 64) }
func (f Float64) Json() string         { return fmt.Sprintf(`"%s"`, f.String()) }
func (f Float64) isConcatenable() bool { return false }

// Boolean represents bool value
type Boolean bool

func newBooleanFromString(value string) (Boolean, error) {
	switch value {
	case "true", "yes", "on":
		return true, nil
	case "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse value: %s to Boolean", value)
	}
}

// Type Boolean
func (b Boolean) Type() Type           { return BooleanType }
func (b Boolean) String() string       { return strconv.FormatBool(bool(b)) }
func (b Boolean) Json() string         { return b.String() }
func (b Boolean) isConcatenable() bool { return true }

// Substitution refers to another value in the configuration tree
type Substitution struct {
	path     string
	optional bool
}

// Type Substitution
func (s *Substitution) Type() Type           { return SubstitutionType }
func (s *Substitution) isConcatenable() bool { return true }

// String method returns the string representation of the Substitution
func (s *Substitution) String() string {
	var builder strings.Builder

	builder.WriteString("${")

	if s.optional {
		builder.WriteString("?")
	}

	builder.WriteString(s.path)
	builder.WriteString("}")

	return builder.String()
}

func (s *Substitution) Json() string {
	return jsonMarshal(s.String())
}

// Null represents a null value
type Null string

const null Null = "null"

// Type Null
func (n Null) Type() Type           { return NullType }
func (n Null) String() string       { return string(null) }
func (n Null) Json() string         { return string(null) }
func (n Null) isConcatenable() bool { return true }

// Duration represents a duration value
type Duration time.Duration

// Type Duration
func (d Duration) Type() Type           { return StringType }
func (d Duration) String() string       { return time.Duration(d).String() }
func (d Duration) Json() string         { return fmt.Sprintf(`%d`, time.Duration(d).Milliseconds()) }
func (d Duration) isConcatenable() bool { return false }

type concatenation Array

func (c concatenation) Type() Type           { return ConcatenationType }
func (c concatenation) isConcatenable() bool { return true }
func (c concatenation) containsObject() bool {
	for _, value := range c {
		if value != nil && value.Type() == ObjectType {
			return true
		}
	}

	return false
}
func (c concatenation) String() string {
	var builder strings.Builder

	for _, value := range c {
		if value != nil {
			builder.WriteString(strings.Trim(value.String(), `"`))
		}
	}

	return builder.String()
}

func (c concatenation) Json() string {
	return jsonMarshal(c.String())
}

func jsonMarshal(v interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	_ = encoder.Encode(v)
	return strings.TrimSpace(buffer.String())
}
