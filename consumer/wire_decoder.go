package consumer

type WireDecoder interface {
	Decode([]byte) (string, error)
}
