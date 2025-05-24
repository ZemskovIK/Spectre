package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"spectre/internal/srv/response"
	"time"
)

type CryptoClient struct {
	EncryptEndpoint string
	DecryptEndpoint string
	Client          *http.Client
}

func NewCryptoClient(epoint, dpoint string) *CryptoClient {
	return &CryptoClient{
		EncryptEndpoint: epoint,
		DecryptEndpoint: dpoint,
		Client:          &http.Client{Timeout: 5 * time.Second}, // ! TODO : cfg
	}
}

func (c *CryptoClient) Encrypt(b64 []string) (response.Response, error) {
	reqBody, err := json.Marshal(map[string][]string{
		"content": b64,
	})
	if err != nil {
		return response.Response{}, err
	}

	resp, err := c.Client.Post(
		c.EncryptEndpoint,
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return response.Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return response.Response{}, errBadStatusCode(resp.StatusCode)
	}

	var res response.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Response{}, err
	}

	return res, nil
}
