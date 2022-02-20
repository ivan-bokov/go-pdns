package stacktrace

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func unwrap(err error, n int) error {
	for i := 0; i < n; i++ {
		err = errors.Unwrap(err)
	}
	return err
}

func TestStdError(t *testing.T) {
	t.Parallel()

	err := errors.New("test")
	err1 := fmt.Errorf("test1: %w", err)
	err2 := fmt.Errorf("test2: %w", err1)

	assert.Equal(t, "test", err.Error())
	assert.Equal(t, "test1: test", err1.Error())
	assert.Equal(t, "test2: test1: test", err2.Error())

	assert.Equal(t, err, unwrap(err1, 1))
	assert.Equal(t, nil, unwrap(err1, 2))

	assert.Equal(t, err1, unwrap(err2, 1))
	assert.Equal(t, err, unwrap(err2, 2))
	assert.Equal(t, nil, unwrap(err2, 3))

	assert.True(t, errors.Is(err2, err2))
	assert.True(t, errors.Is(err2, err1))
	assert.True(t, errors.Is(err2, err))
}

func TestError(t *testing.T) {
	t.Parallel()

	err := New("test")
	err1 := Newf("test1: %w", err)
	err2 := Newf("test2: %w", err1)
	wrapErr := Wrap(err)

	errText := "[stacktrace_test.go:43 stacktrace.TestError] test"
	assert.Equal(t, errText, err.Error())
	err1Text := "[stacktrace_test.go:44 stacktrace.TestError] test1: " + errText
	assert.Equal(t, err1Text, err1.Error())
	err2Text := "[stacktrace_test.go:45 stacktrace.TestError] test2: " + err1Text
	assert.Equal(t, err2Text, err2.Error())
	wrapErrText := "[stacktrace_test.go:46 stacktrace.TestError] " + errText
	assert.Equal(t, wrapErrText, wrapErr.Error())

	assert.NotEqual(t, err, unwrap(err1, 1))
	assert.Equal(t, err, unwrap(err1, 2))
	assert.NotEqual(t, nil, unwrap(err1, 3))
	assert.Equal(t, nil, unwrap(err1, 4))

	assert.NotEqual(t, err1, unwrap(err2, 1))
	assert.Equal(t, err1, unwrap(err2, 2))
	assert.NotEqual(t, err, unwrap(err2, 3))
	assert.Equal(t, err, unwrap(err2, 4))
	assert.NotEqual(t, nil, unwrap(err2, 5))
	assert.Equal(t, nil, unwrap(err2, 6))

	assert.NotEqual(t, err1, unwrap(err2, 1))
	assert.Equal(t, err1, unwrap(err2, 2))
	assert.NotEqual(t, err, unwrap(err2, 3))
	assert.Equal(t, err, unwrap(err2, 4))
	assert.NotEqual(t, nil, unwrap(err2, 5))
	assert.Equal(t, nil, unwrap(err2, 6))

	assert.True(t, errors.Is(err2, err2))
	assert.True(t, errors.Is(err2, err1))
	assert.True(t, errors.Is(err2, err))
}
