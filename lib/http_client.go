package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HttpClient struct {
	Client http.Client
}

func NewHttpClient() HttpClient {
	return HttpClient{}
}

type IHttpClient interface {
	Get(url string, header map[string]string, requestConfig map[string]interface{}) (map[string]interface{}, error)
	Post(url string, header map[string]string, payload map[string]interface{}, requestConfig map[string]interface{}) (map[string]interface{}, error)
}

func (h *HttpClient) Get(url string, header map[string]string, requestConfig map[string]interface{}) (map[string]interface{}, error) {
	result, err := do(http.MethodGet, url, header, nil, requestConfig)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	return result, nil
}

func do(method string, url string, header map[string]string, payload map[string]interface{}, requestConfig map[string]interface{}) (map[string]interface{}, error) {
	var err error

	// Create an HTTP client
	client := &http.Client{}

	var jsonPayload []byte

	if payload != nil {
		jsonPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers (if required)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// Perform the HTTP request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Process the response
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", response.StatusCode)
		err = errors.New(response.Status)
		return nil, err
	}

	var jsonValue map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&jsonValue)

	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}
