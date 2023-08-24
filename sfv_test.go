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

func TestDisplayString(t *testing.T) {
	runTestCases(t, "./testdata/structured-field-tests/display-string.json")
}

func TestExtra(t *testing.T) {
	// This test is not part of the structured-field-tests suite.
	runTestCases(t, "./testdata/extra.json")
}

func TestToken_Valid(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"*", true},
		{"a", true},
		{"z", true},
		{"A", true},
		{"Z", true},
		{"0", false},
		{":", false},
		{"/", false},
		{"!", false},
		{"\"", false},
		{"#", false},
		{"$", false},
		{"%", false},
		{"&", false},
		{"'", false},
		{"(", false},
		{")", false},
		{"aa", true},
		{"aA", true},
		{"a0", true},
		{"a:", true},
		{"a/", true},
	}
	for _, c := range cases {
		if got := Token(c.in).Valid(); got != c.want {
			t.Errorf("Token(%q).Valid() = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestIsValidString(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"*", true},
		{"a", true},
		{"z", true},
		{"A", true},
		{"Z", true},
		{"0", true},
		{"9", true},
		{"~", true},
		{"aa", true},
		{" ", true},
		{"\t", false},
		{"\a", false},
		{"\n", false},
		{"\r", false},
		{"\f", false},
		{"\u0080", false},
	}
	for _, c := range cases {
		if got := IsValidString(c.in); got != c.want {
			t.Errorf("isValidString(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
