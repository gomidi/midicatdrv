// +build !windows

package midicatdrv

import "os/exec"

func execCommand(c string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", "exec "+c)
}

func midiCatCmd(args string) *exec.Cmd {
	return execCommand("midicat " + args)
}
