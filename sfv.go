package sfv

import "errors"

// Token is a token defined in RFC 8941 Section 3.3.4. Tokens.
type Token string

// Value is a bare item.
// It might be Integers, Decimals, Strings, Tokens, Byte Sequences, Booleans or Inner Lists.
// It's type is one of these:
//
//     int64 for Integers
//     float64 for Decimals
//     string for Strings
//     Token for Tokens
//     []byte for Byte Sequences
//     bool for Booleans
//     InnerList for Inner Lists
type Value interface{}

// Parameter is a key-value pair of Parameters.
type Parameter struct {
	Key   string
	Value Value
}

// Parameters are an ordered map of key-value pairs defined in RFC 8941 Section 3.1.2. Parameters.
type Parameters []Parameter

// Get returns the last value associated with the given key.
// If there are no values associated with the key, Get returns Value(nil).
func (param Parameters) Get(key string) Value {
	// In many cases, there are a few parameters.
	// So Linear searching is enough to handle them.

	// We search from the last.
	// because there might be duplicate parameter keys,
	// and we want to get the last instance in this case.
	for i := len(param) - 1; i >= 0; i-- {
		kv := param[i]
		if kv.Key == key {
			return kv.Value
		}
	}
	return nil
}

// Len returns the number of items in the param.
func (param Parameters) Len() int {
	return len(param)
}

// Item is an item defined RFC 8941 Section 3.3. Items.
type Item struct {
	Value      Value
	Parameters Parameters
}

// InnerList is an array defined in RFC 8941 Section 3.1.1. Inner Lists.
type InnerList []Item

// List is an array defined in RFC 8941 Section 3.1. Lists.
type List []Item

// DictMember is a key-value pair of Dictionary.
type DictMember struct {
	Key  string
	Item Item
}

// Dictionary is an ordered map of key-value pairs defined in RFC 8941 Section 3.2. Dictionaries.
type Dictionary []DictMember

// Get returns the last item associated with the given key.
// If there are no items associated with the key, Get returns the zero value of Item.
func (dict Dictionary) Get(key string) Item {
	// In many cases, there are a few items.
	// So Linear searching is enough to handle them.

	// We search from the last.
	// because there might be duplicate parameter keys,
	// and we want to get the last instance in this case.
	for i := len(dict) - 1; i >= 0; i-- {
		kv := dict[i]
		if kv.Key == key {
			return kv.Item
		}
	}
	return Item{}
}

// Len returns the number of items in the dict.
func (dict Dictionary) Len() int {
	return len(dict)
}

func DecodeItem(fields ...string) (Item, error) {
	var item Item
	switch fields[0] {
	case "?0":
		item.Value = false
	case "?1":
		item.Value = true
	default:
		return Item{}, errors.New("TODO: implement")
	}
	return item, nil
}

func DecodeList(fields ...string) (List, error) {
	return nil, errors.New("TODO: implement")
}

func DecodeDictionary(fields ...string) (Dictionary, error) {
	return nil, errors.New("TODO: implement")
}
