package sfv

import (
	"encoding/base64"
	"fmt"
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
