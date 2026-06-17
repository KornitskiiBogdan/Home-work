package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("testdata/env")
	require.NoError(t, err)

	require.Equal(t, EnvValue{Value: `"hello"`}, env["HELLO"])
	require.Equal(t, EnvValue{Value: "bar"}, env["BAR"])
	require.Equal(t, EnvValue{Value: "   foo\nwith new line"}, env["FOO"])
	require.True(t, env["UNSET"].NeedRemove)
	require.True(t, env["EMPTY"].NeedRemove)
}
