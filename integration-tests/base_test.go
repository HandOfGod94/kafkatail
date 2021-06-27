// +build integration

package kafkatail_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKafkatailBase(t *testing.T) {
	testCases := []struct {
		cmd     string
		want    string
		wantErr bool
	}{
		{"kafkatail", `"bootstrap_servers", "topic" not set`, true},
		{"kafkatail --bootstrap_servers=1.1.1.1:9093 --topic=test --wire_format=foo", "must be 'avro', 'plaintext', 'proto'", true},
		{"kafkatail version", "0.1.0", false},
	}

	for _, tc := range testCases {
		t.Run(tc.cmd, func(t *testing.T) {
			tokens := strings.Split(tc.cmd, " ")
			appName := tokens[0]
			args := tokens[1:]

			out, err := exec.Command(appName, args...).Output()

			if tc.wantErr {
				exitErr, ok := err.(*exec.ExitError)
				assert.True(t, ok)
				assert.Contains(t, string(exitErr.Stderr), tc.want)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, string(out), tc.want)
			}
		})
	}
}
