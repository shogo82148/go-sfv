package sfv

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

const endOfInput = -1

var validKeyChars = [256]bool{
	'_': true,
	'-': true,
	'.': true,
	'*': true,

	// DIGIT from RFC 7230
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,

	// lcalpha
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,
}

// valid character for Token except the first character.
var validTokenChars = [256]bool{
	':': true,
	'/': true,

	// tchar from RFC 7230
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'|':  true,
	'~':  true,
	// and DIGIT and ALPHA

	// DIGIT from RFC 7230
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,

	// ALPHA from RFC 7230
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
	'G': true,
	'H': true,
	'I': true,
	'J': true,
	'K': true,
	'L': true,
	'M': true,
	'N': true,
	'O': true,
	'P': true,
	'Q': true,
	'R': true,
	'S': true,
	'T': true,
	'U': true,
	'V': true,
	'W': true,
	'X': true,
	'Y': true,
	'Z': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,
}

// character for base64-decoding
var validBase64Chars = [256]bool{
	'+': true,
	'/': true,
	'=': true,

	// DIGIT from RFC 7230
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,

	// ALPHA from RFC 7230
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
	'G': true,
	'H': true,
	'I': true,
	'J': true,
	'K': true,
	'L': true,
	'M': true,
	'N': true,
	'O': true,
	'P': true,
	'Q': true,
	'R': true,
	'S': true,
	'T': true,
	'U': true,
	'V': true,
	'W': true,
	'X': true,
	'Y': true,
	'Z': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,
}

func isDigit(ch int) bool {
	return ch >= '0' && ch <= '9'
}

type decodeState struct {
	fields     []string
	line, col  int
	endOfField bool
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

func (s *decodeState) skipSPs() {
	for s.peek() == ' ' {
		s.next()
	}
}

// skip OWS in RFC 7230
func (s *decodeState) skipOWS() {
	for ch := s.peek(); ch == ' ' || ch == '\t'; ch = s.peek() {
		s.next()
	}
}

func (s *decodeState) errUnexpectedCharacter() error {
	ch := s.peek()
	if ch == endOfInput {
		return errors.New("sfv: unexpected the end of the input")
	}
	return fmt.Errorf("sfv: unexpected character: %q", ch)
}

func (s *decodeState) decodeItem() (Item, error) {
	v, err := s.decodeBareItem()
	if err != nil {
		return Item{}, err
	}

	param, err := s.decodeParameters()
	if err != nil {
		return Item{}, err
	}

	return Item{
		Value:      v,
		Parameters: param,
	}, nil
}

func (s *decodeState) decodeBareItem() (Value, error) {
	ch := s.peek()
	switch {
	case ch == '-' || isDigit(ch):
		// an Integer or Decimal
		neg := false
		if ch == '-' {
			neg = true
			s.next()

			if !isDigit(s.peek()) {
				return nil, s.errUnexpectedCharacter()
			}
		}

		num := int64(0)
		cnt := 0
		for {
			ch := s.peek()
			if !isDigit(ch) {
				break
			}
			s.next()
			num = num*10 + int64(ch-'0')
			cnt++
			if cnt > 15 {
				return nil, errors.New("sfv: integer overflow")
			}
		}
		if s.peek() != '.' {
			// it is an Integer
			if neg {
				num *= -1
			}
			return num, nil
		}
		// current character is '.'
		s.next() // skip '.'

		// it might be a Decimal
		if cnt > 12 {
			return nil, errors.New("sfv: decimal overflow")
		}

		frac := 0
		ch := s.peek()
		if !isDigit(ch) {
			// fractional part MUST NOT be empty.
			return nil, s.errUnexpectedCharacter()
		}
		s.next()
		frac = frac*10 + int(ch-'0')

		ch = s.peek()
		if !isDigit(ch) {
			ret := float64(num) + float64(frac)/10
			if neg {
				ret *= -1
			}
			return ret, nil
		}
		s.next()
		frac = frac*10 + int(ch-'0')

		ch = s.peek()
		if !isDigit(ch) {
			ret := float64(num) + float64(frac)/100
			if neg {
				ret *= -1
			}
			return ret, nil
		}
		s.next()
		frac = frac*10 + int(ch-'0')

		ch = s.peek()
		if !isDigit(ch) {
			ret := float64(num) + float64(frac)/1000
			if neg {
				ret *= -1
			}
			return ret, nil
		}
		return nil, errors.New("sfv: decimal has too long fractional part")

	case ch == '"':
		// a String
		s.next() // skip '"'
		var buf strings.Builder
		for {
			ch := s.peek()
			switch {
			case ch == '\\':
				s.next() // skip '\\'
				switch s.peek() {
				case '\\':
					s.next() // skip '\\'
					buf.WriteByte('\\')
				case '"':
					s.next() // skip '"'
					buf.WriteByte('"')
				default:
					return nil, s.errUnexpectedCharacter()
				}
			case ch == '"':
				// the end of a String
				s.next() // skip '"'
				return buf.String(), nil
			case ch >= 0x20 && ch < 0x7f:
				s.next()
				buf.WriteByte(byte(ch))
			default:
				return nil, s.errUnexpectedCharacter()
			}
		}

	case ch == '*' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
		// a Token
		var buf strings.Builder
		for {
			ch := s.peek()
			switch {
			case ch == endOfInput:
				return Token(buf.String()), nil
			case validTokenChars[ch]:
				s.next()
				buf.WriteByte(byte(ch))
			default:
				return Token(buf.String()), nil
			}
		}

	case ch == ':':
		// a Byte Sequence
		s.next() // skip ':'
		var buf bytes.Buffer
		for {
			ch := s.peek()
			switch {
			case ch == endOfInput:
				return nil, s.errUnexpectedCharacter()
			case ch == ':':
				// the end of a Binary
				s.next() // skip ':'

				// add missing "=" padding
				// RFC 8941 says that parsers SHOULD NOT fail when "=" padding is not present.
				switch buf.Len() % 4 {
				case 0:
				case 1:
					buf.WriteByte('=')
					fallthrough
				case 2:
					buf.WriteByte('=')
					fallthrough
				case 3:
					buf.WriteByte('=')
				}

				enc := base64.StdEncoding
				ret := make([]byte, enc.DecodedLen(buf.Len()))
				n, err := enc.Decode(ret, buf.Bytes())
				if err != nil {
					return nil, err
				}
				return ret[:n], nil
			case validBase64Chars[ch]:
				s.next()
				buf.WriteByte(byte(ch))
			default:
				return nil, s.errUnexpectedCharacter()
			}
		}

	case ch == '?':
		// a Boolean
		s.next() // skip '?'
		switch s.peek() {
		case '0':
			s.next() // skip '0'
			return false, nil
		case '1':
			s.next() // skip '1'
			return true, nil
		default:
			return nil, s.errUnexpectedCharacter()
		}
	case ch == '@':
		// a Date
		s.next() // skip '@'

		// check sign
		neg := false
		if ch := s.peek(); ch == '-' {
			neg = true
			s.next() // skip '-'
		}

		if !isDigit(s.peek()) {
			return nil, s.errUnexpectedCharacter()
		}

		num := int64(0)
		cnt := 0
		for {
			ch := s.peek()
			if !isDigit(ch) {
				break
			}
			s.next()
			num = num*10 + int64(ch-'0')
			cnt++
			if cnt > 15 {
				return nil, errors.New("sfv: integer overflow")
			}
		}

		if s.peek() == '.' {
			// a Date must not a Decimal.
			return nil, s.errUnexpectedCharacter()
		}
		if neg {
			num *= -1
		}
		return time.Unix(num, 0), nil
	}
	return nil, s.errUnexpectedCharacter()
}

func (s *decodeState) decodeParameters() (Parameters, error) {
	var params Parameters
	seenKeys := map[string]int{}
	for {
		if s.peek() != ';' {
			break
		}
		s.next() // skip ';'
		s.skipSPs()

		key, err := s.decodeKey()
		if err != nil {
			return nil, err
		}
		var value Value
		if s.peek() == '=' {
			s.next() // skip '='
			value, err = s.decodeBareItem()
			if err != nil {
				return nil, err
			}
		} else {
			value = true
		}
		if i, ok := seenKeys[key]; ok {
			// parameters already contains a key,
			// overwrite its value
			params[i] = Parameter{
				Key:   key,
				Value: value,
			}
		} else {
			seenKeys[key] = len(params)
			params = append(params, Parameter{
				Key:   key,
				Value: value,
			})
		}
	}

	return params, nil
}

func (s *decodeState) decodeKey() (string, error) {
	ch := s.peek()
	if (ch < 'a' || ch > 'z') && ch != '*' {
		return "", s.errUnexpectedCharacter()
	}
	s.next()

	var buf strings.Builder
	buf.WriteByte(byte(ch))
	for {
		ch := s.peek()
		if ch == endOfInput {
			break
		}
		if !validKeyChars[ch] {
			break
		}
		s.next()
		buf.WriteByte(byte(ch))
	}
	return buf.String(), nil
}

func (s *decodeState) decodeItemOrInnerItem() (Item, error) {
	if s.peek() != '(' {
		// It might be an Item
		return s.decodeItem()
	}
	s.next() // skip '('

	// parse as an Inner List
	list := InnerList{}
	for {
		s.skipSPs()
		ch := s.peek()
		if ch == ')' {
			s.next() // skip ')'
			break
		}

		item, err := s.decodeItem()
		if err != nil {
			return Item{}, err
		}
		list = append(list, item)
		ch = s.peek()
		if ch != ' ' && ch != ')' {
			return Item{}, s.errUnexpectedCharacter()
		}
	}
	params, err := s.decodeParameters()
	if err != nil {
		return Item{}, err
	}

	return Item{
		Value:      list,
		Parameters: params,
	}, nil
}

func (s *decodeState) decodeList() (List, error) {
	var list List

	if s.peek() == endOfInput {
		// it is an empty list
		return nil, nil
	}

	for {
		item, err := s.decodeItemOrInnerItem()
		if err != nil {
			return nil, err
		}
		list = append(list, item)

		s.skipOWS()
		ch := s.peek()
		if ch == endOfInput {
			break
		}
		if ch != ',' {
			return nil, s.errUnexpectedCharacter()
		}
		s.next() // skip ','
		s.skipOWS()
		if s.peek() == endOfInput {
			// it is trailing comma.
			return nil, errors.New("sfv: trailing comma is not allowed")
		}
	}
	return list, nil
}

func (s *decodeState) decodeDictionary() (Dictionary, error) {
	if s.peek() == endOfInput {
		// it is an empty dictionary
		return nil, nil
	}

	var dict Dictionary
	seenKeys := map[string]int{}
	for {
		// decode keys
		key, err := s.decodeKey()
		if err != nil {
			return nil, err
		}

		// decode items
		var item Item
		if s.peek() == '=' {
			s.next() // skip '='
			item, err = s.decodeItemOrInnerItem()
			if err != nil {
				return nil, err
			}
		} else {
			params, err := s.decodeParameters()
			if err != nil {
				return nil, err
			}
			item = Item{
				Value:      true,
				Parameters: params,
			}
		}
		if i, ok := seenKeys[key]; ok {
			// parameters already contains a key,
			// overwrite its value
			dict[i] = DictMember{
				Key:  key,
				Item: item,
			}
		} else {
			seenKeys[key] = len(dict)
			dict = append(dict, DictMember{
				Key:  key,
				Item: item,
			})
		}

		// skip commas
		s.skipOWS()
		ch := s.peek()
		if ch == endOfInput {
			break
		}
		if ch != ',' {
			return nil, s.errUnexpectedCharacter()
		}
		s.next() // skip ','
		s.skipOWS()
		if s.peek() == endOfInput {
			// it is trailing comma.
			return nil, errors.New("sfv: trailing comma is not allowed")
		}
	}
	return dict, nil
}

func DecodeItem(fields []string) (Item, error) {
	state := &decodeState{
		fields: fields,
	}
	state.skipSPs()
	ret, err := state.decodeItem()
	if err != nil {
		return Item{}, err
	}
	state.skipSPs()
	if state.peek() != endOfInput {
		return Item{}, state.errUnexpectedCharacter()
	}
	return ret, nil
}

func DecodeList(fields []string) (List, error) {
	state := &decodeState{
		fields: fields,
	}
	state.skipSPs()
	ret, err := state.decodeList()
	if err != nil {
		return nil, err
	}
	state.skipSPs()
	if state.peek() != endOfInput {
		return nil, state.errUnexpectedCharacter()
	}
	return ret, nil
}

func DecodeDictionary(fields []string) (Dictionary, error) {
	state := &decodeState{
		fields: fields,
	}
	state.skipSPs()
	ret, err := state.decodeDictionary()
	if err != nil {
		return nil, err
	}
	state.skipSPs()
	if state.peek() != endOfInput {
		return nil, state.errUnexpectedCharacter()
	}
	return ret, nil
}
