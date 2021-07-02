// +build integration

package kafkatail_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKafkatailBase(t *testing.T) {
	testCases := []struct {
		cmd     string
		want    string
		wantErr bool
	}{
		{"kafkatail", `accepts 1 arg(s)`, true},
		{"kafkatail --bootstrap_servers=1.1.1.1:9093 --wire_format=foo test", "must be 'avro', 'plaintext', 'proto'", true},
		{"kafkatail --version", "0.1.0", false},
		{"kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test-base", "hello world", false},
	}

	createTopic(t, context.Background(), "kafkatail-test-base")
	defer deleteTopic(t, context.Background(), "kafkatail-test-base")
	for _, tc := range testCases {
		t.Run(tc.cmd, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			appName, args := appNameAndArgs(tc.cmd)
			cmd := exec.CommandContext(ctx, appName, args...)
			out, err := getOutput(cmd, tc.wantErr)
			if err != nil {
				t.Log("failed to create output pipe:", err)
				t.FailNow()
			}

			err = cmd.Start()
			if err != nil {
				t.Logf("failed to start command: '%v'. error: %v", tc.cmd, err)
				t.FailNow()
			}
			sendMessage(t, context.Background(), []string{"localhost:9093"}, "kafkatail-test-base", []byte(tc.want))
			got, err := io.ReadAll(out)
			if err != nil {
				t.Log("failed to read output:", err)
				t.FailNow()
			}

			assert.Contains(t, string(got), tc.want)

			cmd.Wait()
		})
	}
}
