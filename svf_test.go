package sfv

import "testing"

func TestExamples(t *testing.T) {
	runTestCases(t, "./testdata/examples.json")
}

func TestLargeGenerated(t *testing.T) {
	runTestCases(t, "./testdata/large-generated.json")
}

func TestList(t *testing.T) {
	runTestCases(t, "./testdata/list.json")
}

func TestParamList(t *testing.T) {
	runTestCases(t, "./testdata/param-list.json")
}

func TestListList(t *testing.T) {
	runTestCases(t, "./testdata/listlist.json")
}

func TestParamListList(t *testing.T) {
	runTestCases(t, "./testdata/param-listlist.json")
}

func TestDictionary(t *testing.T) {
	runTestCases(t, "./testdata/dictionary.json")
}

func TestParamDict(t *testing.T) {
	runTestCases(t, "./testdata/param-dict.json")
}

func TestKeyGenerated(t *testing.T) {
	runTestCases(t, "./testdata/key-generated.json")
}

func TestItem(t *testing.T) {
	runTestCases(t, "./testdata/item.json")
}

func TestNumber(t *testing.T) {
	runTestCases(t, "./testdata/number.json")
}

func TestNumberGenerated(t *testing.T) {
	runTestCases(t, "./testdata/number-generated.json")
}

func TestString(t *testing.T) {
	runTestCases(t, "./testdata/string.json")
}

func TestStringGenerated(t *testing.T) {
	runTestCases(t, "./testdata/string-generated.json")
}

func TestToken(t *testing.T) {
	runTestCases(t, "./testdata/token.json")
}

func TestTokenGenerated(t *testing.T) {
	runTestCases(t, "./testdata/token-generated.json")
}

func TestBinary(t *testing.T) {
	runTestCases(t, "./testdata/binary.json")
}

func TestBoolean(t *testing.T) {
	runTestCases(t, "./testdata/boolean.json")
}
