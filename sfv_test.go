package sfv

import "testing"

func TestDecodeBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
