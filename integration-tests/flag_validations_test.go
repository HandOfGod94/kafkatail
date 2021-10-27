// +build integration

package kafkatail_test

import (
	"context"
	"testing"
	"time"

	. "github.com/handofgod94/kafkatail/kafkatest"
	"github.com/stretchr/testify/assert"
)

func TestFlagValidationForPlainText(t *testing.T) {
	testCases := []struct {
		desc string
		cmd  string
		want string
	}{
		{
			desc: "bootstrap_servers is required arg",
			cmd:  "kafkatail kafkatail-test",
			want: `required flag(s) "bootstrap_servers" not set`,
		},
		{
			desc: "bootstrap_servers is required arg with plaintext flag",
			cmd:  "kafkatail --wire_format=plaintext kafkatail-test",
			want: `required flag(s) "bootstrap_servers" not set`,
		},
		{
			desc: "invalid format for from_datetime",
			cmd:  "kafkatail --wire_format=plaintext --from_datetime=foobar --bootstrap_servers=localhost:9093 kafkatail-test",
			want: `invalid datetime provided`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: true}
			cmd.Execute(ctx)
			got := cmd.GetOutput()

			assert.Contains(t, got, tc.want)
		})
	}
}

func TestFlagValidationForProtoFormat(t *testing.T) {
	testCases := []struct {
		desc string
		cmd  string
		want string
	}{
		{
			desc: "bootstrap_servers is required arg with plaintext flag",
			cmd:  "kafkatail --wire_format=plaintext kafkatail-test",
			want: `required flag(s) "bootstrap_servers" not set`,
		},
		{
			desc: "proto_file is required for proto wire_format",
			cmd:  "kafkatail --wire_format=proto --bootstrap_servers=localhost:9093 kafkatail-test",
			want: `required flag(s) "include_paths", "message_type", "proto_file" not set`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: true}
			cmd.Execute(ctx)
			got := cmd.GetOutput()

			assert.Contains(t, got, tc.want)
		})
	}
}
