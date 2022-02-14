package api

import (
	"encoding/json"
	"fmt"
)

type auth_response struct {
	JSON_RPC string      `json:"jsonrpc"`
	Result   string      `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
	Id       int         `json:"id"`
}

func (a Api) Auth(user string, password string) string {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "user.login",
		Params: map[string]interface{}{
			"user":     user,
			"password": password,
		},
		Id: 1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	resp := a.post(body)
	data := auth_response{}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	return data.Result
}
