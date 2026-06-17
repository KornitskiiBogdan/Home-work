package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			continue
		}

		envValue, err := readFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		env[name] = envValue
	}

	return env, nil
}

func readFile(filename string) (EnvValue, error) {
	file, err := os.Open(filename)
	if err != nil {
		return EnvValue{}, err
	}

	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return EnvValue{}, err
	}

	bytesLine := bytes.ReplaceAll([]byte(line), []byte{0x00}, []byte{'\n'})
	line = strings.TrimRight(string(bytesLine), " \t")

	if line == "" {
		return EnvValue{NeedRemove: true}, nil
	}
	return EnvValue{Value: line}, nil
}
