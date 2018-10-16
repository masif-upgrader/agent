package main

import (
	"os/exec"
)

type lazyLogString struct {
	generator func() string
}

func (s lazyLogString) String() string {
	return s.generator()
}

func (s lazyLogString) MarshalText() (text []byte, err error) {
	return []byte(s.generator()), nil
}

func getExePath(exe string) (path string, err error) {
	path, errLP := exec.LookPath(exe)
	if errLP != nil {
		if errCmd, isErrCmd := errLP.(*exec.Error); isErrCmd {
			if errCmd.Err == exec.ErrNotFound {
				return "", nil
			}
		}

		return "", errLP
	}

	return
}
