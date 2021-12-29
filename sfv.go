package sfv

// Token is a token defined in RFC 8941 Section 3.3.4. Tokens.
type Token string

// Valid returns whether the t has valid form.
func (t Token) Valid() bool {
	if t == "" {
		return false
	}
	ch := t[0]
	if ch != '*' && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
		return false
	}
	for _, ch := range []byte(t[1:]) {
		if !validTokenChars[ch] {
			return false
		}
	}
	return true
}

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
	for _, kv := range param {
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

	for _, kv := range dict {
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
