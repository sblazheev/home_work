package main

import (
	"bufio"
	"os"
	"path"
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
	dirItems, err := os.ReadDir(dir)
	envs := make(Environment, len(dirItems))
	if err != nil {
		return nil, err
	}
	for _, item := range dirItems {
		if !item.IsDir() {
			content, errRead := readFirstLine(path.Join(dir, item.Name()))
			if errRead != nil {
				return nil, errRead
			}
			env := EnvValue{
				Value:      strings.TrimRight(content, " \t"),
				NeedRemove: false,
			}
			if len(content) == 0 {
				env.NeedRemove = true
			}
			envName := strings.ReplaceAll(item.Name(), "=", "")
			envs[envName] = env
		}
	}
	return envs, nil
}

func readFirstLine(path string) (res string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buf := scanner.Text()
		buf = strings.ReplaceAll(buf, "\x00", "\n")
		return buf, nil
	}
	return
}
