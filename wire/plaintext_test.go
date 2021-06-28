package wire_test

import (
	"testing"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/stretchr/testify/assert"
)

func TestNewPlaintextDecoder(t *testing.T) {
	decoder := wire.NewPlaintextDecoder()
	assert.NotNil(t, decoder)
}

func TestPlaintextDecoder(t *testing.T) {
	in := []byte("hello world")
	want := "hello world"

	decoder := wire.NewPlaintextDecoder()
	got, err := decoder.Decode(in)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
