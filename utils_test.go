package sfv

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

type headerType string

const (
	headerTypeItem       headerType = "item"
	headerTypeList       headerType = "list"
	headerTypeDictionary headerType = "dictionary"
)

type testCase struct {
	// A string describing the test
	Name string `json:"name"`

	// An array of strings, each representing a field value received
	Raw []string `json:"raw"`

	// One of "item", "list", "dictionary"
	HeaderType headerType `json:"header_type"`

	// The expected data structure after parsing (if successful). Required, unless must_fail is true.
	Expected interface{} `json:"expected"`

	// boolean indicating whether the test is required to fail. Defaults to false.
	MustFail bool `json:"must_fail"`

	// boolean indicating whether failing this test is acceptable; for SHOULDs. Defaults to false.
	CanFail bool `json:"can_fail"`

	// An array of strings representing the canonical form of the field value,
	// if it is different from raw. Not applicable if must_fail is true.
	Canonical []string `json:"canonical"`
}

func runTestCases(t *testing.T, filename string) {
	cases, err := readTestCases(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %v", filename, err)
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			var canonical string
			if tt.Canonical != nil {
				canonical = strings.Join(tt.Canonical, ",")
			} else {
				canonical = strings.Join(tt.Raw, ",")
			}

			switch tt.HeaderType {
			case headerTypeItem:
				// test decoding
				item, err := DecodeItem(tt.Raw)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}
				checkItem(newTestContext(t), item, tt.Expected)

				// test encoding
				encoded, err := EncodeItem(item)
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
					return
				}
				if encoded != canonical {
					t.Errorf("want %q, got %q", canonical, encoded)
				}
			case headerTypeList:
				// test decoding
				list, err := DecodeList(tt.Raw)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}
				checkList(newTestContext(t), list, tt.Expected)

				// test encoding
				encoded, err := EncodeList(list)
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
				}
				if encoded != canonical {
					t.Errorf("want %q, got %q", canonical, encoded)
				}
			case headerTypeDictionary:
				// test decoding
				dict, err := DecodeDictionary(tt.Raw)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}
				checkDictionary(newTestContext(t), dict, tt.Expected)

				// test encoding
				encoded, err := EncodeDictionary(dict)
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
				}
				if encoded != canonical {
					t.Errorf("want %q, got %q", canonical, encoded)
				}
			default:
				t.Errorf("unknown header type: %q", tt.HeaderType)
			}
		})
	}
}

// read the test cases on https://github.com/httpwg/structured-field-tests
func readTestCases(filename string) ([]*testCase, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ret []*testCase
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type testContext struct {
	*testing.T
	Path string
}

func newTestContext(t *testing.T) *testContext {
	return &testContext{
		T:    t,
		Path: "root",
	}
}

func (t testContext) Items() *testContext {
	return &testContext{
		T:    t.T,
		Path: t.Path + "/Items",
	}
}

func (t testContext) Parameters() *testContext {
	return &testContext{
		T:    t.T,
		Path: t.Path + "/Parameters",
	}
}

func (t testContext) Key(name string) *testContext {
	return &testContext{
		T:    t.T,
		Path: t.Path + "." + name,
	}
}

func (t testContext) Index(idx int) *testContext {
	return &testContext{
		T:    t.T,
		Path: t.Path + fmt.Sprintf("[%d]", idx),
	}
}

func (t testContext) Errorf(format string, args ...interface{}) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	t.T.Error(msg, "in", t.Path)
}

func (t testContext) Error(msg string) {
	t.Helper()
	t.T.Error(msg, "in", t.Path)
}

func checkItem(t *testContext, got Item, want interface{}) {
	item, ok := want.([]interface{})
	if !ok || len(item) != 2 {
		t.Errorf("invalid test data: want a (bare_item, parameters) tuple, got %v", want)
		return
	}
	checkValue(t.Items(), got.Value, item[0])

	params, ok := item[1].([]interface{})
	if !ok {
		t.Errorf("invalid test data: want parameters, got %v", item[1])
	}
	checkParameter(t.Parameters(), got.Parameters, params)
}

func checkValue(t *testContext, got Value, want interface{}) {
	switch want := want.(type) {
	case float64:
		switch got := got.(type) {
		case int64:
			i, frac := math.Modf(want)
			if frac != 0 {
				t.Errorf("want %.3f, got %d", want, got)
			}
			if got != int64(i) {
				t.Errorf("want %d, got %d", int64(i), got)
			}
		case float64:
			// convert the numbers into string to avoid calculation errors.
			// the fractional component has at most three digits.
			s1 := fmt.Sprintf("%.3f", got)
			s2 := fmt.Sprintf("%.3f", want)
			if s1 != s2 {
				t.Errorf("want %s, got %s", s2, s1)
			}
		default:
			t.Errorf("unexpected type: %T", got)
		}
	case string:
		if v, ok := got.(string); ok {
			if v != want {
				t.Errorf("want %q, got %q", want, v)
			}
		} else {
			t.Errorf("want %T type, %T type", want, got)
		}
	case bool:
		if v, ok := got.(bool); ok {
			if v != want {
				t.Errorf("want %t, got %t", want, v)
			}
		} else {
			t.Errorf("want %T type, %T type", want, got)
		}
	case map[string]interface{}:
		typ, ok := want["__type"].(string)
		if !ok {
			t.Error("invalid test case: __type is not found")
			return
		}
		switch value := want["value"].(type) {
		case float64:
			switch typ {
			case "date":
				if v, ok := got.(time.Time); ok {
					want := time.Unix(int64(value), 0).UTC()
					if !v.Equal(want) {
						t.Errorf("want %v, got %v", want, v)
					}
				} else {
					t.Errorf("want time.Time, got %T type", got)
				}
			default:
				t.Errorf("invalid test case: unknown __type: %q", typ)
			}
		case string:
			switch typ {
			case "token":
				if v, ok := got.(Token); ok {
					if v != Token(value) {
						t.Errorf("want Token %q, got Token %q", value, v)
					}
				} else {
					t.Errorf("want Token, got %T type", got)
				}
			case "binary":
				if v, ok := got.([]byte); ok {
					want, err := base32.StdEncoding.DecodeString(value)
					if err != nil {
						t.Errorf("invalid test case: %v", err)
						return
					}
					if !bytes.Equal(v, want) {
						t.Errorf("want Binary %x, got Binary %x", want, v)
					}
				} else {
					t.Errorf("want []byte type, got %T type", got)
				}
			case "displaystring":
				if v, ok := got.(DisplayString); ok {
					if v != DisplayString(value) {
						t.Errorf("want DisplayString %q, got DisplayString %q", value, v)
					}
				} else {
					t.Errorf("want DisplayString, got %T type", got)
				}
			default:
				t.Errorf("invalid test case: unknown __type: %q", typ)
			}
		case nil:
			t.Error("invalid test case: value is not found")
			return
		default:
			t.Error("invalid test case: unsupported value type")
			return
		}
	case []interface{}:
		if v, ok := got.(InnerList); ok {
			if len(v) != len(want) {
				t.Errorf("unexpected length: want %d, got %d", len(want), len(v))
			}
			for i := 0; i < len(v) && i < len(want); i++ {
				checkItem(t.Index(i), v[i], want[i])
			}
		} else {
			t.Errorf("want InnerList type, got %T type", got)
		}

	default:
		t.Errorf("error while parsing test case, unknown type: %T", want)
	}
}

func checkParameter(t *testContext, got Parameters, want []interface{}) {
	if len(got) != len(want) {
		t.Errorf("invalid length: want %d, got %d", len(want), len(got))
		return
	}
	for i := 0; i < len(got) && i < len(want); i++ {
		kv, ok := want[i].([]interface{})
		if !ok || len(kv) != 2 {
			t.Errorf("invalid test case: want (key, value) tuple, got %v", want[i])
			return
		}
		key, ok := kv[0].(string)
		if !ok {
			t.Errorf("invalid test case: invalid key type: %T", kv[0])
			return
		}
		if got[i].Key != key {
			t.Errorf("unexpected key: want %q, got %q", key, got[i].Key)
		}
		checkValue(t.Key(key), got[i].Value, kv[1])
	}
}

func checkList(t *testContext, got List, want interface{}) {
	list, ok := want.([]interface{})
	if !ok {
		t.Errorf("error while parsing test case, unknown type: %T", want)
		return
	}
	if len(got) != len(list) {
		t.Errorf("invalid length: want %d, got %d", len(list), len(got))
		return
	}
	for i := 0; i < len(got) && i < len(list); i++ {
		want := list[i]
		got := got[i]
		checkItem(t.Index(i), got, want)
	}
}

func checkDictionary(t *testContext, got Dictionary, want interface{}) {
	dict, ok := want.([]interface{})
	if !ok {
		t.Errorf("error while parsing test case, unknown type: %T", want)
		return
	}
	if len(got) != len(dict) {
		t.Errorf("invalid length: want %d, got %d", len(dict), len(got))
		return
	}
	for i := 0; i < len(got) && i < len(dict); i++ {
		kv, ok := dict[i].([]interface{})
		if !ok || len(kv) != 2 {
			t.Errorf("invalid test case: want (key, value) tuple, got %v", dict[i])
			return
		}
		key, ok := kv[0].(string)
		if !ok {
			t.Errorf("invalid test case: invalid key type: %T", kv[0])
			return
		}
		got := got[i]
		if got.Key != key {
			t.Errorf("unexpected key: want %q, got %q", key, got.Key)
		}
		checkItem(t.Index(i), got.Item, kv[1])
	}
}
