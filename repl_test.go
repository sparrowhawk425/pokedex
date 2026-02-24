package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "                              hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  Madam I'm  adaM  ",
			expected: []string{"madam", "i'm", "adam"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual slice %d doesn't match expected length %d", len(actual), len(c.expected))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Actual %v doesn't match expected %v", word, expectedWord)
				t.Fail()
			}
		}
	}
}
