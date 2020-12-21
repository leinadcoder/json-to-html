package jsontohtml

import (
	"encoding/json"
)

func toStruct(contents string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte (contents), &data)

	return data, err
}

func toMap(elements interface{}) map[string]string {
	m := make(map[string]string)

	switch elements.(type) {
	case map[string]interface{}:
        for k, v := range elements.(map[string]interface{}) {
            m[k] = v.(string)
        }
	}

	return m
}
