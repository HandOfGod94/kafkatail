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

func TestKafkatailPlaintext(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		message string
		want    string
		wantErr bool
	}{
		{
			desc:    "when topic doesn't exist",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 foobar",
			message: "hello world",
			want:    "",
		},
		{
			desc:    "when messages are present in topic",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 kafkatail-test",
			message: "hello world",
			want:    "hello world",
		},
		{
			desc:    "with offset option",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test",
			message: "hello world",
			want:    "Offset: 0",
		},
		{
			desc:    "with partition option",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 --partition=0 kafkatail-test",
			message: "hello world",
			want:    "Partition: 0",
		},
		{
			desc:    "with `from_datetime` option",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test",
			message: "hello world",
			want:    "hello world",
		},
	}

	createTopic(t, context.Background(), "kafkatail-test")
	defer deleteTopic(t, context.Background(), "kafkatail-test")
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			appName, args := appNameAndArgs(tc.cmd)
			cmd := exec.CommandContext(ctx, appName, args...)
			stdout, err := cmd.StdoutPipe()
			stderr, err := cmd.StderrPipe()
			if err != nil {
				t.Log("failed to create stdout pipe:", err)
				t.FailNow()
			}

			err = cmd.Start()
			if err != nil {
				t.Logf("failed to start command: '%v'. error: %v", tc.cmd, err)
				t.FailNow()
			}
			sendMessage(t, context.Background(), []string{"localhost:9093"}, "kafkatail-test", []byte(tc.message))
			stream := streamToRead(tc.wantErr, stdout, stderr)
			got, err := io.ReadAll(stream)
			if err != nil {
				t.Log("failed to read stdout:", err)
				t.FailNow()
			}

			assert.Contains(t, string(got), tc.want)

			cmd.Wait()
		})
	}
}
