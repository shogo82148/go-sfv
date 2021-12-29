package sfv

import "testing"

func TestDecodeNumber(t *testing.T) {
	runTestCases(t, "./testdata/number.json")
}

func TestDecodeString(t *testing.T) {
	runTestCases(t, "./testdata/string.json")
}

// func TestDecodeStringGenerated(t *testing.T) {
// 	runTestCases(t, "./testdata/string-generated.json")
// }

func TestDecodeToken(t *testing.T) {
	runTestCases(t, "./testdata/token.json")
}

func TestDecodeBinary(t *testing.T) {
	runTestCases(t, "./testdata/binary.json")
}

func TestDecodeBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
