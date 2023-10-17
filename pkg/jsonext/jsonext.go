package jsonext

import "encoding/json"

func Raw(data any) []byte {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	return b
}

func String(data any) string {
	return string(Raw(data))
}
