package clam

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RunError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	re := RunError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("wrapped: %w", re)

	r.True(errors.Is(wrapped, io.EOF))

	var ex RunError
	r.True(errors.As(wrapped, &ex))

	r.Equal(io.EOF, errors.Unwrap(re))
}
