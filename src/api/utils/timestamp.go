package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Timestamp type wrapping time.Time
type Timestamp time.Time

// Standard SQLite date-time format
const internalTimeFormat = "2006-01-02 15:04:05"

var externalFormats = []string{
	"2006-01-02 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 GMT",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"2006-01-02T15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	time.RFC3339,
}

func ParseTimeStr(timeStr string) (time.Time, error) {
	timeStr = strings.Replace(timeStr, "+00:00", "+0000", 1)
	for _, format := range externalFormats {
		parsedTime, err := time.Parse(format, timeStr)
		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid RSS time format: %s", timeStr)
}

// MarshalJSON controls how Timestamp is serialized to JSON
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(internalTimeFormat))
}

// UnmarshalJSON controls how Timestamp is deserialized from JSON
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Unmarshal JSON into a string first
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	// Parse the string into time.Time
	parsedTime, err := time.Parse(internalTimeFormat, str)
	if err != nil {
		return err
	}
	*t = Timestamp(parsedTime)
	return nil
}

// Scan converts database value into time.Time
func (t *Timestamp) Scan(value any) error {
	if value == nil {
		*t = Timestamp(time.Time{}) // Set to zero time
		return nil
	}
	switch v := value.(type) {
	case string:
		parsedTime, err := ParseTimeStr(v)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		*t = Timestamp(parsedTime)
	case []byte:
		parsedTime, err := ParseTimeStr(string(v))
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		*t = Timestamp(parsedTime)
	case time.Time:
		*t = Timestamp(v)
	default:
		return fmt.Errorf("unexpected type for timestamp: %T", value)
	}
	return nil
}

// Value converts Go time.Time into SQLite format
func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t).Format(internalTimeFormat), nil
}
