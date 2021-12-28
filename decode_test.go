package sfv

import "testing"

func TestDecodeString(t *testing.T) {
	runTestCases(t, "./testdata/string.json")
}

// func TestDecodeStringGenerated(t *testing.T) {
// 	runTestCases(t, "./testdata/string-generated.json")
// }

func TestDecodeBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
