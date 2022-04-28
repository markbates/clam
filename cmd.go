package clam

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Env   []string  // environment variables, defaults to os.Environ()
	Dir   string    // working directory, defaults to os.Getwd()
	Stdin io.Reader // standard input, defaults to os.Stdin
}

func (c *Cmd) Run(ctx context.Context, args ...string) (*Result, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	if len(c.Dir) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		c.Dir = pwd
	}

	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}

	name := args[0]
	args = args[1:]

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = c.Dir
	cmd.Env = os.Environ()

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	cmd.Stdin = c.Stdin

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	if len(c.Env) > 0 {
		cmd.Env = c.Env
	}

	res := &Result{
		Args: cmd.Args,
		Dir:  cmd.Dir,
		Env:  cmd.Env,
		Exit: -1,
	}

	res.Err = cmd.Run()

	if cmd.ProcessState != nil {
		res.Exit = cmd.ProcessState.ExitCode()
	}

	res.Stderr = stderr.Bytes()
	res.Stdout = stdout.Bytes()

	return res, res.Err
}

type Result struct {
	Args   []string
	Dir    string
	Env    []string
	Err    error
	Exit   int
	Stderr []byte
	Stdout []byte
}

// CmdString returns a string representation of the command.
// 		$ echo hello
// 		$ go run maing.go
func (res Result) CmdString() string {
	return fmt.Sprintf("$ %s", strings.Join(res.Args, " "))
}