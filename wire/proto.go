package wire

type protoDecoder struct {
	includes  []string
	protoFile string
}

func NewProtoDecoder(protoFile string, includes []string) *protoDecoder {
	return &protoDecoder{includes, protoFile}
}
