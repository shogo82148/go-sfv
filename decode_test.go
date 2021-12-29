package sfv

import "testing"

func TestDecodeExamples(t *testing.T) {
	runTestCases(t, "./testdata/examples.json")
}

func TestDecodeLargeGenerated(t *testing.T) {
	runTestCases(t, "./testdata/large-generated.json")
}

func TestDecodeList(t *testing.T) {
	runTestCases(t, "./testdata/list.json")
}

func TestDecodeParamList(t *testing.T) {
	runTestCases(t, "./testdata/param-list.json")
}

func TestDecodeListList(t *testing.T) {
	runTestCases(t, "./testdata/listlist.json")
}

func TestDecodeParamListList(t *testing.T) {
	runTestCases(t, "./testdata/param-listlist.json")
}

func TestDecodeDictionary(t *testing.T) {
	runTestCases(t, "./testdata/dictionary.json")
}

func TestDecodeParamDict(t *testing.T) {
	runTestCases(t, "./testdata/param-dict.json")
}

func TestDecodeKeyGenerated(t *testing.T) {
	runTestCases(t, "./testdata/key-generated.json")
}

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
