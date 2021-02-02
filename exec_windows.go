// +build windows

package midicatdrv

import (
	"os/exec"
)

func execCommand(c string) *exec.Cmd {
	//return exec.Command("powershell.exe", "/Command",  `$Process = [Diagnostics.Process]::Start("` + c + `") ; echo $Process.Id `)
	//return exec.Command("powershell.exe", "/Command", `$Process = [Diagnostics.Process]::Start("fluidsynth.exe", "-i -q -n $_file") ; echo $Process.Id `)
	return exec.Command("cmd.exe", "/C", c)
}

func midiCatCmd(args string) *exec.Cmd {
	return execCommand("midicat.exe " + args)
}
