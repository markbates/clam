package clam

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Result struct {
	Args   []string `json:"args,omitempty"`
	Dir    string   `json:"dir,omitempty"`
	Env    []string `json:"-"`
	Err    error    `json:"err,omitempty"`
	Exit   int      `json:"exit,omitempty"`
	Stderr []byte   `json:"stderr,omitempty"`
	Stdout []byte   `json:"stdout,omitempty"`
}

func (res Result) String() string {
	b, _ := json.Marshal(res)
	return string(b)
}

// CmdString returns a string representation of the command.
// 		$ echo hello
// 		$ go run maing.go
func (res Result) CmdString() string {
	return fmt.Sprintf("$ %s", strings.Join(res.Args, " "))
}
