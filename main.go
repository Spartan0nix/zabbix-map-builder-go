package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"zabbix.builder/main/src/api"
)

// -------------------------------------------------
// Static
// -------------------------------------------------
type Fixture struct {
	local_hostname   string
	local_interface  string
	remote_hostname  string
	remote_interface string
}

func main() {
	FIXURES := Fixture{
		local_hostname:   "routeur-1",
		local_interface:  "FastEthernet1/0",
		remote_hostname:  "routeur-2",
		remote_interface: "FastEthernet2/0",
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	ZABBIX_URL := os.Getenv("ZABBIX_URL")
	MAP_NAME := os.Getenv("ZABBIX_MAP_NAME")
	api := api.Api{
		Url:   ZABBIX_URL,
		Token: "",
	}
	api.Token = api.Auth(os.Getenv("ZABBIX_USER"), os.Getenv("ZABBIX_USER_PASSWORD"))

	fmt.Println(FIXURES)
	fmt.Println(ZABBIX_URL)

	zabbix_map := api.Get_by_name(MAP_NAME)
	var map_id string

	if len(zabbix_map.Result) == 0 {
		fmt.Println("Map not existing. Creating the map...")
		map_id = api.Create_map(MAP_NAME)
	} else {
		map_id = zabbix_map.Result[0].Sysmapid
	}

	fmt.Println(map_id)

}
