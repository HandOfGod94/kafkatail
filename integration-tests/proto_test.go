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

func TestKafkatalProto(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			desc:    "fails when user doesn't provide required proto args",
			cmd:     "kafkatail  --wire_format=proto kafkatail-test",
			want:    `"protoFile", "incldues" not set`,
			wantErr: true,
		},
		// {
		// 	desc: "prints decoded messages based on valid proto provided via args",
		// },
		// {
		// 	desc: "prints decode error when it fails to decode",
		// },
	}
	t.Skip()

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
