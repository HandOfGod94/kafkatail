package wire

import (
	"fmt"

	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type ProtoDecoder struct {
	includes    []string
	protoFile   string
	messageType string
}

func NewProtoDecoder(protoFile, messageType string, includes []string) *ProtoDecoder {
	return &ProtoDecoder{includes, protoFile, messageType}
}

func (pd *ProtoDecoder) Decode(raw []byte) (string, error) {
	parser := protoparse.Parser{
		ImportPaths:      pd.includes,
		InferImportPaths: true,
	}

	descriptors, err := parser.ParseFiles(pd.protoFile)
	if err != nil {
		return "", ProtoDecodeError(err, "proto parser failed")
	}

	protoDesc, err := protodesc.FileOptions{AllowUnresolvable: true}.New(descriptors[0].AsFileDescriptorProto(), nil)
	if err != nil {
		return "", ProtoDecodeError(err, "failed to genereate proto descriptors")
	}

	messageDesc := protoDesc.Messages().ByName(protoreflect.Name(pd.messageType))
	if messageDesc == nil {
		return "", ProtoDecodeError(nil, "message type not found")
	}

	decodedMsg := dynamicpb.NewMessage(messageDesc)
	err = proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: false}.Unmarshal(raw, decodedMsg)
	if err != nil {
		return "", ProtoDecodeError(err, fmt.Sprintf("failed to unmarshal to type %v", pd.messageType))
	}

	return prototext.Format(decodedMsg), nil
}
