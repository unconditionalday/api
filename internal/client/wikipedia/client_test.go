package wikipedia_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unconditionalday/server/internal/client/wikipedia"
)

type TestInput struct {
	query string
	lang  string
}

type TestExpect struct {
	validRes bool
	err      error
}

func TestFetchContextDetails(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		input  TestInput
		output TestExpect
	}{
		{
			name:  "empty query string",
			input: TestInput{query: "", lang: "en"},
			output: TestExpect{
				validRes: false,
				err:      errors.New("query string must not be empty"),
			},
		},
		{
			name:  "empty language",
			input: TestInput{query: "Lorem ipsum", lang: ""},
			output: TestExpect{
				validRes: false,
				err:      errors.New("language string must not be empty"),
			},
		},
		{
			name:  "valid query",
			input: TestInput{query: "Italy", lang: "en"},
			output: TestExpect{
				validRes: true,
				err:      nil,
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := wikipedia.NewClient()
			actual, err := w.FetchContextDetails(tc.input.query, tc.input.lang)

			if tc.output.err != nil {
				assert.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.output.validRes, actual.IsValid())
			}
		})
	}
}
