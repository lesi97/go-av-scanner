package utils

import (
	"encoding/json"
	"fmt"
)

func PrintPrettyJSON(value interface{}) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
