package wire

type plaintextDecoder struct{}

func NewPlaintextDecoder() *plaintextDecoder {
	return &plaintextDecoder{}
}

func (ptd *plaintextDecoder) Decode(raw []byte) (string, error) {
	return string(raw), nil
}
