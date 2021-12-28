package sfv

import (
	"encoding/json"
	"fmt"
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
		T: t,
	}
}

func checkItem(t *testContext, got Item, want interface{}) {
	item, ok := want.([]interface{})
	if !ok || len(item) != 2 {
		t.Errorf("invalid test data: want a (bare_item, parameters) tuple, got %v", want)
		return
	}
	checkValue(t, got.Value, item[0])
}

func checkValue(t *testContext, got Value, want interface{}) {
	switch want := want.(type) {
	case float64:
		switch got := got.(type) {
		case int64:
			t.Error("TODO: implement")
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
		t.Error("TODO: implement")
	case bool:
		if got, ok := got.(bool); ok {
			if got != want {
				t.Errorf("want %t, got %t", want, got)
			}
		} else {
			t.Errorf("want %T type, %T type", want, got)
		}
	case map[string]interface{}:
		t.Error("TODO: implement")
	case []interface{}:
		t.Error("TODO: implement")
	default:
		t.Errorf("error while parsing test case, unknown type: %T", want)
	}
}
