package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPResponseError extracts any error from a response, or returns nil if there is no error
func HTTPResponseError(response *http.Response) error {
	if response != nil && response.StatusCode >= 300 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("got status %s", response.Status)
		}

		var responseContent any

		err = json.Unmarshal(responseBody, &responseContent)
		if err != nil {
			responseContent = string(responseBody)
		}

		return fmt.Errorf("%s: %s", response.Status, responseContent)
	}

	return nil
}
