package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
)

func Unfold(data string) *utils.Success {
	var response = &utils.Success{}
	err := json.Unmarshal([]byte(data), response)
	if err != nil {
		panic(err)
	}

	return response
}

func UnfoldAccess(data string) *utils.AuthTokensObject {
	var res = &utils.AuthTokensObject{}
	err := json.Unmarshal([]byte(data), res)
	if err != nil {
		panic(err)
	}

	return res
}

func ErrorHandle() string {

	ErrorArray := utils.Success{Success: false, Message: "Access Denied", Data: ""}

	byteArray, err := json.Marshal(ErrorArray)
	if err != nil {
		fmt.Println(err)
	}

	return string(byteArray)
}
