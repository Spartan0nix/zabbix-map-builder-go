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
		remote_interface: "FastEthernet1/0",
	})
	FIXTURES = append(FIXTURES, Fixture{
		local_hostname:   "routeur-2",
		local_interface:  "FastEthernet2/0",
		remote_hostname:  "routeur-3",
		remote_interface: "FastEthernet2/0",
	})
	FIXTURES = append(FIXTURES, Fixture{
		local_hostname:   "routeur-2",
		local_interface:  "FastEthernet2/0",
		remote_hostname:  "routeur-3",
		remote_interface: "FastEthernet2/0",
	})

	var Api, MAP_NAME = init_api()

	zabbix_map := Api.Map_get_by_name(MAP_NAME)
	var map_id string

	if len(zabbix_map) == 0 {
		fmt.Println("Map not existing. Creating the map...")
		map_id = Api.Map_create(MAP_NAME)
	} else {
		map_id = zabbix_map[0].Sysmapid
	}

	current_map := Api.Map_get_by_id(map_id)[0]
	current_map.Selements = make([]api.Map_selement, 0)
	current_map.Links = make([]api.Map_link, 0)

	Api.Map_update(current_map)

	// ------------------------------------------------------------
	// 	Build_map function
	// ------------------------------------------------------------
	for _, router := range FIXTURES {
		// fmt.Println(router)

		local_hostid := Api.Host_get_id(router.local_hostname)
		if len(local_hostid) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.local_hostname)
			break
		}
		remote_hostid := Api.Host_get_id(router.remote_hostname)
		if len(remote_hostid) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.remote_hostname)
			break
		}

		local_router_hostid := local_hostid[0].Hostid
		remote_router_hostid := remote_hostid[0].Hostid

		local_router_exist := Api.Map_selement_exist(current_map, local_router_hostid)
		if local_router_exist == "" {
			local_selement := Api.Map_build_selement(router.local_hostname, local_router_hostid)
			current_map = Api.Map_add_selement(current_map, local_selement)
		} else {
			local_router_hostid = local_router_exist
		}

		remote_router_exist := Api.Map_selement_exist(current_map, remote_router_hostid)
		if remote_router_exist == "" {
			remote_selement := Api.Map_build_selement(router.remote_hostname, remote_router_hostid)
			current_map = Api.Map_add_selement(current_map, remote_selement)
		} else {
			remote_router_hostid = remote_router_exist
		}

		local_trigger := Api.Trigger_get_id(local_router_hostid, router.local_interface)
		remote_trigger := Api.Trigger_get_id(remote_router_hostid, router.remote_interface)

		if !Api.Map_link_exist(current_map, router.local_interface, router.remote_interface, local_router_hostid, remote_router_hostid) {
			link := Api.Map_build_link(local_router_hostid,
				remote_router_hostid,
				router.local_interface,
				router.remote_interface,
				local_trigger[0].Triggerid,
				remote_trigger[0].Triggerid)

			current_map = Api.Map_add_link(current_map, link)
		}

		Api.Map_update(current_map)
	}

}
