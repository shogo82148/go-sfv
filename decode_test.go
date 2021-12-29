package sfv

import "testing"

func TestDecodeItem(t *testing.T) {
	runTestCases(t, "./testdata/item.json")
}

func TestDecodeNumber(t *testing.T) {
	runTestCases(t, "./testdata/number.json")
}

func TestDecodeNumberGenerated(t *testing.T) {
	runTestCases(t, "./testdata/number-generated.json")
}

func TestDecodeString(t *testing.T) {
	runTestCases(t, "./testdata/string.json")
}

func TestDecodeStringGenerated(t *testing.T) {
	runTestCases(t, "./testdata/string-generated.json")
}

func TestDecodeToken(t *testing.T) {
	runTestCases(t, "./testdata/token.json")
}

func TestDecodeTokenGenerated(t *testing.T) {
	runTestCases(t, "./testdata/token-generated.json")
}

func TestDecodeBinary(t *testing.T) {
	runTestCases(t, "./testdata/binary.json")
}

func TestDecodeBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
