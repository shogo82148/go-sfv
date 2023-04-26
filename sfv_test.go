package sfv

import "testing"

func TestExamples(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/examples.json")
}

func TestLargeGenerated(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/large-generated.json")
}

func TestList(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/list.json")
}

func TestParamList(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/param-list.json")
}

func TestListList(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/listlist.json")
}

func TestParamListList(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/param-listlist.json")
}

func TestDictionary(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/dictionary.json")
}

func TestParamDict(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/param-dict.json")
}

func TestKeyGenerated(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/key-generated.json")
}

func TestItem(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/item.json")
}

func TestNumber(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/number.json")
}

func TestNumberGenerated(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/number-generated.json")
}

func TestString(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/string.json")
}

func TestStringGenerated(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/string-generated.json")
}

func TestToken(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/token.json")
}

func TestTokenGenerated(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/token-generated.json")
}

func TestBinary(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/binary.json")
}

func TestBoolean(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/boolean.json")
}

func TestDate(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/date.json")
}
