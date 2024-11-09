package fctools

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type MarshalOptions struct {
	Bytes2Hash     bool
	Timestamp2Date bool
}

func _marshal(msg proto.Message) ([]byte, error) {
	options := protojson.MarshalOptions{
		Indent:          "  ", // Pretty-print with indentation
		EmitUnpopulated: true, // Include fields with zero values
	}
	return options.Marshal(msg)
}

func replaceHashFields(data interface{}, opts MarshalOptions) {
	switch value := data.(type) {
	case map[string]interface{}:
		for k, v := range value {
			if opts.Bytes2Hash && (k == "hash" || k == "signature" || k == "signer") {
				bytes, err := base64.StdEncoding.DecodeString(v.(string))
				if err != nil {
					panic("Field value is not base64")
				}
				value[k] = "0x" + hex.EncodeToString(bytes)
			} else if opts.Timestamp2Date && (k == "timestamp") {
				value[k] = time.Unix(int64(v.(float64))+FARCASTER_EPOCH, 0)
			} else {
				replaceHashFields(v, opts)
			}
		}
	case []interface{}:
		for _, v := range value {
			replaceHashFields(v, opts)
		}
	}
}

func Marshal(msg proto.Message, opts MarshalOptions) ([]byte, error) {
	json_bytes, err := protojson.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal(json_bytes, &jsonData)
	if err != nil {
		return nil, err
	}

	replaceHashFields(jsonData, opts)

	updatedJsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}
	return updatedJsonBytes, nil
}
