package sfv

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
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
		t.Run(tt.Name, func(t *testing.T) {
			switch tt.HeaderType {
			case headerTypeItem:
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
			case headerTypeList:
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
			case headerTypeDictionary:
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
		if got, ok := got.(string); ok {
			if got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		} else {
			t.Errorf("want %T type, %T type", want, got)
		}
	case bool:
		if got, ok := got.(bool); ok {
			if got != want {
				t.Errorf("want %t, got %t", want, got)
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
		value, ok := want["value"].(string)
		if !ok {
			t.Error("invalid test case: value is not found")
			return
		}
		switch typ {
		case "token":
			if got, ok := got.(Token); ok {
				if got != Token(value) {
					t.Errorf("want Token %q, got Token %q", value, got)
				}
			} else {
				t.Errorf("want Token, got %T type", got)
			}
		case "binary":
			if got, ok := got.([]byte); ok {
				want, err := base32.StdEncoding.DecodeString(value)
				if err != nil {
					t.Errorf("invalid test case: %v", err)
					return
				}
				if !bytes.Equal(got, want) {
					t.Errorf("want Binary %x, got Binary %x", want, got)
				}
			} else {
				t.Errorf("want []byte type, got %T type", got)
			}
		default:
			t.Errorf("invalid test case: unknown __type: %q", typ)
		}
	case []interface{}:
		t.Error("TODO: implement")
	default:
		t.Errorf("error while parsing test case, unknown type: %T", want)
	}
}

func checkParameter(t *testContext, got Parameters, want []interface{}) {
	if len(got) != len(want) {
		t.Errorf("invalid length: want %d, got %d", len(want), len(got))
		return
	}
	for i := range want {
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
}
