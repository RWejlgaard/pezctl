package util

import (
	"bytes"
	"os/exec"
)

func RunCommand(command string, shell string) (string, error) {
    if shell == "" {
        shell = "/bin/bash"
    }

    cmd := exec.Command(shell, "-c", command)
    var stdout bytes.Buffer
    cmd.Stdout = &stdout
    err := cmd.Run()

    return stdout.String(), err
}
