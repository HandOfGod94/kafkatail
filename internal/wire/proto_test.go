package wire_test

import (
	"testing"

	"github.com/handofgod94/kafkatail/internal/wire"
	"github.com/handofgod94/kafkatail/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func TestNewProtoDecoder(t *testing.T) {
	protoFile := "foo.proto"
	includes := []string{"./protos"}

	decoder := wire.NewProtoDecoder(protoFile, "", includes)
	assert.NotNil(t, decoder)
}

func TestProtoDecoder(t *testing.T) {
	type fields struct {
		protoFile   string
		includes    []string
		messageType string
	}

	testCases := []struct {
		desc     string
		fields   fields
		protoMsg proto.Message
		wantErr  bool
	}{
		{
			desc:   "with valid messageType, proto file and includes path",
			fields: fields{"starwars.proto", []string{"../testdata"}, "Human"},
			protoMsg: &testdata.Human{
				HomePlanet: "earth",
				Mass:       32.0,
				Height: &testdata.Human_Height{
					Unit:  testdata.LengthUnit_METER,
					Value: 172.0,
				},
			},
			wantErr: false,
		},
		{
			desc:     "when proto file is not present",
			fields:   fields{"foo.proto", []string{"../testdata"}, "Droid"},
			protoMsg: &testdata.Human{HomePlanet: "earth", Mass: 32},
			wantErr:  true,
		},
		{
			desc:     "when message type is not present in proto path",
			fields:   fields{"starwars.proto", []string{"../testdata"}, "Foo"},
			protoMsg: &testdata.Human{Mass: 32, HomePlanet: "earth"},
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			bytes, _ := proto.Marshal(tc.protoMsg)
			p := wire.NewProtoDecoder(tc.fields.protoFile, tc.fields.messageType, tc.fields.includes)

			gotMsg, err := p.Decode(bytes)

			if tc.wantErr {
				var decodeErr *wire.DecodError
				assert.Error(t, err)
				assert.ErrorAs(t, err, &decodeErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, prototext.Format(tc.protoMsg), gotMsg)
			}
		})
	}
}
