package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty command", func(t *testing.T) {
		require.Equal(t, 1, RunCmd([]string{}, Environment{}))
	})

	t.Run("success exit code", func(t *testing.T) {
		code := RunCmd([]string{"/bin/sh", "-c", "exit 0"}, Environment{})
		require.Equal(t, 0, code)
	})

	t.Run("child exit code", func(t *testing.T) {
		code := RunCmd([]string{"/bin/sh", "-c", "exit 17"}, Environment{})
		require.Equal(t, 17, code)
	})

	t.Run("sets env variable", func(t *testing.T) {
		code := RunCmd([]string{"/bin/sh", "-c", `test "$TEMP = "bar"`}, Environment{"TEMP": {Value: "bar"}})
		require.Equal(t, 0, code)
	})

	t.Run("removes env variable", func(t *testing.T) {
		t.Setenv("TEMP", "old")
		code := RunCmd(
			[]string{"/bin/sh", "-c", `test -z "TEMP"`},
			Environment{"TEMP": {NeedRemove: true}},
		)
		require.Equal(t, 0, code)
	})

}
