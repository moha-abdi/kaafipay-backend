package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON is a custom type for handling JSONB in PostgreSQL
type JSON json.RawMessage

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
		return nil
	case string:
		*j = JSON([]byte(v))
		return nil
	}

	return errors.New("invalid scan source for JSON")
}
