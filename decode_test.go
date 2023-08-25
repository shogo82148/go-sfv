package sfv

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkDecodeInteger(b *testing.B) {
	v := []string{"-123456789012345"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeDecimal(b *testing.B) {
	v := []string{"-123456789012.345"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeString(b *testing.B) {
	v := []string{`"hello"`}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeToken(b *testing.B) {
	v := []string{"hello"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeBinary(b *testing.B) {
	v := []string{":AQID:"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeBoolean(b *testing.B) {
	v := []string{"?1"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeDate(b *testing.B) {
	v := []string{"@1659578233"}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

func BenchmarkDecodeDisplayString(b *testing.B) {
	v := []string{`%"%e3%81%93%e3%82%93%e3%81%ab%e3%81%a1%e3%82%8f%e3%80%9co(^^)o"`}
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(v)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
	}
}

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
	vv := []string{v}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got, err := DecodeItem(vv)
		if err != nil {
			b.Error(err)
		}
		runtime.KeepAlive(got)
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
	v, err := EncodeItem(item)
	if err != nil {
		b.Error(err)
	}

	var list []string
	for i := 0; i < 1024; i++ {
		list = append(list, v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := DecodeList(list); err != nil {
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
	var dict []string
	for i := 0; i < 1024; i++ {
		member, err := EncodeDictionary([]DictMember{
			{
				Key:  fmt.Sprintf("key%d", i),
				Item: item,
			},
		})
		if err != nil {
			b.Error(err)
		}
		dict = append(dict, member)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := DecodeDictionary(dict); err != nil {
			b.Error(err)
		} else {
			runtime.KeepAlive(got)
		}
	}
}
