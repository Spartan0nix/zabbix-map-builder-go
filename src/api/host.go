package api

import (
	"encoding/json"
	"fmt"
)

type host_response struct {
	JSON_RPC string      `json:"jsonrpc"`
	Result   []Host      `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
	Id       int         `json:"id"`
}

func handle_host_response(resp []byte) host_response {
	data := host_response{}
	// fmt.Println(string(resp))

	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data
}

func (a Api) Get_host_id(host string) host_response {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "host.get",
		Params: map[string]interface{}{
			"filter": map[string]interface{}{
				"host": []string{
					host,
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

	return handle_host_response(resp)
}
