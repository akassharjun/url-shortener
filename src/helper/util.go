package helper

import "encoding/json"

func ConvertMapToStruct(data map[string]interface{}, v interface{}) error {
	// Convert the map to a JSON payload
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Unmarshal the JSON payload into the struct
	err = json.Unmarshal(jsonPayload, v)
	if err != nil {
		return err
	}

	return nil
}
