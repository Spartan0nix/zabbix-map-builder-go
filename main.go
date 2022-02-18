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

func init_api() (api.Api, string) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	zabbix_url := os.Getenv("ZABBIX_URL")
	map_name := os.Getenv("ZABBIX_MAP_NAME")

	api := api.Api{
		Url:   zabbix_url,
		Token: "",
	}
	api.Token = api.Auth(os.Getenv("ZABBIX_USER"), os.Getenv("ZABBIX_USER_PASSWORD"))

	return api, map_name
}

func main() {
	FIXTURES := make([]Fixture, 0)
	FIXTURES = append(FIXTURES, Fixture{
		local_hostname:   "routeur-1",
		local_interface:  "FastEthernet1/0",
		remote_hostname:  "routeur-2",
		remote_interface: "FastEthernet2/0",
	})

	var Api, MAP_NAME = init_api()

	zabbix_map := Api.Get_map_by_name(MAP_NAME)
	var map_id string

	if len(zabbix_map.Result) == 0 {
		fmt.Println("Map not existing. Creating the map...")
		map_id = Api.Create_map(MAP_NAME)
	} else {
		map_id = zabbix_map.Result[0].Sysmapid
	}

	current_map := Api.Get_map_by_id(map_id).Result[0]
	current_map.Selements = make([]api.Map_selement, 0)
	current_map.Links = make([]api.Map_link, 0)

	Api.Update_map(current_map)

	for _, router := range FIXTURES {
		fmt.Println(router)

		local_hostid := Api.Get_host_id(router.local_hostname)
		if len(local_hostid.Result) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.local_hostname)
			break
		}
		remote_hostid := Api.Get_host_id(router.remote_hostname)
		if len(remote_hostid.Result) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.remote_hostname)
			break
		}

		local_router_hostid := local_hostid.Result[0].Hostid
		remote_router_hostid := remote_hostid.Result[0].Hostid

		fmt.Println(local_router_hostid)
		fmt.Println(remote_router_hostid)
	}

}
