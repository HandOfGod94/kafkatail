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
		desc    string
		cmd     string
		want    string
		msg     string
		wantErr bool
	}{
		{
			desc: "check version",
			cmd:  "kafkatail --version",
			want: "dev",
		},
		{
			desc:    "print error for missing required args",
			cmd:     "kafkatail",
			want:    `requires at least 1 arg(s)`,
			wantErr: true,
		},
		{
			desc:    "print error for invalid wire_fomrat",
			cmd:     "kafkatail --bootstrap_servers=1.1.1.1:9093 --wire_format=foo test",
			want:    "must be 'avro', 'plaintext', 'proto'",
			wantErr: true,
		},
		{
			desc: "prints messages for valid args",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with offset option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "Offset: 0",
		},
		{
			desc: "with partition option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --partition=0 kafkatail-test-base",
			msg:  "hello world",
			want: "Partition: 0",
		},
		{
			desc: "with `from_datetime` option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with shorthand flag",
			cmd:  "kafkatail -b=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with spaces instead of `=` in command",
			cmd:  "kafkatail -b localhost:9093 --offset=-2 --group_id foo_id kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
	}

	const topic = "kafkatail-test-base"
	createTopic(t, context.Background(), topic)
	defer deleteTopic(t, context.Background(), topic)
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			appName, args := appNameAndArgs(tc.cmd)
			cmd := exec.CommandContext(ctx, appName, args...)
			out, err := getOutput(cmd, tc.wantErr)
			if err != nil {
				t.Log("failed to create output pipe:", err)
				t.FailNow()
			}

			if err := cmd.Start(); err != nil {
				t.Logf("failed to start command: '%v'. error: %v", tc.cmd, err)
				t.FailNow()
			}

			sendMessage(t, context.Background(), localBroker, topic, []byte(tc.msg))
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
