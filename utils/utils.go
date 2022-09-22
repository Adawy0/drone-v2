package utils

import (
	"encoding/json"
	"errors"
	"regexp"
)

// utils.TypeConverter
func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func ValidateMedicationName(name string) error {
	vaildName, _ := regexp.MatchString("^[a-zA-Z0-9_.-]*$", name)
	if !vaildName {
		return errors.New("invald format for medciation name, can be only combination of (letters, numbers, -, _)")
	}
	return nil
}
