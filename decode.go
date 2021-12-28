package sfv

import (
	"errors"
	"fmt"
)

const endOfInput = -1

type decodeState struct {
	fields     []string
	line, col  int
	endOfField bool

	err error
}

func (s *decodeState) peek() int {
	if s.line >= len(s.fields) {
		return endOfInput
	}
	if s.endOfField {
		// insert commas between fields.
		return ','
	}
	f := s.fields[s.line]
	if s.col >= len(f) {
		return endOfInput
	}
	return int(f[s.col])
}

func (s *decodeState) next() {
	if s.line >= len(s.fields) {
		// no more inputs.
		return
	}
	if s.endOfField {
		s.endOfField = false
		return
	}

	f := s.fields[s.line]
	s.col++
	if s.col >= len(f) {
		// goto next the field.
		s.col = 0
		s.line++
		s.endOfField = true
	}
}

func (s *decodeState) decodeItem() (Item, error) {
	v, err := s.decodeBareItem()
	if err != nil {
		return Item{}, err
	}

	// TODO: parse parameters

	return Item{
		Value: v,
	}, nil
}

func (s *decodeState) decodeBareItem() (Value, error) {
	ch := s.peek()
	switch {
	case ch == '-' || (ch >= '0' && ch <= '9'):
		// an Integer or Decimal
	case ch == '"':
		// a String
	case ch == '*' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
		// a Token
	case ch == ':':
		// a Byte Sequence
	case ch == '?':
		// a Boolean
		s.next()
		switch s.peek() {
		case '0':
			return false, nil
		case '1':
			return true, nil
		}
	}
	ch = s.peek()
	if ch == endOfInput {
		return nil, errors.New("unexpected the end of the input")
	}
	return nil, fmt.Errorf("unexpected character: %q", ch)
}

func DecodeItem(fields []string) (Item, error) {
	state := &decodeState{
		fields: fields,
	}
	return state.decodeItem()
}

func DecodeList(fields []string) (List, error) {
	return nil, errors.New("TODO: implement")
}

func DecodeDictionary(fields []string) (Dictionary, error) {
	return nil, errors.New("TODO: implement")
}
