package db

import (
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/datautils"
)

// A Decoder provides unmarshaling AttributeValues to Go value types.
type Decoder struct {
	// Instructs the decoder to decode AttributeValue Numbers as
	// Number type instead of float64 when the destination type
	// is interface{}. Similar to encoding/json.Number
	UseNumber bool
}

// UnmarshalMap is an alias for Unmarshal which unmarshals from
// a map of AttributeValues.
//
// The output value provided must be a non-nil pointer
func UnmarshalMap(m map[string]*dynamodb.AttributeValue, out *datautils.JSONObject) error {
	return NewDecoder().Decode(&dynamodb.AttributeValue{M: m}, out)
}

// NewDecoder creates a new Decoder with default configuration. Use
// the `opts` functional options to override the default configuration.
func NewDecoder(opts ...func(*Decoder)) *Decoder {
	d := &Decoder{}
	for _, o := range opts {
		o(d)
	}

	return d
}

// Decode will unmarshal an AttributeValue into a Go value type. An error
// will be return if the decoder is unable to unmarshal the AttributeValue
// to the provide Go value type.
//
// The output value provided must be a non-nil pointer
func (d *Decoder) Decode(av *dynamodb.AttributeValue, out *datautils.JSONObject, opts ...func(*Decoder)) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() || !v.IsValid() {
		return &dynamodbattribute.InvalidUnmarshalError{Type: reflect.TypeOf(out)}
	}

	return d.decode(av, v)
}

var stringInterfaceMapType = reflect.TypeOf(map[string]interface{}(nil))
var timeType = reflect.TypeOf(time.Time{})

func (d *Decoder) decode(av *dynamodb.AttributeValue, v reflect.Value) error {
	var u dynamodbattribute.Unmarshaler
	if av == nil || av.NULL != nil {
		u, v = indirect(v, true)
		if u != nil {
			return u.UnmarshalDynamoDBAttributeValue(av)
		}
		return d.decodeNull(v)
	}

	u, v = indirect(v, false)
	if u != nil {
		return u.UnmarshalDynamoDBAttributeValue(av)
	}

	switch {
	case len(av.B) != 0:
		panic("Not implemented B")
	case av.BOOL != nil:
		return d.decodeBool(av.BOOL, v)
	case len(av.BS) != 0:
		panic("Not implemented BS")
	case av.L != nil:
		return d.decodeList(av.L, v)
	case len(av.M) != 0:
		return d.decodeMap(av.M, v)
	case av.N != nil:
		return d.decodeNumber(av.N, v)
	case len(av.NS) != 0:
		panic("Not implemented NS")
	case av.S != nil:
		return d.decodeString(av.S, v)
	case len(av.SS) != 0:
		panic("Not implemented SS")
	}

	return nil
}

func (d *Decoder) decodeBool(b *bool, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool, reflect.Interface:
		v.Set(reflect.ValueOf(*b).Convert(v.Type()))
	default:
		return &dynamodbattribute.UnmarshalTypeError{Value: "bool", Type: v.Type()}
	}

	return nil
}

func (d *Decoder) decodeNumber(n *string, v reflect.Value) error {
	// Default to float64 for all numbers
	i, err := strconv.ParseFloat(*n, 64)
	if err != nil {
		return &dynamodbattribute.UnmarshalTypeError{Value: "number", Type: v.Type()}
	}
	v.Set(reflect.ValueOf(i))
	return nil
}

func (d *Decoder) decodeList(avList []*dynamodb.AttributeValue, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface:
		s := make([]interface{}, len(avList))
		for i, av := range avList {
			if err := d.decode(av, reflect.ValueOf(&s[i]).Elem()); err != nil {
				return err
			}
		}
		v.Set(reflect.ValueOf(s))
		return nil
	default:
		return &dynamodbattribute.UnmarshalTypeError{Value: "list", Type: v.Type()}
	}
}

func (d *Decoder) decodeMap(avMap map[string]*dynamodb.AttributeValue, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		t := v.Type()
		if t.Key().Kind() != reflect.String {
			return &dynamodbattribute.UnmarshalTypeError{Value: "map string key", Type: t.Key()}
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(t))
		}
	case reflect.Struct:
	case reflect.Interface:
		v.Set(reflect.MakeMap(stringInterfaceMapType))
		v = v.Elem()
	default:
		return &dynamodbattribute.UnmarshalTypeError{Value: "map", Type: v.Type()}
	}

	if v.Kind() == reflect.Map {
		for k, av := range avMap {
			key := reflect.ValueOf(k)
			elem := reflect.New(v.Type().Elem()).Elem()
			if err := d.decode(av, elem); err != nil {
				return err
			}
			v.SetMapIndex(key, elem)
		}
	} else if v.Kind() == reflect.Struct {
		panic("Not implemented struct")
	}

	return nil
}

func (d *Decoder) decodeNull(v reflect.Value) error {
	if v.IsValid() && v.CanSet() {
		v.Set(reflect.Zero(v.Type()))
	}

	return nil
}

func (d *Decoder) decodeString(s *string, v reflect.Value) error {
	// To maintain backwards compatibility with ConvertFrom family of methods which
	// converted strings to time.Time structs
	if v.Type().ConvertibleTo(timeType) {
		t, err := time.Parse(time.RFC3339, *s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t).Convert(v.Type()))
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(*s)
	case reflect.Interface:
		// Ensure type aliasing is handled properly
		v.Set(reflect.ValueOf(*s).Convert(v.Type()))
	default:
		return &dynamodbattribute.UnmarshalTypeError{Value: "string", Type: v.Type()}
	}

	return nil
}

// indirect will walk a value's interface or pointer value types. Returning
// the final value or the value a unmarshaler is defined on.
//
// Based on the enoding/json type reflect value type indirection in Go Stdlib
// https://golang.org/src/encoding/json/decode.go indirect func.
func indirect(v reflect.Value, decodingNull bool) (dynamodbattribute.Unmarshaler, reflect.Value) {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(dynamodbattribute.Unmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		v = v.Elem()
	}

	return nil, v
}
