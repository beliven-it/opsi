package helpers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

func Request(method string, endpoint string, body []byte, queryMap map[string]string, headers map[string]string) ([]byte, error) {
	client := http.Client{}

	// Add body
	var bodyAsReader io.Reader
	if body != nil {
		bodyAsReader = bytes.NewReader(body)
	}

	// Create request
	request, err := http.NewRequest(method, endpoint, bodyAsReader)
	if err != nil {
		return nil, err
	}

	// Add query params
	if queryMap != nil {
		query := request.URL.Query()
		for key, value := range queryMap {
			query.Add(key, value)
		}
		request.URL.RawQuery = query.Encode()
	}

	// Add additional headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// Execute
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	messageAsBytes, _ := io.ReadAll(response.Body)
	message := string(messageAsBytes)

	defer response.Body.Close()
	if response.StatusCode >= http.StatusOK && response.StatusCode <= http.StatusIMUsed {
		return messageAsBytes, nil
	}

	return nil, errors.New(message)
}
