package wire

import "fmt"

type avroDecoder struct {
	schemaFile string
}

func NewAvroDecoder(schemaFile string) *avroDecoder {
	return &avroDecoder{schemaFile}
}

func (ad *avroDecoder) Decode(raw []byte) (string, error) {
	return "", fmt.Errorf("NotImplementedYet. unsupported format avro decoder")
}
