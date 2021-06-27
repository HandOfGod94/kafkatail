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
		desc string
		cmd  string
		want string
	}{
		{
			desc: "when topic doesn't exist",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --topic=foobar",
			want: "",
		},
		{
			desc: "when messages are present in topic",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --topic=kafkatail-test",
			want: "hello world",
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
			if err != nil {
				t.Log("failed to create stdout pipe:", err)
				t.FailNow()
			}

			err = cmd.Start()
			if err != nil {
				t.Logf("failed to start command: '%v'. error: %v", tc.cmd, err)
				t.FailNow()
			}
			sendMessage(t, context.Background(), []string{"localhost:9093"}, "kafkatail-test", tc.want)
			got, err := io.ReadAll(stdout)
			if err != nil {
				t.Log("failed to read stdout:", err)
				t.FailNow()
			}

			assert.Contains(t, string(got), tc.want)

			cmd.Wait()
		})
	}
}
