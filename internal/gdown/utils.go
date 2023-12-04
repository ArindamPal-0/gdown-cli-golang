package gdown

import (
	"encoding/json"
	"fmt"
)

func Prettify(data interface{}) string {
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("❌ Could not convert data to json")
		return ""
	}
	return string(d)
}
