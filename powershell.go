package main

import (
	"bytes"
	"os/exec"
)

// PowerShell represents a PowerShell instance.
type PowerShell struct {
	powerShell string
}

// NewPowerShell creates and returns a new PowerShell instance. It uses the
// 'exec' package's LookPath function to find the path to the 'powershell.exe'
// executable on the system, and uses this path to create the PowerShell instance.
//
// Returns:
//     *PowerShell: A pointer to the newly created PowerShell instance.
//
// Usage:
//     ps := NewPowerShell()
func NewPowerShell() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

// execute runs a given PowerShell command and returns its standard output,
// standard error, and any error that occurred during execution. The function
// constructs the command using the PowerShell instance and the provided arguments.
//
// The '-NoProfile' and '-NonInteractive' flags are added to the command to
// prevent the loading of the PowerShell profile and to ensure the command runs
// without requiring interactive user input.
//
// Parameters:
//     args (string): Zero or more string arguments that represent the command to be run.
//
// Returns:
//     stdOut (string): The standard output from the executed command.
//     stdErr (string): The standard error from the executed command.
//     err (error): An error object indicating any errors that occurred during command execution.
//
// Usage:
//     stdOut, stdErr, err := ps.execute("Get-Process")
func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
