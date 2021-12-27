// +build integration

package kafkatail_test

import (
	"context"
	"testing"
	"time"

	. "github.com/handofgod94/kafkatail/kafkatest"
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
			cmd:  `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-proto-test-topic`,
			want: "height:",
		},
		{
			desc:    "prints decode error when it fails to decode",
			cmd:     `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=foo.proto --include_paths="../testdata" --message_type=Human kafkatail-proto-test-topic`,
			want:    "failed to decode message",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			msg, err := proto.Marshal(&starwarsHuman)
			if err != nil {
				t.Log("failed to marshal message:", err)
				t.FailNow()
			}

			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMessage(t, []string{LocalBroker}, kafkaProtoTopic, nil, msg)

			got := cmd.GetOutput()
			assert.Contains(t, got, tc.want)
		})
	}
}
