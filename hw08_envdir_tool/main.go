package main

import (
	"os"
)

func main() {
	args := os.Args
	lEnvs, err := ReadDir(args[1])
	if err != nil {
		os.Exit(1)
	}
	envs := MergeOSEnv(os.Environ(), lEnvs)
	exitCode := RunCmd(args[2:], envs)
	if exitCode > 0 {
		os.Exit(exitCode)
	}
}
