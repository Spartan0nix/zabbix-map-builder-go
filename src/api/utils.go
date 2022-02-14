package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a Api) post(body []byte) []byte {
	resp, err := http.Post(a.Url, "application/json-rpc", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	return data
}
