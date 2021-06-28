// +build integration

package kafkatail_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var starwarsHuman = testdata.Human{
	Mass:   32.0,
	Height: &testdata.Human_Height{Unit: testdata.LengthUnit_METER, Value: 1.0},
}

func TestKafkatalProto(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			desc: "prints decoded messages based on valid proto provided via args",
			cmd:  `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto`,
			want: "height:",
		},
		{
			desc:    "prints decode error when it fails to decode",
			cmd:     `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=foo.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto`,
			want:    "failed to decode message",
			wantErr: true,
		},
	}
	createTopic(t, context.Background(), "kafkatail-test-proto")
	defer deleteTopic(t, context.Background(), "kafkatail-test-proto")
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

			msg, err := proto.Marshal(&starwarsHuman)
			if err != nil {
				t.Log("failed to marshal message:", err)
				t.FailNow()
			}

			sendMessage(t, context.Background(), []string{"localhost:9093"}, "kafkatail-test-proto", msg)
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
