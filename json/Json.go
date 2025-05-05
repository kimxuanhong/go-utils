package json

import (
	"encoding/json"
	"fmt"
)

func ToJson[T any](obj *T) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error converting struct to JSON:", err)
		return "", err
	}
	return string(data), nil
}

func ToStruct[T any](data string, obj *T) error {
	err := json.Unmarshal([]byte(data), obj)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		return err
	}
	return nil
}
