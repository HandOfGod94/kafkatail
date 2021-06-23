package consumer_test

import "testing"

func TestConsume(t *testing.T) {
	testCases := []struct {
		desc             string
		bootstrapServers []string
		topic            string
		wantErr          bool
	}{
		{
			"with invalid consumer config",
			[]string{"localhost:9093"},
			"test_topic",
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

		})
	}
}
