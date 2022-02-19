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

func map_handle_response(resp []byte) map_response {
	data := map_response{}
	// fmt.Println(string(resp))

	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data
}

func (a Api) map_extract_id(resp []byte) string {
	data := map_create_response{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return data.Result["sysmapids"][0]
}

func (a Api) Map_get_by_name(name string) []Map {
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

	return map_handle_response(response).Result
}

func (a Api) Map_get_by_id(id string) []Map {
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

	return map_handle_response(response).Result
}

func (a Api) Map_create(name string) string {
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

	return a.map_extract_id(resp)
}

func (a Api) Map_update(raw_map Map) string {
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

	return a.map_extract_id(resp)
}

func (a Api) Map_build_selement(name string, hostid string) Map_selement {
	element := Map_selement{}

	element.Selementid = hostid
	element.Elementtype = "0"
	element.Elements = append(element.Elements, map[string]string{
		"hostid": hostid,
	})
	element.Iconid_off = "156"
	element.Label = name
	element.Label_location = "-1"
	element.Inherited_label = name
	element.Label_type = "2"
	element.ElementName = name

	return element
}

func (a Api) Map_selement_exist(zabbix_map Map, selement_hostid string) string {
	for _, selement := range zabbix_map.Selements {
		if selement.Elements[0]["hostid"] == selement_hostid {
			return selement.Selementid
		}
	}
	return ""
}

func (a Api) Map_add_selement(zabbix_map Map, selement Map_selement) Map {
	zabbix_map.Selements = append(zabbix_map.Selements, selement)
	return zabbix_map
}

func (a Api) Map_build_link(local_hostid string,
	remote_hostid string,
	local_interface string,
	remote_interface string,
	local_trigger string,
	remote_trigger string,
) Map_link {
	link := Map_link{}

	link.Selementid1 = local_hostid
	link.Selementid2 = remote_hostid
	link.Label = local_interface + " - " + remote_interface
	link.Linktriggers = append(link.Linktriggers, map[string]string{
		"triggerid": local_trigger,
		"color":     "DD0000",
	})
	link.Linktriggers = append(link.Linktriggers, map[string]string{
		"triggerid": remote_trigger,
		"color":     "DD0000",
	})
	link.Color = "00CC00"

	return link
}

func (a Api) Map_link_exist(zabbix_map Map,
	local_interface string,
	remote_interface string,
	local_hostid string,
	remote_hostid string,
) bool {
	label := local_interface + " - " + remote_interface
	reverse_label := remote_interface + " - " + local_interface
	for _, link := range zabbix_map.Links {
		if link.Label == label {
			if link.Selementid1 == local_hostid && link.Selementid2 == remote_hostid {
				return true
			}
		}
		if link.Label == reverse_label {
			if link.Selementid1 == remote_hostid && link.Selementid2 == local_hostid {
				return true
			}
		}
	}
	return false
}

func (a Api) Map_add_link(zabbix_map Map, link Map_link) Map {
	zabbix_map.Links = append(zabbix_map.Links, link)
	return zabbix_map
}
