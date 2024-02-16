package clam

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
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

	c.Env = append(c.Env, "GOWORK=off")
	c.Env = append(c.Env, os.Environ()...)

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

	now := time.Now()
	res.Err = cmd.Run()
	res.Duration = time.Since(now)

	if cmd.ProcessState != nil {
		res.Exit = cmd.ProcessState.ExitCode()
	}

	res.Stderr = stderr.Bytes()
	res.Stderr = bytes.TrimSpace(res.Stderr)

	res.Stdout = stdout.Bytes()
	res.Stdout = bytes.TrimSpace(res.Stdout)

	if res.Err == nil {
		return res, nil
	}

	err := RunError{
		Err:    res.Err,
		Args:   res.Args,
		Dir:    res.Dir,
		Env:    res.Env,
		Exit:   res.Exit,
		Output: append(res.Stdout, res.Stderr...),
	}

	return res, err
}
