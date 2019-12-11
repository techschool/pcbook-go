package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// ProtobufToJSON converts protocol buffer message to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	return marshaler.MarshalToString(message)
}

// JSONToProtobufMessage converts JSON string to protocol buffer message
func JSONToProtobufMessage(data string, message proto.Message) error {
	return jsonpb.UnmarshalString(data, message)
}
