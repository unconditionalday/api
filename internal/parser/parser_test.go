package parser_test

import (
	"testing"

	"github.com/unconditionalday/server/internal/parser"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "single word with new line",
			input:    "hello\n",
			expected: "hello",
		},
		{
			name:     "single word with html tags",
			input:    "<p>hello</p>",
			expected: "hello",
		},
		{
			name:     "single word with html tags and new line",
			input:    "<p>hello</p>\n",
			expected: "hello",
		},
		{
			name:     "single word with html tags and new line and bloated text",
			input:    "<p>hello</p>\n[Continue reading...]",
			expected: "hello",
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := parser.NewParser()
			actual := p.Parse(tc.input)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
