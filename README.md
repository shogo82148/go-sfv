# go-sfv

[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/go-sfv.svg)](https://pkg.go.dev/github.com/shogo82148/go-sfv)
[![Test](https://github.com/shogo82148/go-sfv/actions/workflows/test.yml/badge.svg)](https://github.com/shogo82148/go-sfv/actions/workflows/test.yml)

Go implementation for [RFC 9651 Structured Field Values for HTTP](https://www.rfc-editor.org/rfc/rfc9651.html) (SFV).

## Synopsis

### Decoding Structured Field Values

```go
h := make(http.Header)
// Decoding Items
item, err := sfv.DecodeItem(h.Values("Example-Hdr"))
switch val := item.Value.(type) {
case int64:             // Integers
case float64:           // Decimals
case string:            // Strings
case sfv.Token:         // Tokens
case bool:              // Booleans
case time.Time:         // Dates
case sfv.DisplayString: // Display Strings
}

// Decoding Lists
list, err := sfv.DecodeList(h.Values("Example-Hdr"))

// Decoding Dictionaries
dict, err := sfv.DecodeDictionary(h.Values("Example-Hdr"))
```

### Encoding Structured Field Values

```go
// Encoding Items
val, err := sfv.EncodeItem(item)

// Encoding Lists
val, err := sfv.EncodeList(list)

// Encoding Dictionaries
val, err := sfv.EncodeDictionary(dict)
```

## Supported Data Types

SFV types are mapped to Go types as described in this section.
Note that only **Lists**(`sfv.List`), **Dictionaries**(`sfv.Dictionary`), and **Items** (`sfv.Item`) can be in a top-level.

### Values of Items

The `sfv.Value` is defined as the following:

```go
type Value any
```

The actual type might be one of them:

| Type of SFV   | Example of SFV     | Type in Go          | Example in Go              |
| ------------- | ------------------ | ------------------- | -------------------------- |
| Integer       | `10`               | `int64`             | `int64(10)`                |
| Decimal       | `3.14`             | `float64`           | `float64(3.14)`            |
| String        | `"hello"`          | `string`            | `"hello"`                  |
| Token         | `x`                | `sfv.Token`         | `sfv.Token("x")`           |
| Byte Seq      | `:AQID:`           | `[]byte`            | `[]byte{1, 2, 3}`          |
| Boolean       | `?1`               | `bool`              | `true`                     |
| Date          | `@1659578233`      | `time.Time`         | `time.Unix(1659578233, 0)` |
| DisplayString | `%"f%c3%bc%c3%bc"` | `sfv.DisplayString` | `sfv.DisplayString("füü")` |
| Inner List    | `(1 2)`            | `sfv.InnerList`     | `sfv.InnerList{}`          |

### Parameters of Items

**Parameters** are ordered map of key-value pairs, however Go's `map` types are unordered.
So `sfv.Parameters` is defined by a slice of `sfv.Parameter` that is a key-value pair.

```go
type Parameter struct {
	Key   string
	Value Value
}

type Parameters []Parameter
```

### Lists

**Lists** are decoded to `sfv.List`.

```go
type List []Item
```

### Inner Lists

**Inner Lists** are decoded to `sfv.InnerList`.

```go
type InnerList []Item
```

Note that `sfv.InnerList` can't contain `sfv.InnerList` itself.

```go
// Encoding this will fail.
innerList := sfv.InnerList{
    {
        Value: sfv.InnerList{},
    },
}
```

### Dictionaries

**Dictionaries** are ordered maps of key-value pairs, however Go's `map` types are unordered.
So `sfv.Dictionary` is defined by a slice of `sfv.DictMember` that is a key-value pair.

```go
type DictMember struct {
	Key  string
	Item Item
}

type Dictionary []DictMember
```

## References

- [RFC 9651 Structured Field Values for HTTP](https://www.rfc-editor.org/rfc/rfc9651.html)
- [Structured Field Values による Header Field の構造化](https://blog.jxck.io/entries/2021-01-31/structured-field-values.html)
- [IETF RFC における ABNF と Parsing Algorithm の関係](https://blog.jxck.io/entries/2023-05-17/abnf-or-algorithm-in-rfc.html)
