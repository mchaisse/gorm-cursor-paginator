package paginator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// CursorEncoder encoder for cursor
type CursorEncoder interface {
	Encode(v interface{}) (string, error)
}

// NewCursorEncoder creates cursor encoder
func NewCursorEncoder(overKeys map[string]string, keys ...string) CursorEncoder {
	return &cursorEncoder{keys, overKeys}
}

type cursorEncoder struct {
	keys     []string
	overKeys map[string]string
}

func (e *cursorEncoder) Encode(v interface{}) (string, error) {
	b, err := e.marshalJSON(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (e *cursorEncoder) marshalJSON(value interface{}) ([]byte, error) {
	rv := toReflectValue(value)
	// reduce reflect value to underlying value
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	fields := make([]interface{}, len(e.keys))
	for i, key := range e.keys {
		fields[i] = rv.FieldByName(key).Interface()
	}
	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/* deprecated */

func encodeOld(rv reflect.Value, keys []string) string {
	fields := make([]string, len(keys))
	for index, key := range keys {
		if rv.Kind() == reflect.Ptr {
			fields[index] = convert(reflect.Indirect(rv).FieldByName(key).Interface())
		} else {
			fields[index] = convert(rv.FieldByName(key).Interface())
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(fields, ",")))
}

func convert(field interface{}) string {
	switch field.(type) {
	case time.Time:
		return fmt.Sprintf("%s?%s", field.(time.Time).UTC().Format(time.RFC3339Nano), fieldTime)
	default:
		return fmt.Sprintf("%v?%s", field, fieldString)
	}
}
