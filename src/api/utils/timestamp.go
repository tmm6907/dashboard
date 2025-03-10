package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp type wrapping time.Time
type Timestamp time.Time

// Standard SQLite date-time format
const timeFormat = "2006-01-02 15:04:05"

// MarshalJSON controls how Timestamp is serialized to JSON
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(timeFormat))
}

// UnmarshalJSON controls how Timestamp is deserialized from JSON
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Unmarshal JSON into a string first
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	// Parse the string into time.Time
	parsedTime, err := time.Parse(timeFormat, str)
	if err != nil {
		return err
	}
	*t = Timestamp(parsedTime)
	return nil
}

// Scan converts database value into time.Time
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*t = Timestamp(time.Time{}) // Set to zero time
		return nil
	}
	switch v := value.(type) {
	case string:
		parsedTime, err := time.Parse(timeFormat, v)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %v", err)
		}
		*t = Timestamp(parsedTime)
	case []byte:
		parsedTime, err := time.Parse(timeFormat, string(v))
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
	return time.Time(t).Format(timeFormat), nil
}
