package sfv

import (
	"encoding/base32"
	"errors"
	"fmt"
	"math"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestEncodeKeyGenerated(t *testing.T) {
	runEncodeTestCases(t, "./testdata/structured-field-tests/serialisation-tests/key-generated.json")
}

func TestEncodeNumber(t *testing.T) {
	runEncodeTestCases(t, "./testdata/structured-field-tests/serialisation-tests/number.json")
}

func TestEncodeStringGenerated(t *testing.T) {
	runEncodeTestCases(t, "./testdata/structured-field-tests/serialisation-tests/string-generated.json")
}

func TestEncodeTokenGenerated(t *testing.T) {
	runEncodeTestCases(t, "./testdata/structured-field-tests/serialisation-tests/token-generated.json")
}

func TestEncodeExtra(t *testing.T) {
	// This test is not part of the structured-field-tests suite.
	runEncodeTestCases(t, "./testdata/extra-serialization.json")
}

func TestEncodeIntegers(t *testing.T) {
	test := func(item Item) {
		t.Helper()
		val, err := EncodeItem(item)
		if err != nil {
			t.Error(err)
			return
		}
		if val != "123" {
			t.Errorf("want %q, got %q", "123", val)
		}
	}

	test(Item{
		Value: int(123),
	})
	test(Item{
		Value: uint(123),
	})
	test(Item{
		Value: int8(123),
	})
	test(Item{
		Value: uint8(123),
	})
	test(Item{
		Value: int16(123),
	})
	test(Item{
		Value: uint16(123),
	})
	test(Item{
		Value: int32(123),
	})
	test(Item{
		Value: uint32(123),
	})
	test(Item{
		Value: int64(123),
	})
	test(Item{
		Value: uint64(123),
	})
}

// Boundary value check
func TestEncodeIntegers_boundary(t *testing.T) {
	val, err := EncodeItem(Item{
		Value: int64(MaxInteger),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "999999999999999" {
		t.Errorf("want 999999999999999, got %q", val)
	}

	_, err = EncodeItem(Item{
		Value: int64(MaxInteger + 1),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	val, err = EncodeItem(Item{
		Value: int64(MinInteger),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "-999999999999999" {
		t.Errorf("want -999999999999999, got %q", val)
	}

	_, err = EncodeItem(Item{
		Value: int64(MinInteger - 1),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	val, err = EncodeItem(Item{
		Value: uint64(MaxInteger),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "999999999999999" {
		t.Errorf("want 999999999999999, got %q", val)
	}

	_, err = EncodeItem(Item{
		Value: uint64(MaxInteger + 1),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	_, err = EncodeItem(Item{
		Value: uint64(math.MaxInt64 + 1),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	if math.MaxInt >= MaxInteger {
		v := uint64(math.MaxInt64 + 1)
		_, err = EncodeItem(Item{
			Value: uint(v),
		})
		if err == nil {
			t.Error("want error, not not")
		}
	}
}

func TestEncodeDecimals(t *testing.T) {
	// Boundary value check
	val, err := EncodeItem(Item{
		Value: float64(MaxDecimal),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "999999999999.999" {
		t.Errorf("want 999999999999.999, got %q", val)
	}

	_, err = EncodeItem(Item{
		Value: math.Nextafter(float64(MaxDecimal), math.Inf(1)),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	val, err = EncodeItem(Item{
		Value: float64(MinDecimal),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "-999999999999.999" {
		t.Errorf("want -999999999999.999, got %q", val)
	}

	_, err = EncodeItem(Item{
		Value: math.Nextafter(float64(MinDecimal), math.Inf(-1)),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// encoding float32
	val, err = EncodeItem(Item{
		Value: float32(1.0),
	})
	if err != nil {
		t.Error(err)
	}
	if val != "1.0" {
		t.Errorf("want 1.0, got %q", val)
	}
}

func TestEncodeBadDisplayString(t *testing.T) {
	_, err := EncodeItem(Item{
		Value: DisplayString("\x80"),
	})
	if err == nil {
		t.Error("want error, not not")
	}
}

func runEncodeTestCases(t *testing.T, filename string) {
	cases, err := readTestCases(filename)
	if err != nil {
		t.Fatalf("failed to read %q: %v", filename, err)
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			canonical := strings.Join(tt.Canonical, ",")

			switch tt.HeaderType {
			case headerTypeItem:
				item, err := readItem(tt.Expected)
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}

				// test encoding
				encoded, err := EncodeItem(item)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
					return
				}
				if encoded != canonical {
					t.Errorf("want %q, got %q", canonical, encoded)
				}
			case headerTypeList:
				// test decoding
				list, err := readList(tt.Expected)
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}

				// test encoding
				encoded, err := EncodeList(list)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
					return
				}
				if encoded != canonical {
					t.Errorf("want %q, got %q", canonical, encoded)
				}
			case headerTypeDictionary:
				// test decoding
				dict, err := readDictionary(tt.Expected)
				if err != nil {
					t.Errorf("unexpected parse error: %v", err)
					return
				}

				// test encoding
				encoded, err := EncodeDictionary(dict)
				if tt.MustFail {
					if err == nil {
						t.Error("must fail, but no errors")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected encode error: %v", err)
					return
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

func readItem(v interface{}) (Item, error) {
	var ret Item
	item, ok := v.([]interface{})
	if !ok || len(item) != 2 {
		return Item{}, fmt.Errorf("invalid test data: want a (bare_item, parameters) tuple, got %v", v)
	}
	var err error
	ret.Value, err = readBareItem(item[0])
	if err != nil {
		return Item{}, err
	}
	ret.Parameters, err = readParameters(item[1])
	if err != nil {
		return Item{}, err
	}

	return ret, nil
}

func readBareItem(v interface{}) (Value, error) {
	switch v := v.(type) {
	case float64:
		i, frac := math.Modf(v)
		if frac == 0 {
			return int64(i), nil
		}
		return v, nil
	case string:
		return v, nil
	case bool:
		return v, nil
	case map[string]interface{}:
		typ, ok := v["__type"].(string)
		if !ok {
			return nil, errors.New("invalid test case: __type is not found")
		}
		value, ok := v["value"].(string)
		if !ok {
			return nil, errors.New("invalid test case: value is not found")
		}
		switch typ {
		case "token":
			return Token(value), nil
		case "binary":
			bin, err := base32.StdEncoding.DecodeString(value)
			if err != nil {
				return nil, fmt.Errorf("invalid test case: %v", err)
			}
			return bin, nil
		default:
			return nil, fmt.Errorf("invalid test case: unknown __type: %q", typ)
		}
	case []interface{}:
		var ret InnerList
		for _, item := range v {
			v, err := readItem(item)
			if err != nil {
				return nil, err
			}
			ret = append(ret, v)
		}
		return ret, nil
	}
	return nil, fmt.Errorf("error while parsing test case, unknown type: %T", v)
}

func readParameters(v interface{}) (Parameters, error) {
	dict, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid test case: unknown type: %T", v)
	}

	var params Parameters
	for _, item := range dict {
		kv, ok := item.([]interface{})
		if !ok || len(kv) != 2 {
			return nil, fmt.Errorf("invalid test case: want (key, value) tuple, got %v", item)
		}
		key, ok := kv[0].(string)
		if !ok {
			return nil, fmt.Errorf("invalid test case: invalid key type: %T", kv[0])
		}
		v, err := readBareItem(kv[1])
		if err != nil {
			return nil, err
		}
		params = append(params, Parameter{
			Key:   key,
			Value: v,
		})
	}
	return params, nil
}

func readList(v interface{}) (List, error) {
	list, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid test case: unknown type: %T", v)
	}

	var ret List
	for _, item := range list {
		v, err := readItem(item)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func readDictionary(v interface{}) (Dictionary, error) {
	dict, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid test case: unknown type: %T", v)
	}

	var ret Dictionary
	for _, item := range dict {
		kv, ok := item.([]interface{})
		if !ok || len(kv) != 2 {
			return nil, fmt.Errorf("invalid test case: want (key, value) tuple, got %v", item)
		}
		key, ok := kv[0].(string)
		if !ok {
			return nil, fmt.Errorf("invalid test case: invalid key type: %T", kv[0])
		}
		v, err := readItem(kv[1])
		if err != nil {
			return nil, err
		}
		ret = append(ret, DictMember{
			Key:  key,
			Item: v,
		})
	}
	return ret, nil
}

func TestEncode_invalidTypes(t *testing.T) {
	var err error

	// in Items
	_, err = EncodeItem(Item{
		Value: make(chan int),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// NaN
	_, err = EncodeItem(Item{
		Value: math.NaN(),
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// in Parameters
	_, err = EncodeItem(Item{
		Value: 1,
		Parameters: []Parameter{
			{
				Key:   "invalid",
				Value: make(chan int),
			},
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// in InnerLists
	_, err = EncodeList(List{
		Item{
			Value: InnerList{
				Item{
					Value: make(chan int),
				},
			},
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// in Lists
	_, err = EncodeList(List{
		Item{
			Value: make(chan int),
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}

	_, err = EncodeList(List{
		Item{
			Value: 1,
			Parameters: []Parameter{
				{
					Key:   "invalid",
					Value: make(chan int),
				},
			},
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}

	// in Dictionaries
	_, err = EncodeDictionary(Dictionary{
		DictMember{
			Key: "invalid",
			Item: Item{
				Value: make(chan int),
			},
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}

	_, err = EncodeDictionary(Dictionary{
		DictMember{
			Key: "invalid",
			Item: Item{
				Value: 1,
				Parameters: []Parameter{
					{
						Key:   "invalid",
						Value: make(chan int),
					},
				},
			},
		},
	})
	if err == nil {
		t.Error("want error, not not")
	}
}

func BenchmarkEncodeInteger(b *testing.B) {
	item := Item{
		Value: int64(-MaxInteger),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeDecimal(b *testing.B) {
	item := Item{
		Value: float64(-MaxDecimal),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeString(b *testing.B) {
	item := Item{
		Value: "hello world!",
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeToken(b *testing.B) {
	item := Item{
		Value: Token("hello"),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeBinary(b *testing.B) {
	item := Item{
		Value: []byte("こんにちわ〜o(^^)o"),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeBoolean(b *testing.B) {
	item := Item{
		Value: true,
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeDate(b *testing.B) {
	item := Item{
		Value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeDisplayString(b *testing.B) {
	item := Item{
		Value: DisplayString("こんにちわ〜o(^^)o"),
	}
	for i := 0; i < b.N; i++ {
		got, err := EncodeItem(item)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkEncodeItem(b *testing.B) {
	item := Item{
		Value: []byte("こんにちわ〜o(^^)o"),
		Parameters: []Parameter{
			{
				Key: "integer", Value: int64(1),
			},
			{
				Key: "decimal", Value: 1.234,
			},
			{
				Key: "binary", Value: []byte("こんにちわ〜o(^^)o"),
			},
			{
				Key: "token", Value: Token("hello"),
			},
			{
				Key: "string", Value: "hello world!",
			},
			{
				Key: "boolean", Value: false,
			},
			{
				Key: "date", Value: time.Unix(1659578233, 0),
			},
			{
				Key: "display-string", Value: DisplayString("こんにちわ〜o(^^)o"),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := EncodeItem(item); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}

func BenchmarkEncodeList(b *testing.B) {
	item := Item{
		Value: []byte("こんにちわ〜o(^^)o"),
		Parameters: []Parameter{
			{
				Key: "integer", Value: int64(1),
			},
			{
				Key: "decimal", Value: 1.234,
			},
			{
				Key: "binary", Value: []byte("こんにちわ〜o(^^)o"),
			},
			{
				Key: "token", Value: Token("hello"),
			},
			{
				Key: "string", Value: "hello world!",
			},
			{
				Key: "boolean", Value: false,
			},
			{
				Key: "date", Value: time.Unix(1659578233, 0),
			},
			{
				Key: "display-string", Value: DisplayString("こんにちわ〜o(^^)o"),
			},
		},
	}
	var list List
	for i := 0; i < 1024; i++ {
		list = append(list, item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := EncodeList(list); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}

func BenchmarkEncodeDictionary(b *testing.B) {
	item := Item{
		Value: []byte("こんにちわ〜o(^^)o"),
		Parameters: []Parameter{
			{
				Key: "integer", Value: int64(1),
			},
			{
				Key: "decimal", Value: 1.234,
			},
			{
				Key: "binary", Value: []byte("こんにちわ〜o(^^)o"),
			},
			{
				Key: "token", Value: Token("hello"),
			},
			{
				Key: "string", Value: "hello world!",
			},
			{
				Key: "boolean", Value: false,
			},
			{
				Key: "date", Value: time.Unix(1659578233, 0),
			},
			{
				Key: "display-string", Value: DisplayString("こんにちわ〜o(^^)o"),
			},
		},
	}
	var dict Dictionary
	for i := 0; i < 1024; i++ {
		dict = append(dict, DictMember{
			Key:  fmt.Sprintf("key%d", i),
			Item: item,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := EncodeDictionary(dict); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}
