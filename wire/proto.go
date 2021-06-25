package wire

import (
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type protoDecoder struct {
	includes  []string
	protoFile string
}

func (pd *protoDecoder) Decode(raw []byte, messageType string) (string, error) {
	parser := protoparse.Parser{
		ImportPaths:      pd.includes,
		InferImportPaths: true,
	}

	descriptors, err := parser.ParseFiles(pd.protoFile)
	if err != nil {
		return "", &DecodError{Format: "proto", Reason: err}
	}

	protoDesc, err := protodesc.FileOptions{AllowUnresolvable: true}.New(descriptors[0].AsFileDescriptorProto(), nil)
	if err != nil {
		return "", &DecodError{Format: "proto", Reason: err}
	}

	messageDesc := protoDesc.Messages().ByName(protoreflect.Name(messageType))
	decodeMsg := dynamicpb.NewMessage(messageDesc)

	err = proto.UnmarshalOptions{AllowPartial: true}.Unmarshal(raw, decodeMsg)
	if err != nil {
		return "", &DecodError{Format: "proto", Reason: err}
	}

	return prototext.Format(decodeMsg), nil
}

func NewProtoDecoder(protoFile string, includes []string) *protoDecoder {
	return &protoDecoder{includes, protoFile}
}
