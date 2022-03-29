package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

func GetGitConfig(key string) (string, error) {
	return ExecGit([]string{"config", "--get", "--null", key})
}

func ExecGit(gitArgs []string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", errors.New("wait status returned non-zero")
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000\n"), nil
}
