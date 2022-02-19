package api

import (
	"encoding/json"
	"fmt"
)

type trigger_response struct {
	JSON_RPC string      `json:"jsonrpc"`
	Result   []Trigger   `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
	Id       int         `json:"id"`
}

func trigger_handle_response(resp []byte) trigger_response {
	data := trigger_response{}
	// fmt.Println(string(resp))

	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data
}

func (a Api) Trigger_get_id(hostid string, host_interface string) []Trigger {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "trigger.get",
		Params: map[string]interface{}{
			"hostids": hostid,
			"filter": map[string]interface{}{
				"description": map[string]string{
					"name": "Interface " + host_interface + " is down",
				},
			},
		},
		Auth: a.Token,
		Id:   1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	resp := a.post(body)

	return trigger_handle_response(resp).Result
}
