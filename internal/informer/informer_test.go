package informer_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unconditionalday/server/internal/informer"
	"go.uber.org/zap"
)

type MockRunner struct {
	expectedCmd  string
	expectedArgs []string
	output       []byte
	err          error
}

func (m *MockRunner) Run(name string, args ...string) ([]byte, error) {
	return m.output, m.err
}

func (m *MockRunner) Version() string {
	return ""
}

type TestInput struct {
	text       string
	mockRunner *MockRunner
}

type TestExpect struct {
	expectedEmbeddings []float32
	err                error
}

func TestGetSimilarity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		input  TestInput
		output TestExpect
	}{
		{
			name: "valid input",
			input: TestInput{text: "example text", mockRunner: &MockRunner{
				expectedCmd:  "python",
				expectedArgs: []string{"path/to/your_script.py", "example text"},
				output:       []byte(`[0.1, 0.2, 0.3]`),
				err:          nil,
			}},
			output: TestExpect{
				expectedEmbeddings: []float32{0.1, 0.2, 0.3},
				err:                nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			z,_ := zap.NewDevelopment(nil)
			// Crea un informer con l'executorMock
			informerInstance := informer.NewInformer(tc.input.mockRunner, "",z)

			// Esegui la funzione da test
			embeddings, err := informerInstance.GetSimilarity(tc.input.text)

			// Verifica che l'errore restituito sia conforme all'output atteso
			if tc.output.err != nil {
				assert.Equal(t, tc.output.err.Error(), err.Error())
			} else {
				// Verifica che le embeddings restituite siano conformi all'output atteso
				expectedJSON, _ := json.Marshal(tc.output.expectedEmbeddings)
				actualJSON, _ := json.Marshal(embeddings)
				assert.Equal(t, string(expectedJSON), string(actualJSON))
			}

		})
	}
}
