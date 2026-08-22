package utils

import (
	"encoding/json"
	"fmt"
)

func JsonPrettyPrint(json_string any) {
	var res, err = json.MarshalIndent(json_string, "", "    ")
	if err != nil {
		panic("can't format the json string")
	}
	fmt.Println(string(res))
}
