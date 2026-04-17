package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSON(r *http.Request, target interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, target)
}

func ReadJSONFromBytes(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
