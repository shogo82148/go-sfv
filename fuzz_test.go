//go:build go1.18
// +build go1.18

package sfv

import (
	"reflect"
	"testing"
)

func FuzzDecodeItem(f *testing.F) {
	f.Add(`2; foourl="https://foo.example.com/"`)
	f.Add(`1; a; b=?0`)
	f.Add(`5`)
	f.Add(`5; foo=bar`)
	f.Add(`42`)
	f.Add(`4.5`)
	f.Add(`"Hello world"`)
	f.Add(`:cHJldGVuZCB0aGlzIGlzIGJpbmFyeSBjb250ZW50Lg==:`)

	f.Add("")
	f.Add("  1  ")
	f.Add("42")
	f.Add("0")
	f.Add("-0")
	f.Add("-42")
	f.Add("042")
	f.Add("-042")
	f.Add("00")
	f.Add("123456789012345")
	f.Add("-123456789012345")
	f.Add("1.23")
	f.Add("-1.23")
	f.Add("123456789012.1")
	f.Add("1.123")
	f.Add("-1.123")
	f.Add(`"foobar"`)
	f.Add(`""`)
	f.Add(`"foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo foo "`)
	f.Add(`"   "`)
	f.Add(`"foo \"bar\" \\ baz"`)
	f.Add("a_b-c.d3:f%00/*")
	f.Add("fooBar")
	f.Add("FooBar")
	f.Add(":aGVsbG8=:")
	f.Add("::")
	f.Add(":aGVsbG8:")
	f.Add(":/+Ah:")
	f.Add("?0")
	f.Add("?1")
	f.Fuzz(func(t *testing.T, field string) {
		item, err := DecodeItem([]string{field})
		if err != nil {
			t.Skip()
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
	f.Add(`"foo", "bar", "It was the best of times."`)
	f.Add("foo, bar")
	f.Add(`("foo" "bar"), ("baz"), ("bat" "one"), ()`)
	f.Add(`("foo"; a=1;b=2);lvl=5, ("bar" "baz");lvl=1`)
	f.Add(`abc;a=1;b=2; cde_456, (ghi;jk=4 l);q="9";r=w`)

	f.Add("1, 42")
	f.Add("")
	f.Add("  42, 43")
	f.Add("42")
	f.Add("1,42")
	f.Add("1 , 42")
	f.Add("1\t,\t42")
	f.Add("(1 2), (42 43)")
	f.Add("(42)")
	f.Add("()")
	f.Add("(1),(),(42)")
	f.Add("(  1  42  )")
	f.Fuzz(func(t *testing.T, field string) {
		list, err := DecodeList([]string{field})
		if err != nil {
			t.Skip()
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
	f.Add(`en="Applepie", da=:w4ZibGV0w6ZydGU=:`)
	f.Add("a=?0, b, c; foo=bar")
	f.Add("rating=1.5, feelings=(joy sadness)")
	f.Add("a=(1 2), b=3, c=4;aa=bb, d=(5 6);valid")
	f.Add("foo=1, bar=2")

	f.Add(`en="Applepie", da=:w4ZibGV0w6ZydGUK:`)
	f.Add("")
	f.Add("a=1")
	f.Add("a=(1 2)")
	f.Add("a=(1)")
	f.Add("a=()")
	f.Add("a=1,b=2")
	f.Add("a=1 ,  b=2")
	f.Add("a=1\t,\tb=2")
	f.Add("     a=1 ,  b=2")
	f.Add("a=1, b, c=3")
	f.Add("a, b, c")
	f.Add("a, b=2")
	f.Add("a=1, b")
	f.Add("a=1, b;foo=9, c=3")
	f.Add("a=1, b=?1;foo=9, c=3")
	f.Add("a=1,b=2,a=3")
	f.Fuzz(func(t *testing.T, field string) {
		dict, err := DecodeDictionary([]string{field})
		if err != nil {
			t.Skip()
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
