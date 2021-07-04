package wire

type Decoder interface {
	Decode([]byte) (string, error)
}
