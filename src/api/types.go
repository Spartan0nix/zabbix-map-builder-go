package api

type Api struct {
	Url   string
	Token string
}

type Payload struct {
	JSON_RPC string      `json:"jsonrpc"`
	Method   string      `json:"method"`
	Params   interface{} `json:"params"`
	Auth     string      `json:"auth,omitempty"`
	Id       int         `json:"id"`
}

type Response struct {
	JSON_RPC string      `json:"jsonrpc"`
	Result   interface{} `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
	Id       int         `json:"id"`
}

type Map_selement struct {
	Selementid      string              `json:"selementid"`
	Elementtype     string              `json:"elementtype"`
	Elements        []map[string]string `json:"elements"`
	Iconid_off      string              `json:"iconid_off"`
	Label           string              `json:"label"`
	Label_location  string              `json:"label_location"`
	Inherited_label string              `json:"inherited_label"`
	Label_type      string              `json:"label_type"`
	ElementName     string              `json:"elementName"`
	Color           string              `json:"color"`
}

type Map_link struct {
	Link_id      string              `json:"linkid"`
	Sysmapid     string              `json:"sysmapid"`
	Selementid1  string              `json:"selementid1"`
	Selementid2  string              `json:"selementid2"`
	Color        string              `json:"color"`
	Label        string              `json:"label"`
	Linktriggers []map[string]string `json:"linktriggers"`
}

type Map struct {
	Sysmapid  string         `json:"sysmapid"`
	Name      string         `json:"name"`
	Height    string         `json:"height"`
	Width     string         `json:"width"`
	Selements []Map_selement `json:"selements"`
	Links     []Map_link     `json:"links"`
}

type Host struct {
	Hostid      string `json:"hostid"`
	Host        string `json:"host"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Trigger struct {
	Triggerid   string `json:"triggerid"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
