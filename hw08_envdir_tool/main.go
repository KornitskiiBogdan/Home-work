package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(RunCmd(cmd, env))
}
