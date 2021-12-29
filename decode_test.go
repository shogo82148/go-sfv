package sfv

import "testing"

func TestDecodeString(t *testing.T) {
	runTestCases(t, "./testdata/string.json")
}

// func TestDecodeStringGenerated(t *testing.T) {
// 	runTestCases(t, "./testdata/string-generated.json")
// }

func TestDecodeToken(t *testing.T) {
	runTestCases(t, "./testdata/token.json")
}

func TestDecodeBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
