package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Try to unmarshal a request body into a value.
// Reading a request body empties it.
func TryUnmarshal(req *http.Request, out any) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}
