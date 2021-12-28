package kafkatest

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"
)

var LocalBroker = "localhost:9093"

type Command struct {
	T       *testing.T
	outer   io.ReadCloser
	kmd     *exec.Cmd
	Cmd     string
	WantErr bool
}

func (c *Command) Execute(ctx context.Context) {
	appName, args := AppNameAndArgs(c.Cmd)
	c.kmd = exec.CommandContext(ctx, appName, args...)
	out, err := GetOutput(c.kmd, c.WantErr)
	if err != nil {
		c.T.Log("failed to create output pipe:", err)
		c.T.FailNow()
	}

	if err := c.kmd.Start(); err != nil {
		c.T.Logf("failed to start command: '%v'. error: %v", c.Cmd, err)
		c.T.FailNow()
	}

	c.outer = out
}

func (c *Command) GetOutput() string {
	got, err := io.ReadAll(c.outer)
	if err != nil {
		c.T.Log("failed to read output:", err)
		c.T.FailNow()
	}
	c.kmd.Wait()

	return string(got)
}

func AppNameAndArgs(cmd string) (appName string, args []string) {
	tokens := strings.Split(cmd, " ")
	appName = tokens[0]
	args = tokens[1:]
	return
}

func GetOutput(cmd *exec.Cmd, wantErr bool) (io.ReadCloser, error) {
	if wantErr {
		return cmd.StderrPipe()
	} else {
		return cmd.StdoutPipe()
	}
}

func StreamToRead(wantErr bool, stdout, stderr io.ReadCloser) io.ReadCloser {
	if wantErr {
		return stderr
	}
	return stdout
}

func MinifyString(str string) string {
	miniStr := strings.TrimSpace(str)
	miniStr = strings.ReplaceAll(miniStr, " ", "")
	miniStr = strings.ReplaceAll(miniStr, "\t", "")
	miniStr = strings.ReplaceAll(miniStr, "\n", "")
	return miniStr
}
