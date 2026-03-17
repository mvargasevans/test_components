package contracttests

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/example/components/ipc"
)

// dialWithRetry retries dialing the socket up to maxAttempts times.
func dialWithRetry(path string, maxAttempts int) (*ipc.Conn, error) {
	var (
		conn *ipc.Conn
		err  error
	)
	for i := 0; i < maxAttempts; i++ {
		conn, err = ipc.Dial(path)
		if err == nil {
			return conn, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, err
}

// assertFieldValue checks that a field exists, has the expected type, and equals the example value.
func assertFieldValue(t *testing.T, fields map[string]interface{}, name, expectedType, example string) {
	t.Helper()
	assertField(t, fields, name, expectedType)
	v, ok := fields[name]
	if !ok {
		return
	}
	switch expectedType {
	case "string":
		if got, ok := v.(string); ok && got != example {
			t.Errorf("field %q: want %q, got %q", name, example, got)
		}
	case "float64":
		expected, _ := strconv.ParseFloat(example, 64)
		if got, ok := v.(float64); ok && got != expected {
			t.Errorf("field %q: want %v, got %v", name, expected, got)
		}
	case "int":
		expected, _ := strconv.ParseInt(example, 10, 64)
		if got, ok := v.(float64); ok && int64(got) != expected {
			t.Errorf("field %q: want %v, got %v", name, expected, int64(got))
		}
	case "bool":
		expected, _ := strconv.ParseBool(example)
		if got, ok := v.(bool); ok && got != expected {
			t.Errorf("field %q: want %v, got %v", name, expected, got)
		}
	default:
		t.Errorf("field %q: unknown type %q", name, fmt.Sprintf("%T", v))
	}
}

// assertField checks that a field exists in the decoded JSON map and has the expected Go type.
func assertField(t *testing.T, fields map[string]interface{}, name, expectedType string) {
	t.Helper()
	v, ok := fields[name]
	if !ok {
		t.Errorf("response missing field %q", name)
		return
	}
	switch expectedType {
	case "string":
		if _, ok := v.(string); !ok {
			t.Errorf("field %q: want string, got %T", name, v)
		}
	case "int":
		switch v.(type) {
		case float64, json.Number:
			// JSON numbers decode as float64 or json.Number; both are acceptable for int fields.
		default:
			t.Errorf("field %q: want int (numeric), got %T", name, v)
		}
	case "float64":
		if _, ok := v.(float64); !ok {
			t.Errorf("field %q: want float64, got %T", name, v)
		}
	case "bool":
		if _, ok := v.(bool); !ok {
			t.Errorf("field %q: want bool, got %T", name, v)
		}
	}
}
