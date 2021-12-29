package sfv

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type encodeState struct {
	buf strings.Builder
}

func (s *encodeState) encodeItem(item Item) error {
	if err := s.encodeBareItem(item.Value); err != nil {
		return err
	}
	return nil
}

func (s *encodeState) encodeBareItem(v Value) error {
	switch v := v.(type) {
	case int64:
		if v > MaxInteger || v < MinInteger {
			return fmt.Errorf("integer %d is out of range", v)
		}
		var buf [20]byte
		dst := strconv.AppendInt(buf[:0], v, 10)
		s.buf.Write(dst)
	case float64:
		i := int64(math.RoundToEven(v * 1000))
		if i > MaxInteger || i < MinInteger {
			return fmt.Errorf("decimal %f is out of range", v)
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
			break // omit trailing zeros
		}
		s.buf.WriteByte(byte(frac/10) + '0')
		frac %= 10
		if frac == 0 {
			break // omit trailing zeros
		}
		s.buf.WriteByte(byte(frac) + '0')
	case string:
	case Token:
		if !v.Valid() {
			return fmt.Errorf("token %q is invalid form", v)
		}
		s.buf.WriteString(string(v))
	case []byte:
		s.buf.WriteByte(':')
		w := base64.NewEncoder(base64.StdEncoding, &s.buf)
		w.Write(v)
		w.Close()
		s.buf.WriteByte(':')
	case bool:
		if v {
			s.buf.WriteString("?1")
		} else {
			s.buf.WriteString("?0")
		}
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (s *encodeState) encodeList(list List) error {
	return nil
}

func (s *encodeState) encodeDictionary(dict Dictionary) error {
	return nil
}

func EncodeItem(item Item) (string, error) {
	state := &encodeState{}
	if err := state.encodeItem(item); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

func EncodeList(list List) (string, error) {
	state := &encodeState{}
	if err := state.encodeList(list); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}

func EncodeDictionary(dict Dictionary) (string, error) {
	state := &encodeState{}
	if err := state.encodeDictionary(dict); err != nil {
		return "", err
	}
	return state.buf.String(), nil
}
