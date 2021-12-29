package sfv

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const endOfInput = -1

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

func (s *decodeState) skipSPs() {
	for s.peek() == ' ' {
		s.next()
	}
}

func (s *decodeState) errUnexpectedCharacter() error {
	ch := s.peek()
	if ch == endOfInput {
		return errors.New("unexpected the end of the input")
	}
	return fmt.Errorf("unexpected character: %q", ch)
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
				return nil, errors.New("integer overflow")
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
			return nil, errors.New("decimal overflow")
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
		return nil, errors.New("decimal has too long fractional part")

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
		}
	}
	return nil, s.errUnexpectedCharacter()
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
	return nil, errors.New("TODO: implement")
}

func DecodeDictionary(fields []string) (Dictionary, error) {
	return nil, errors.New("TODO: implement")
}
