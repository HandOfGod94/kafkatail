package wire

type PlaintextDecoder struct{}

func NewPlaintextDecoder() *PlaintextDecoder {
	return &PlaintextDecoder{}
}

func (ptd *PlaintextDecoder) Decode(raw []byte) (string, error) {
	return string(raw), nil
}
