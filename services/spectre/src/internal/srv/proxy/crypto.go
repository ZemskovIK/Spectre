package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"spectre/internal/srv/lib/response"
	"time"
)

// CryptoClient provides methods for interacting with an external encryption service.
type CryptoClient struct {
	EncryptEndpoint string       // URL of the encryption endpoint
	DecryptEndpoint string       // URL of the decryption endpoint
	Client          *http.Client // HTTP client for requests
}

// NewCryptoClient creates a new CryptoClient instance.
// epoint is the encryption endpoint URL, dpoint is the decryption endpoint URL.
func NewCryptoClient(epoint, dpoint string) *CryptoClient {
	return &CryptoClient{
		EncryptEndpoint: epoint,
		DecryptEndpoint: dpoint,
		Client:          &http.Client{Timeout: 5 * time.Second}, // TODO: move to config
	}
}

// Encrypt sends a slice of base64 strings to the external encryption service.
// Returns a response.Response struct with the result or an error.
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

// Decrypt sends the request body to the external decryption service and decodes the response.
// Returns a response.Response struct with the result or an error if the request fails or the response is invalid.
func (c *CryptoClient) Decrypt(r *http.Request) (response.Response, error) {
	resp, err := c.Client.Post(
		c.DecryptEndpoint,
		"application/json",
		r.Body,
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
