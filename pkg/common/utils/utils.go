package utils

import (
	"encoding/json"
	"io"
	"time"
)

// ParseBody reads and parses the JSON data from the request body into the given interface.
func ParseBody(r io.Reader, x interface{}) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, x); err != nil {
		return err
	}

	return nil
}

func ConvertToUnix(timeValue time.Time) int64 {
	return timeValue.UTC().Unix()
}
