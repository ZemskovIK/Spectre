package lib

import (
	"encoding/base64"
	"encoding/json"
)

func ToBase64(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)

	return b64, nil
}

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
