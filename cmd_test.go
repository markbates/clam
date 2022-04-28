package clam

import (
	"context"
	"os"
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

func Test_Cmd(t *testing.T) {
	t.Parallel()

	cmd := &Cmd{}
	res, err := cmd.Run(context.Background(), "echo", "hello")

	noErr(t, err)

	eq(t, 0, res.Exit)

	act := strings.TrimSpace(string(res.Stdout))
	exp := `hello`

	eq(t, exp, act)
	eq(t, "$ echo hello", res.CmdString())
}

func Test_Cmd_Unknown(t *testing.T) {
	t.Parallel()

	cmd := &Cmd{}
	_, err := cmd.Run(context.Background(), "ehho", "hello")

	if err == nil {
		t.Fatal("expected error")
	}

}

func Test_Cmd_Bad(t *testing.T) {
	t.Parallel()

	cmd := &Cmd{}
	res, err := cmd.Run(context.Background(), "go", "hello")

	if err == nil {
		t.Fatal("expected error")
	}

	eq(t, 2, res.Exit)

	act := strings.TrimSpace(string(res.Stderr))
	exp := "go hello: unknown command\nRun 'go help' for usage."

	eq(t, exp, act)
}

func Test_Cmd_Env(t *testing.T) {
	t.Parallel()

	exp := "FOO=bar"
	env := append(os.Environ(), exp)
	cmd := &Cmd{
		Env: env,
		Dir: "testdata/env",
	}

	args := []string{"go", "run", "main.go"}
	res, err := cmd.Run(context.Background(), args...)

	noErr(t, err)

	eq(t, 0, res.Exit)

	sliceEq(t, args, res.Args)
	sliceEq(t, env, res.Env)

	eq(t, cmd.Dir, res.Dir)

	act := strings.TrimSpace(string(res.Stdout))

	if !strings.Contains(act, exp) {
		t.Fatalf("expected %q, got %q", exp, act)
	}
}

func noErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func sliceEq[T constraints.Ordered](t testing.TB, exp []T, act []T) {
	t.Helper()

	if len(exp) != len(act) {
		t.Fatalf("expected %d, got %d", len(exp), len(act))
	}

	for i, e := range exp {
		if e != act[i] {
			t.Fatalf("expected %v, got %v", exp, act)
		}
	}
}

func eq[T constraints.Ordered](t testing.TB, exp T, act T) {
	t.Helper()

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}
}
