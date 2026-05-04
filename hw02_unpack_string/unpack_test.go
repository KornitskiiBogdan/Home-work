package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "a1b1c1", expected: "abc"},
		{input: "a3", expected: "aaa"},
		{input: "a0b0c0", expected: ""},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "🙃5", expected: "🙃🙃🙃🙃🙃"},
		{input: "世3界1", expected: "世世世界"},
		{input: "aaф0b", expected: "aab"},
	}

	for _, tc := range tests {
		testcase := tc
		t.Run(testcase.input, func(t *testing.T) {
			result, err := Unpack(testcase.input)
			require.NoError(t, err)
			require.Equal(t, testcase.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "5", "a05", "123"}
	for _, tc := range invalidStrings {
		testcase := tc
		t.Run(testcase, func(t *testing.T) {
			_, err := Unpack(testcase)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
