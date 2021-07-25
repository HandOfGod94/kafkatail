// +build integration

package kafkatail_test

import (
	"io"
	"os/exec"
	"strings"
)

var localBroker = "localhost:9093"

func appNameAndArgs(cmd string) (appName string, args []string) {
	tokens := strings.Split(cmd, " ")
	appName = tokens[0]
	args = tokens[1:]
	return
}

func getOutput(cmd *exec.Cmd, wantErr bool) (io.ReadCloser, error) {
	if wantErr {
		return cmd.StderrPipe()
	} else {
		return cmd.StdoutPipe()
	}
}

func streamToRead(wantErr bool, stdout, stderr io.ReadCloser) io.ReadCloser {
	if wantErr {
		return stderr
	}
	return stdout
}
