// Package util
/**
  @author: zqwh_yw1
  @date: 2021/12/10
  @comment:
**/
package util

import (
	"os/exec"
)

func ExecUnixCmd(cmd string) ([]byte, error) {
	return exec.Command("/bin/bash", "-c", cmd).Output()
}

func ExecWinCmd(cmdWithArgs ...string) ([]byte, error) {
	var fullCmd []string
	fullCmd = append(fullCmd, "/c")
	for _, arg := range cmdWithArgs {
		fullCmd = append(fullCmd, arg)
	}
	return exec.Command("cmd.exe", fullCmd...).Output()
}
