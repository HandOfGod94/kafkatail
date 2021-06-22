package main_test

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

func TestRequiredArgs(t *testing.T) {
	testCases := []struct {
		cmd     string
		want    string
		wantErr bool
	}{
		{"kafkatail", `"bootstrap_servers, topic" not set`, true},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			testcli.Run(tc.cmd)

			assert.Equal(t, tc.wantErr, testcli.Failure())
			assert.Contains(t, testcli.Stderr(), tc.want)
		})
	}
}
