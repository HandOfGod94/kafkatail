package cmd

import (
	"fmt"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
)

func TestPartitionFlag(t *testing.T) {
	testCases := []struct {
		desc    string
		flagVal interface{}
		want    mo.Either[int, string]
		wantErr bool
	}{
		{
			desc:    "should set value with `all`",
			flagVal: "all",
			want:    mo.Right[int]("all"),
		},
		{
			desc:    "should set value with integer value",
			flagVal: 1,
			want:    mo.Left[int, string](1),
		},
		{
			desc:    "should return error for strings other than `all`",
			flagVal: "foobar",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fl := PartitionFlag{}
			err := fl.Set(fmt.Sprintf("%v", tc.flagVal))

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, fl.Get())
			}

		})
	}
}
