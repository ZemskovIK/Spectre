package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"spectre/internal/models"
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

func FetchLetterFromB64(l *models.Letter, b64 interface{}) error {
	sl, ok := b64.([]interface{})
	if !ok {
		return errors.New("b64 is not a slice of")
	}
	sb64, ok := sl[0].(string)
	if !ok {
		return errors.New("not a string in slice")
	}
	if sb64 == "" {
		return errors.New("b64 string is empty")
	}
	data, err := base64.StdEncoding.DecodeString(sb64)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("decoded base64 is empty")
	}
	if err := json.Unmarshal(data, l); err != nil {
		return err
	}
	return nil
}
