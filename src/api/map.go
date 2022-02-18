package api

import (
	"encoding/json"
	"fmt"
)

type map_response struct {
	JSON_RPC string      `json:"jsonrpc"`
	Result   []Map       `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
	Id       int         `json:"id"`
}

type map_create_response struct {
	JSON_RPC string              `json:"jsonrpc"`
	Result   map[string][]string `json:"result,omitempty"`
	Error    interface{}         `json:"error,omitempty"`
	Id       int                 `json:"id"`
}

func handle_map_response(resp []byte) map_response {
	data := map_response{}
	// fmt.Println(string(resp))

	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data
}

func (a Api) extract_map_id(resp []byte) string {
	data := map_create_response{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data.Result["sysmapids"][0]
}

func (a Api) Get_map_by_name(name string) map_response {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "map.get",
		Params: map[string]interface{}{
			"search": map[string]interface{}{
				"name": name,
			},
		},
		Auth: a.Token,
		Id:   1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	response := a.post(body)

	return handle_map_response(response)
}

func (a Api) Get_map_by_id(id string) map_response {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "map.get",
		Params: map[string]string{
			"output":           "extend",
			"selectSelements":  "extend",
			"selectLinks":      "extend",
			"selectUsers":      "extend",
			"selectUserGroups": "extend",
			"selectShapes":     "extend",
			"selectLines":      "extend",
			"sysmapids":        id,
		},
		Auth: a.Token,
		Id:   1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	response := a.post(body)

	return handle_map_response(response)
}

func (a Api) Create_map(name string) string {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "map.create",
		Params: map[string]interface{}{
			"name":   name,
			"width":  600,
			"height": 600,
		},
		Auth: a.Token,
		Id:   1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	resp := a.post(body)

	return a.extract_map_id(resp)
}

func (a Api) Update_map(raw_map Map) string {
	payload := Payload{
		JSON_RPC: "2.0",
		Method:   "map.update",
		Params:   raw_map,
		Auth:     a.Token,
		Id:       1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	resp := a.post(body)

	return a.extract_map_id(resp)
}
