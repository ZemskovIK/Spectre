package lib

import (
	"encoding/base64"
	"encoding/json"
)

// ToBase64 marshals any Go value to JSON and encodes the result as a base64 string.
// Returns the base64-encoded string or an error if marshalling fails.
func ToBase64(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)

	return b64, nil
}

// ToBase64Slice encodes a slice of any type to a slice of base64 strings.
// Each element is marshaled to JSON and then encoded as base64.
// Returns a slice of base64-encoded strings or an error if any element fails to encode.
func ToBase64Slice[T any](data []T) ([]string, error) {
	res := make([]string, 0, len(data))
	for _, v := range data {
		b64, err := ToBase64(v)
		if err != nil {
			return nil, err
		}
		res = append(res, b64)
	}
	return res, nil
}
