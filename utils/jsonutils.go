package utils

import (
	"encoding/json"
	"go/types"
	"io/ioutil"
)

func UnmarshalFileJson(filename string, v interface{}) error {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return types.Error{}
	} else {
		if len(data) < 1 {
			return nil
		}

		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		return nil
	}
}

func MarshalFileJson(filename string, v interface{}) error {
	data, err := json.Marshal(v)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0755)

	if err != nil {
		return err
	}

	return nil
}
