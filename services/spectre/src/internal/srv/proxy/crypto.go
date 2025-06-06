package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spectre/internal/srv/lib/response"
	"time"
)

const (
	PROTO = "http://"
)

type Request struct {
	response.ResponseWithContent
	From string `json:"from"`
}

// CryptoClient provides methods for interacting with an external encryption service.
type CryptoClient struct {
	EncryptEndpoint string // URL of the encryption endpoint
	DecryptEndpoint string // URL of the decryption endpoint
	ECDHEndpoint    string
	Client          *http.Client // HTTP client for requests
}

// NewCryptoClient creates a new CryptoClient instance.
// epoint is the encryption endpoint URL, dpoint is the decryption endpoint URL.
func NewCryptoClient(epoint, dpoint, ecdhPoint string) *CryptoClient {
	return &CryptoClient{
		EncryptEndpoint: epoint,
		DecryptEndpoint: dpoint,
		ECDHEndpoint:    ecdhPoint,
		Client:          &http.Client{Timeout: 5 * time.Second}, // TODO: move to config
	}
}

// Encrypt sends a slice of base64 strings to the external encryption service.
// Returns a response.Response struct with the result or an error.
func (c *CryptoClient) Encrypt(b64 []string, from string) (response.ResponseWithContent, error) {
	reqBody, err := json.Marshal(Request{
		ResponseWithContent: response.ResponseWithContent{
			Content: b64,
		},
		From: from,
	})
	if err != nil {
		return response.EmptyWithContent, err
	}

	resp, err := c.Client.Post(
		PROTO+c.EncryptEndpoint,
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return response.EmptyWithContent, err
	}

	if resp.StatusCode != http.StatusOK {
		return response.EmptyWithContent, errBadStatusCode(resp.StatusCode)
	}

	var res response.ResponseWithContent
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.EmptyWithContent, err
	}

	return res, nil
}

// Decrypt sends the request body to the external decryption service and decodes the response.
// Returns a response.Response struct with the result or an error if the request fails or the response is invalid.
func (c *CryptoClient) Decrypt(r *http.Request) (response.ResponseWithContent, error) {

	if err := augmentRequestBody(r); err != nil {
		return response.EmptyWithContent, err
	}

	resp, err := c.Client.Post(
		PROTO+c.DecryptEndpoint,
		"application/json",
		r.Body,
	)
	if err != nil {
		return response.EmptyWithContent, err
	}

	if resp.StatusCode != http.StatusOK {
		return response.EmptyWithContent, errBadStatusCode(resp.StatusCode)
	}

	var res response.ResponseWithContent
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.EmptyWithContent, err
	}

	return res, nil
}

func (c *CryptoClient) GetK(r *http.Request) (response.ECDHResponse, error) {
	data := map[string]string{
		"from": r.Host,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return response.EmptyEDCH, err
	}

	resp, err := c.Client.Post(
		PROTO+c.ECDHEndpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return response.EmptyEDCH, err
	}

	if resp.StatusCode != http.StatusOK {
		return response.EmptyEDCH, errBadStatusCode(resp.StatusCode)
	}

	var res response.ECDHResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.EmptyEDCH, err
	}

	return res, nil
}

func (c *CryptoClient) SetA(r *http.Request) error {
	type setAreq struct {
		Key string `json:"key"`
	}

	var reqData setAreq
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return err
	}

	data := map[string]string{
		"from": r.Host,
		"key":  reqData.Key,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(
		PROTO+c.ECDHEndpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errBadStatusCode(resp.StatusCode)
	}

	return nil
}

func augmentRequestBody(r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()

	var rawJsonString string
	if err := json.Unmarshal(bodyBytes, &rawJsonString); err != nil {
		return fmt.Errorf("не удалось распарсить как JSON-строку: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rawJsonString), &data); err != nil {
		return fmt.Errorf("не удалось распарсить вложенный JSON: %w", err)
	}

	data["from"] = r.Host

	newBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(bytes.NewReader(newBody))
	r.ContentLength = int64(len(newBody))

	fmt.Println("Updated JSON:")
	fmt.Println(string(newBody))

	return nil
}
