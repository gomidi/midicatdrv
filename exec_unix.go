// +build !windows

package midicatdrv

import (
	"os/exec"
	"syscall"
)

func execCommand(c string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", "exec "+c)
}

func midiCatCmd(args string) *exec.Cmd {
	cmd := execCommand("midicat " + args)
	// important! prevents that signals such as interrupt send to the main program gets passed
	// to midicat (which would not allow us to shutdown properly, e.g. stop hanging notes)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
