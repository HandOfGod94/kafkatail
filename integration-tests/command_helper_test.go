// +build integration

package kafkatail_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
)

type command struct {
	t       *testing.T
	outer   io.ReadCloser
	kmd     *exec.Cmd
	Cmd     string
	WantErr bool
}

func (c *command) execute(ctx context.Context) {
	appName, args := appNameAndArgs(c.Cmd)
	c.kmd = exec.CommandContext(ctx, appName, args...)
	out, err := getOutput(c.kmd, c.WantErr)
	if err != nil {
		c.t.Log("failed to create output pipe:", err)
		c.t.FailNow()
	}

	if err := c.kmd.Start(); err != nil {
		c.t.Logf("failed to start command: '%v'. error: %v", c.Cmd, err)
		c.t.FailNow()
	}

	c.outer = out
}

func (c *command) getOutput() string {
	got, err := io.ReadAll(c.outer)
	if err != nil {
		c.t.Log("failed to read output:", err)
		c.t.FailNow()
	}
	c.kmd.Wait()

	return string(got)
}
