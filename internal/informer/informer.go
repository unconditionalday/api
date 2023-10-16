package informer

import (
	"encoding/json"
	"fmt"

	"github.com/unconditionalday/server/internal/x/exec"
	"go.uber.org/zap"
)

type Informer struct {
	scriptsPath   string
	similarityCmd exec.Runner
	logger        *zap.Logger
}

func NewInformer(Runner exec.Runner, scriptsPath string, logger *zap.Logger) *Informer {
	return &Informer{
		scriptsPath:   scriptsPath,
		similarityCmd: Runner,
		logger:        logger,
	}
}

func (i *Informer) GetSimilarity(text string) ([]float32, error) {
	p := fmt.Sprintf("%s/similarity.py", i.scriptsPath)

	o, err := i.similarityCmd.Run(p, text)
	if err != nil {
		return nil, err
	}

	var embeddings []float32
	if err := json.Unmarshal(o, &embeddings); err != nil {
		return nil, err
	}

	return embeddings, nil
}
