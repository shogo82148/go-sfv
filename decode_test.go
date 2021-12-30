package sfv

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkDecodeItem(b *testing.B) {
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
		},
	}
	v, err := EncodeItem(item)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := DecodeItem([]string{v}); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}

func BenchmarkDecodeList(b *testing.B) {
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
		},
	}
	var list List
	for i := 0; i < 1024; i++ {
		list = append(list, item)
	}

	v, err := EncodeList(list)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := DecodeList([]string{v}); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}

func BenchmarkDecodeDictionary(b *testing.B) {
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
		},
	}
	var dict Dictionary
	for i := 0; i < 1024; i++ {
		dict = append(dict, DictMember{
			Key:  fmt.Sprintf("key%d", i),
			Item: item,
		})
	}

	v, err := EncodeDictionary(dict)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := DecodeDictionary([]string{v}); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}
