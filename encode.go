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
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type encodeState struct {
	buf *bytes.Buffer
}

func (s *encodeState) encodeItem(item Item) error {
	if err := s.encodeBareItem(item.Value); err != nil {
		return err
	}
	if err := s.encodeParams(item.Parameters); err != nil {
		return err
	}
	return nil
}

func (s *encodeState) encodeInteger(v int64) error {
	if v > MaxInteger || v < MinInteger {
		return fmt.Errorf("sfv: integer %d is out of range", v)
	}
	var buf [20]byte
	dst := strconv.AppendInt(buf[:0], v, 10)
	s.buf.Write(dst)
	return nil
}

func (s *encodeState) encodeDecimal(v float64) error {
	i := int64(math.RoundToEven(v * 1000))
	if i > MaxInteger || i < MinInteger {
		return fmt.Errorf("sfv: decimal %f is out of range", v)
	}

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
		s.buf.WriteByte(':')
		w := base64.NewEncoder(base64.StdEncoding, s.buf)
		w.Write(v)
		w.Close()
		s.buf.WriteByte(':')

	case bool:
		if v {
			s.buf.WriteString("?1")
		} else {
			s.buf.WriteString("?0")
		}
	case time.Time:
		s.buf.WriteString("@")
		return s.encodeInteger(v.Unix())

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
	s.buf.WriteString(key)
	return nil
}

func (s *encodeState) encodeBareItemOrInnerList(value Value) error {
	if list, ok := value.(InnerList); ok {
		return s.encodeInnerList(list)
	}
	return s.encodeBareItem(value)
}

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

func EncodeItem(item Item) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()
	state := &encodeState{buf: buf}
	if err := state.encodeItem(item); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

func EncodeList(list List) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()
	state := &encodeState{buf: buf}
	if err := state.encodeList(list); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

func EncodeDictionary(dict Dictionary) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()
	state := &encodeState{buf: buf}
	if err := state.encodeDictionary(dict); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}
