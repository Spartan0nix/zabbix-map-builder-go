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

type Map struct {
	Sysmapid string `json:"sysmapid"`
	Name     string `json:"name"`
	Height   string `json:"height"`
	Width    string `json:"width"`
}
