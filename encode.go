package sfv

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

var hexBytes = "0123456789abcdef"

// bytes.Buffer provides AvailableBuffer from Go 1.21.0.
type availableBuffer interface {
	AvailableBuffer() []byte
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(encodeState)
	},
}

func getEncodeState() *encodeState {
	s := bufPool.Get().(*encodeState)
	s.buf.Reset()
	return s
}

func putEncodeState(s *encodeState) {
	bufPool.Put(s)
}

type encodeState struct {
	buf bytes.Buffer
}

// encodeItem serializes an item according to RFC 8941 Section 4.1.3.
func (s *encodeState) encodeItem(item Item) error {
	if err := s.encodeBareItem(item.Value); err != nil {
		return err
	}
	if err := s.encodeParams(item.Parameters); err != nil {
		return err
	}
	return nil
}

// encodeInteger serializes an integer according to RFC 8941 Section 4.1.4.
func (s *encodeState) encodeInteger(v int64) error {
	if v > MaxInteger || v < MinInteger {
		return fmt.Errorf("sfv: integer %d is out of range", v)
	}
	var buf [20]byte
	dst := strconv.AppendInt(buf[:0], v, 10)
	s.buf.Write(dst)
	return nil
}

// encodeDecimal serializes an decimal according to RFC 8941 Section 4.1.5.
func (s *encodeState) encodeDecimal(v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("sfv: decimal %f is not a finite number", v)
	}
	if v > MaxDecimal || v < MinDecimal {
		return fmt.Errorf("sfv: decimal %f is out of range", v)
	}
	i := int64(math.RoundToEven(v * 1000))

	// write the sign
	if i < 0 {
		s.buf.WriteByte('-')
		i *= -1
	}

	// integer component
	var buf [20]byte
	dst := strconv.AppendInt(buf[:0], i/1000, 10)
	s.buf.Write(dst)

	// fractional component
	frac := i % 1000
	s.buf.WriteByte('.')
	s.buf.WriteByte(byte(frac/100) + '0')
	frac %= 100
	if frac == 0 {
		return nil // omit trailing zeros
	}
	s.buf.WriteByte(byte(frac/10) + '0')
	frac %= 10
	if frac == 0 {
		return nil // omit trailing zeros
	}
	s.buf.WriteByte(byte(frac) + '0')

	return nil
}

// encodeBinary serializes a byte sequence according to RFC 8941 Section 4.1.8.
func (s *encodeState) encodeByteSequence(v []byte) error {
	// allocate a buffer
	l := base64.StdEncoding.EncodedLen(len(v)) + 2
	var b []byte
	if ab, ok := any(&s.buf).(availableBuffer); ok {
		s.buf.Grow(l)
		b = ab.AvailableBuffer()
	}
	if cap(b) < l {
		b = make([]byte, l)
	} else {
		b = b[:l]
	}

	// encode the byte sequence as base64.
	b[0] = ':'
	b[l-1] = ':'
	base64.StdEncoding.Encode(b[1:], v)
	s.buf.Write(b[:l])
	return nil
}

// encodeDisplayString serializes a display string according to [sfbis-03 Section 4.1.11].
//
// [sfbis-03 Section 4.1.11]: https://www.ietf.org/archive/id/draft-ietf-httpbis-sfbis-03.html#name-serializing-a-display-strin
func (s *encodeState) encodeDisplayString(v string) error {
	if !utf8.ValidString(v) {
		return fmt.Errorf("sfv: display string %q has invalid characters", v)
	}
	s.buf.WriteByte('"')
	for _, ch := range []byte(v) {
		if ch == '%' || ch == '"' || ch <= 0x1f || ch >= 0x7f {
			s.buf.WriteByte('%')
			s.buf.WriteByte(hexBytes[ch>>4])
			s.buf.WriteByte(hexBytes[ch&0xf])
		} else {
			s.buf.WriteByte(ch)
		}
	}
	s.buf.WriteByte('"')
	return nil
}

// encodeBareItem serializes a bare item according to RFC 8941 Section 4.1.3.1.
func (s *encodeState) encodeBareItem(v Value) error {
	switch v := v.(type) {
	case int8:
		return s.encodeInteger(int64(v))
	case uint8:
		return s.encodeInteger(int64(v))
	case int16:
		return s.encodeInteger(int64(v))
	case uint16:
		return s.encodeInteger(int64(v))
	case int32:
		return s.encodeInteger(int64(v))
	case uint32:
		return s.encodeInteger(int64(v))
	case int:
		return s.encodeInteger(int64(v))
	case uint:
		w := int(v) // this cast may overflow,
		if w < 0 {  // so we need to check it.
			return fmt.Errorf("sfv: integer %d is out of range", v)
		}
		return s.encodeInteger(int64(v))
	case int64:
		return s.encodeInteger(v)
	case uint64:
		w := int64(v) // this cast may overflow,
		if w < 0 {    // so we need to check it.
			return fmt.Errorf("sfv: integer %d is out of range", v)
		}
		return s.encodeInteger(int64(v))

	case float64:
		return s.encodeDecimal(v)
	case float32:
		return s.encodeDecimal(float64(v))

	case string:
		if !IsValidString(v) {
			return fmt.Errorf("sfv: string %q has invalid characters", v)
		}
		s.buf.WriteByte('"')
		for _, ch := range []byte(v) {
			switch ch {
			case '\\':
				s.buf.WriteString("\\\\")
			case '"':
				s.buf.WriteString("\\\"")
			default:
				s.buf.WriteByte(ch)
			}
		}
		s.buf.WriteByte('"')

	case Token:
		if !v.Valid() {
			return fmt.Errorf("sfv: token %q has invalid characters", v)
		}
		s.buf.WriteString(string(v))

	case []byte:
		return s.encodeByteSequence(v)

	case bool:
		if v {
			s.buf.WriteString("?1")
		} else {
			s.buf.WriteString("?0")
		}

	case time.Time:
		s.buf.WriteByte('@')
		return s.encodeInteger(v.Unix())

	case DisplayString:
		s.buf.WriteByte('%')
		return s.encodeDisplayString(string(v))

	default:
		return fmt.Errorf("sfv: unsupported type: %T", v)
	}
	return nil
}

func (s *encodeState) encodeParams(params Parameters) error {
	for _, param := range params {
		s.buf.WriteByte(';')
		if err := s.encodeKey(param.Key); err != nil {
			return err
		}
		if param.Value == true {
			continue
		}
		s.buf.WriteByte('=')
		if err := s.encodeBareItem(param.Value); err != nil {
			return err
		}
	}
	return nil
}

func (s *encodeState) encodeKey(key string) error {
	// validation
	if len(key) == 0 {
		return errors.New("sfv: key is an empty string")
	}
	if (key[0] < 'a' || key[0] > 'z') && key[0] != '*' {
		return fmt.Errorf("sfv: key %q has invalid characters", key)
	}
	for _, ch := range []byte(key[1:]) {
		if !validKeyChars[ch] {
			return fmt.Errorf("sfv: key %q has invalid characters", key)
		}
	}

	// encode the key
	s.buf.WriteString(key)
	return nil
}

func (s *encodeState) encodeBareItemOrInnerList(value Value) error {
	if list, ok := value.(InnerList); ok {
		return s.encodeInnerList(list)
	}
	return s.encodeBareItem(value)
}

// encodeInnerList serializes an inner list according to RFC 8941 Section 4.1.1.1.
func (s *encodeState) encodeInnerList(list InnerList) error {
	s.buf.WriteByte('(')
	for i, item := range list {
		if err := s.encodeItem(item); err != nil {
			return err
		}
		if i+1 < len(list) {
			s.buf.WriteRune(' ')
		}
	}
	s.buf.WriteByte(')')
	return nil
}

// encodeList serializes a list according to RFC 8941 Section 4.1.1.
func (s *encodeState) encodeList(list List) error {
	for i, item := range list {
		if err := s.encodeBareItemOrInnerList(item.Value); err != nil {
			return err
		}
		if err := s.encodeParams(item.Parameters); err != nil {
			return err
		}
		if i+1 < len(list) {
			s.buf.WriteString(", ")
		}
	}
	return nil
}

// encodeDictionary serializes a dictionary according to RFC 8941 Section 4.1.2.
func (s *encodeState) encodeDictionary(dict Dictionary) error {
	for i, item := range dict {
		if err := s.encodeKey(item.Key); err != nil {
			return err
		}
		if item.Item.Value != true {
			s.buf.WriteByte('=')
			if err := s.encodeBareItemOrInnerList(item.Item.Value); err != nil {
				return err
			}
		}
		if err := s.encodeParams(item.Item.Parameters); err != nil {
			return err
		}
		if i+1 < len(dict) {
			s.buf.WriteString(", ")
		}
	}
	return nil
}

// EncodeItem encodes the given item to Structured Field Values.
func EncodeItem(item Item) (string, error) {
	state := getEncodeState()
	defer putEncodeState(state)

	if err := state.encodeItem(item); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

// EncodeList encodes the given list to Structured Field Values.
func EncodeList(list List) (string, error) {
	state := getEncodeState()
	defer putEncodeState(state)

	if err := state.encodeList(list); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

// EncodeDictionary encodes the given dictionary to Structured Field Values.
func EncodeDictionary(dict Dictionary) (string, error) {
	state := getEncodeState()
	defer putEncodeState(state)

	if err := state.encodeDictionary(dict); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}
