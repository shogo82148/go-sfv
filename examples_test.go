package sfv_test

import (
	"fmt"
	"net/http"

	"github.com/shogo82148/go-sfv"
)

func ExampleEncodeList() {
	list := sfv.List{
		{
			Value: sfv.Token("foo"),
		},
		{
			Value: sfv.Token("bar"),
		},
	}
	val, err := sfv.EncodeList(list)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	//Output:
	// foo, bar
}

func ExampleDecodeList() {
	h := make(http.Header)
	h.Add("Example-Hdr", "foo")
	h.Add("Example-Hdr", "bar")

	list, err := sfv.DecodeList(h.Values("Example-Hdr"))
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		fmt.Println(item.Value)
	}

	//Output:
	// foo
	// bar
}

func ExampleEncodeDictionary() {
	dict := sfv.Dictionary{
		{
			Key: "foo",
			Item: sfv.Item{
				Value: int64(1),
			},
		},
		{
			Key: "bar",
			Item: sfv.Item{
				Value: int64(2),
			},
		},
	}

	val, err := sfv.EncodeDictionary(dict)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	//Output:
	// foo=1, bar=2
}

func ExampleDecodeDictionary() {
	h := make(http.Header)
	h.Add("Example-Hdr", "foo=1")
	h.Add("Example-Hdr", "bar=2")

	dict, err := sfv.DecodeDictionary(h.Values("Example-Hdr"))
	if err != nil {
		panic(err)
	}

	for _, member := range dict {
		fmt.Println(member.Key, member.Item.Value)
	}

	//Output:
	// foo 1
	// bar 2
}

func ExampleEncodeItem() {
	item := sfv.Item{
		Value: int64(2),
		Parameters: sfv.Parameters{
			{
				Key:   "foourl",
				Value: "https://foo.example.com/",
			},
		},
	}
	val, err := sfv.EncodeItem(item)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	//Output:
	// 2;foourl="https://foo.example.com/"
}

func ExampleDecodeItem() {
	item, err := sfv.DecodeItem([]string{`2; foourl="https://foo.example.com/"`})
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Value)
	fmt.Println(item.Parameters.Get("foourl"))

	//Output:
	// 2
	// https://foo.example.com/
}
