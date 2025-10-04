package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	args := make([]string, 0, (len(cmd) - 1))
	envAr := make([]string, 0, len(env))
	for key, val := range env {
		envAr = append(envAr, key+"="+val.Value)
	}
	command := cmd[0]
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	execCmd := exec.Command(command, args...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = envAr
	err := execCmd.Run()
	if err != nil {
		returnCode = 1
	}
	return
}

func MergeOSEnv(envOS []string, localEnv Environment) Environment {
	envs := make(Environment, len(envOS)+len(localEnv))
	for _, envString := range envOS {
		env := strings.Split(envString, "=")
		envs[env[0]] = EnvValue{Value: env[1]}
	}
	for key, env := range localEnv {
		envs[key] = env
	}
	return envs
}
