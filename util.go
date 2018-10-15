package main

import (
	"fmt"
	"os/exec"
)

type lazyLogString struct {
	stringer fmt.Stringer
}

func (s lazyLogString) String() string {
	return s.stringer.String()
}

func (s lazyLogString) MarshalText() (text []byte, err error) {
	return []byte(s.stringer.String()), nil
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
