package clam

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/markbates/hepa"
)

type RunError struct {
	Err error `json:"err,omitempty"`

	Args   []string `json:"args,omitempty"`
	Dir    string   `json:"dir,omitempty"`
	Env    []string `json:"env,omitempty"`
	Exit   int      `json:"exit,omitempty"`
	Output []byte   `json:"output,omitempty"`
}

func (e RunError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"err":    e.Err,
		"args":   e.Args,
		"dir":    e.Dir,
		"exit":   e.Exit,
		"output": e.Output,
	}

	p := hepa.Deep()

	env := make([]string, 0, len(e.Env))
	for _, e := range e.Env {
		b, _ := p.Clean(strings.NewReader(e))
		env = append(env, string(b))
	}
	mm["env"] = env

	return json.MarshalIndent(mm, "", "  ")
}

func (e RunError) Error() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

func (e RunError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := e.Err.(Unwrapper); ok {
		return errors.Unwrap(e.Err)
	}

	return e.Err
}

func (e RunError) Is(target error) bool {
	if _, ok := target.(RunError); ok {
		return true
	}

	return errors.Is(e.Err, target)
}

func (e RunError) As(target any) bool {
	ex, ok := target.(*RunError)
	if !ok {
		return errors.As(e.Err, target)
	}

	(*ex) = e
	return true
}
