package main_test

import (
	"strings"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

const appName = "kafkatail"

func TestRequiredArgs(t *testing.T) {
	testCases := []struct {
		cliArgs string
		want    string
		wantErr bool
	}{
		{"", `"bootstrap_servers, topic" not set`, true},
		{"--bootstrap_servers=1.1.1.1:9093 --topic=test", "Hello World", false},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			testcli.Run(appName, strings.Split(tc.cliArgs, " ")...)

			if tc.wantErr {
				assert.True(t, testcli.Failure())
				assert.Contains(t, testcli.Stderr(), tc.want)
			} else {
				assert.True(t, testcli.Success())
				assert.Contains(t, testcli.Stdout(), tc.want)
			}
		})
	}
}
