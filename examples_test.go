package sfv_test

import (
	"fmt"
	"net/http"

	"github.com/shogo82148/go-sfv"
)

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
