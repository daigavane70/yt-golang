package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"sprint/go/pkg/common/logger"
	"sprint/go/pkg/constants"
	"time"
)

// ParseBody reads and parses the JSON data from the request body into the given interface.
func ParseBody(r io.Reader, x interface{}) error {
	// Read the request body
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the provided interface
	if err := json.Unmarshal(body, x); err != nil {
		return err
	}

	return nil
}

// sendJSONResponse marshals data into JSON format and sends it as HTTP response with the specified status code.
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error marshaling JSON response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}

// ConvertToUnix converts a time.Time value to Unix timestamp (seconds since January 1, 1970 UTC).
func ConvertToUnix(timeValue time.Time) int64 {
	return timeValue.UTC().Unix()
}

// EpochToUTC converts an epoch timestamp (seconds since January 1, 1970 UTC) to UTC time in RFC3339 format.
func EpochToRFC3339(epoch int64) string {
	utcTime := time.Unix(epoch, 0).UTC()
	return utcTime.Format(time.RFC3339)
}

// FormatToRFC3339 formats a time.Time value to RFC3339 format in UTC.
func FormatToRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func FormatAsReadableTime(epoch int64) string {
	if epoch == 0 {
		return time.Now().Format(constants.READABLE_DATA_TIME_FORMAT)
	}
	return time.Unix(epoch, 0).Format(constants.READABLE_DATA_TIME_FORMAT)
}
