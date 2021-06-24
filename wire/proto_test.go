package wire_test

import (
	"testing"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/stretchr/testify/assert"
)

func TestNewProtoDecoder(t *testing.T) {
	protoFile := "foo.proto"
	includes := []string{"./protos"}

	decoder := wire.NewProtoDecoder(protoFile, includes)
	assert.NotNil(t, decoder)
}
