package exec

import (
	"errors"
	"os/exec"
	"strings"
)

var (
	ErrInterpreterNotFound = errors.New("python interpreter not found in PATH")
)

type Runner interface {
	Run(name string, args ...string) ([]byte, error)
	Version() string
}

type PythonRunner struct {
	path    string
	version string
}

func NewPythonRunner() (*PythonRunner, error) {
	pythonPath, err := exec.LookPath("python")
	if err != nil {
		return nil, ErrInterpreterNotFound
	}

	v := exec.Command(pythonPath, "--version").String()
	version := strings.Split(v, " ")[0]

	return &PythonRunner{
		path:    pythonPath,
		version: version,
	}, nil
}

func (p *PythonRunner) Run(name string, args ...string) ([]byte, error) {
	args = append([]string{name}, args...)

	o, err := exec.Command(p.path, args...).CombinedOutput()
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (p *PythonRunner) Version() string {
	return p.version
}
