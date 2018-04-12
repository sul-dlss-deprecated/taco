package db

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/datautils"
)

type emptyOrigError struct{}

// An InvalidMarshalError is an error type representing an error
// occurring when marshaling a Go value type to an AttributeValue.
type InvalidMarshalError struct {
	emptyOrigError
	msg string
}

// Error returns the string representation of the error.
// satisfying the error interface
func (e *InvalidMarshalError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Message())
}

// Code returns the code of the error, satisfying the awserr.Error
// interface.
func (e *InvalidMarshalError) Code() string {
	return "InvalidMarshalError"
}

// Message returns the detailed message of the error, satisfying
// the awserr.Error interface.
func (e *InvalidMarshalError) Message() string {
	return e.msg
}

// MarshalMap converts a datautils.JSONObject to a map of dynamodb.AttributeValue
// values   suitable for making an insert.
// This is requred, because the dynamodbattribute.MarshalMap casts empty slices
// to nil, and we think the empty slice is important.
// See https://github.com/aws/aws-sdk-go/issues/682 and https://github.com/aws/aws-sdk-go-v2/issues/115
func MarshalMap(in datautils.JSONObject) (map[string]*dynamodb.AttributeValue, error) {
	av, err := NewEncoder().Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return av.M, nil
}

// An Encoder provides marshaling Go value types to AttributeValues.
type Encoder struct {
	// Empty strings, "", will be marked as NULL AttributeValue types.
	// Empty strings are not valid values for DynamoDB. Will not apply
	// to lists, sets, or maps. Use the struct tag `omitemptyelem`
	// to skip empty (zero) values in lists, sets and maps.
	//
	// Enabled by default.
	NullEmptyString bool
}

// NewEncoder creates a new Encoder with default configuration. Use
// the `opts` functional options to override the default configuration.
func NewEncoder(opts ...func(*Encoder)) *Encoder {
	e := &Encoder{
		NullEmptyString: true,
	}
	for _, o := range opts {
		o(e)
	}

	return e
}

// Encode will marshal a Go value type to an AttributeValue. Returning
// the AttributeValue constructed or error.
func (e *Encoder) Encode(in interface{}) (*dynamodb.AttributeValue, error) {
	av := &dynamodb.AttributeValue{}
	if err := e.encode(av, reflect.ValueOf(in)); err != nil {
		return nil, err
	}

	return av, nil
}

func (e *Encoder) encode(av *dynamodb.AttributeValue, v reflect.Value) error {
	// Handle both pointers and interface conversion into types
	v = valueElem(v)

	switch v.Kind() {
	case reflect.Invalid:
		encodeNull(av)
	case reflect.Struct:
		panic("invalid type Struct") // return e.encodeStruct(av, v, fieldTag)
	case reflect.Map:
		return e.encodeMap(av, v)
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(av, v)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// do nothing for unsupported types
	default:
		return e.encodeScalar(av, v)
	}

	return nil
}

func (e *Encoder) encodeMap(av *dynamodb.AttributeValue, v reflect.Value) error {
	av.M = map[string]*dynamodb.AttributeValue{}
	for _, key := range v.MapKeys() {
		keyName := fmt.Sprint(key.Interface())
		if keyName == "" {
			return &InvalidMarshalError{msg: "map key cannot be empty"}
		}

		elemVal := v.MapIndex(key)
		elem := &dynamodb.AttributeValue{}
		err := e.encode(elem, elemVal)
		if err != nil {
			return err
		}

		av.M[keyName] = elem
	}

	return nil
}

var byteSliceType = reflect.ValueOf([]byte{}).Type()
var numberType = reflect.TypeOf(dynamodbattribute.Number(""))

func (e *Encoder) encodeSlice(av *dynamodb.AttributeValue, v reflect.Value) error {
	switch v.Type().Elem().Kind() {
	case reflect.Uint8:
		slice := reflect.MakeSlice(byteSliceType, v.Len(), v.Len())
		reflect.Copy(slice, v)

		b := slice.Bytes()
		av.B = append([]byte{}, b...)
	default:
		var elemFn func(dynamodb.AttributeValue) error

		// List
		av.L = make([]*dynamodb.AttributeValue, 0, v.Len())
		elemFn = func(elem dynamodb.AttributeValue) error {
			av.L = append(av.L, &elem)
			return nil
		}

		if _, err := e.encodeList(v, elemFn); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeList(v reflect.Value, elemFn func(dynamodb.AttributeValue) error) (int, error) {
	count := 0
	for i := 0; i < v.Len(); i++ {
		elem := dynamodb.AttributeValue{}
		err := e.encode(&elem, v.Index(i))
		if err != nil {
			return 0, err
		}

		if err := elemFn(elem); err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

func (e *Encoder) encodeScalar(av *dynamodb.AttributeValue, v reflect.Value) error {
	if v.Type() == numberType {
		s := v.String()
		av.N = &s
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		av.BOOL = new(bool)
		*av.BOOL = v.Bool()
	case reflect.String:
		if err := e.encodeString(av, v); err != nil {
			return err
		}
	default:
		// Fallback to encoding numbers, will return invalid type if not supported
		if err := e.encodeNumber(av, v); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeNumber(av *dynamodb.AttributeValue, v reflect.Value) error {
	var out string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = encodeInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = encodeUint(v.Uint())
	case reflect.Float32, reflect.Float64:
		out = encodeFloat(v.Float())
	default:
		return &unsupportedMarshalTypeError{Type: v.Type()}
	}

	av.N = &out

	return nil
}

func (e *Encoder) encodeString(av *dynamodb.AttributeValue, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if len(s) == 0 && e.NullEmptyString {
			encodeNull(av)
		} else {
			av.S = &s
		}
	default:
		return &unsupportedMarshalTypeError{Type: v.Type()}
	}

	return nil
}

func encodeInt(i int64) string {
	return strconv.FormatInt(i, 10)
}
func encodeUint(u uint64) string {
	return strconv.FormatUint(u, 10)
}
func encodeFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func encodeNull(av *dynamodb.AttributeValue) {
	t := true
	*av = dynamodb.AttributeValue{NULL: &t}
}

func valueElem(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v
}

// An unsupportedMarshalTypeError represents a Go value type
// which cannot be marshaled into an AttributeValue and should
// be skipped by the marshaler.
type unsupportedMarshalTypeError struct {
	emptyOrigError
	Type reflect.Type
}

// Error returns the string representation of the error.
// satisfying the error interface
func (e *unsupportedMarshalTypeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Message())
}

// Code returns the code of the error, satisfying the awserr.Error
// interface.
func (e *unsupportedMarshalTypeError) Code() string {
	return "unsupportedMarshalTypeError"
}

// Message returns the detailed message of the error, satisfying
// the awserr.Error interface.
func (e *unsupportedMarshalTypeError) Message() string {
	return "Go value type " + e.Type.String() + " is not supported"
}
