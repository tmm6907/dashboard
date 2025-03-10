package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type UUID []byte

func (u UUID) MarshalJSON() ([]byte, error) {
	id, err := uuid.FromBytes(u)
	if err != nil {
		return nil, err
	}
	return json.Marshal(id.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*u = id[:]
	return nil
}

func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UUID: %v", value)
	}
	*u = bytes
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return []byte(u), nil
}
