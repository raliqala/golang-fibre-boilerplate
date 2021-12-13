package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
)

func Unfold(data string) *utils.Success {
	response := &utils.Success{}
	err := json.Unmarshal([]byte(data), response)
	if err != nil {
		panic(err)
	}

	return response
}

func UnfoldAccess(data string) *utils.AuthTokensObject {
	res := &utils.AuthTokensObject{}
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

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")

	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}
