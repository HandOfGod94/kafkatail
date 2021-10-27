// +build integration

package kafkatail_test

import (
	"context"
	"testing"
	"time"

	. "github.com/handofgod94/kafkatail/kafkatest"
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
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 kafkatail-plaintext-test-topic",
			message: "hello world",
			want:    "hello world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMessage(t, context.Background(), []string{LocalBroker}, kafkaPlainTextTopic, nil, []byte(tc.message))
			got := cmd.GetOutput()
			assert.Contains(t, string(got), tc.want)
		})
	}
}
