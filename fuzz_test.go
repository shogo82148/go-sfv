//go:build go1.18
// +build go1.18

package sfv

import (
	"reflect"
	"testing"
)

func FuzzDecodeItem(f *testing.F) {
	addFuzzingData(f)
	f.Fuzz(func(t *testing.T, field string) {
		item, err := DecodeItem([]string{field})
		if err != nil {
			t.Skip(field)
		}
		field2, err := EncodeItem(item)
		if err != nil {
			t.Fatal(err)
		}
		item2, err := DecodeItem([]string{field2})
		if err != nil {
			t.Fatalf("DecodeItem failed to decode %q: %v", field2, err)
		}
		if !reflect.DeepEqual(item, item2) {
			t.Errorf("DecodeItem different query after being encoded\nbefore: %v\nafter: %v", item, item2)
		}
	})
}

func FuzzDecodeList(f *testing.F) {
	addFuzzingData(f)
	f.Fuzz(func(t *testing.T, field string) {
		list, err := DecodeList([]string{field})
		if err != nil {
			t.Skip(field)
		}
		field2, err := EncodeList(list)
		if err != nil {
			t.Fatal(err)
		}
		list2, err := DecodeList([]string{field2})
		if err != nil {
			t.Fatalf("DecodeList failed to decode %q: %v", field2, err)
		}
		if !reflect.DeepEqual(list, list2) {
			t.Errorf("DecodeList different query after being encoded\nbefore: %v\nafter: %v", list, list2)
		}
	})
}

func FuzzDecodeDictionary(f *testing.F) {
	addFuzzingData(f)
	f.Fuzz(func(t *testing.T, field string) {
		dict, err := DecodeDictionary([]string{field})
		if err != nil {
			t.Skip(field)
		}
		field2, err := EncodeDictionary(dict)
		if err != nil {
			t.Fatal(err)
		}
		dict2, err := DecodeDictionary([]string{field2})
		if err != nil {
			t.Fatalf("DecodeDictionary failed to decode %q: %v", field2, err)
		}
		if !reflect.DeepEqual(dict, dict2) {
			t.Errorf("DecodeDictionary different query after being encoded\nbefore: %v\nafter: %v", dict, dict2)
		}
	})
}

func addFuzzingData(f *testing.F) {
	addFuzzingDataFile(f, "./testdata/structured-field-tests/serialisation-tests/key-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/serialisation-tests/number.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/serialisation-tests/string-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/serialisation-tests/token-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/binary.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/binary.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/binary.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/boolean.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/date.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/dictionary.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/examples.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/item.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/key-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/large-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/list.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/number-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/number.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/param-dict.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/param-list.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/param-listlist.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/string-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/string.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/token-generated.json")
	addFuzzingDataFile(f, "./testdata/structured-field-tests/token.json")
	addFuzzingDataFile(f, "./testdata/extra.json")
}

func addFuzzingDataFile(f *testing.F, filename string) {
	cases, err := readTestCases(filename)
	if err != nil {
		f.Fatalf("failed to read %q: %v", filename, err)
	}
	for _, tt := range cases {
		if len(tt.Raw) == 1 {
			f.Add(tt.Raw[0])
		}
		if len(tt.Canonical) == 1 {
			f.Add(tt.Canonical[0])
		}
	}
}
