package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]interface{}

func (data JSONB) Value() (driver.Value, error) {
	return json.Marshal(data)
}

func (data *JSONB) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var newData interface{}
	err := json.Unmarshal(source, &newData)
	if err != nil {
		return err
	}

	*data, ok = newData.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}
