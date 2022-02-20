package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"zabbix.builder/main/src/api"
	"zabbix.builder/main/src/snmp"
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

func init_app() (api.Api, string, string) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	zabbix_url := os.Getenv("ZABBIX_URL")
	map_name := os.Getenv("ZABBIX_MAP_NAME")
	router_ip := os.Getenv("ROUTER_IP")

	if router_ip == "" {
		log.Fatal("!! Env variable 'ROUTER_IP' not set !!")
		fmt.Printf("Stopping the program...")
		os.Exit(1)
	}

	api := api.Api{
		Url:   zabbix_url,
		Token: "",
	}
	api.Token = api.Auth(os.Getenv("ZABBIX_USER"), os.Getenv("ZABBIX_USER_PASSWORD"))

	return api, map_name, router_ip
}

func get_router_connections(router string) []snmp.Router {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	SNMP_COMMUNITY := os.Getenv("SNMP_COMMUNITY")
	SNMP_PORT := os.Getenv("SNMP_PORT")
	snmp_port, _ := strconv.ParseUint(SNMP_PORT, 10, 16)
	routers := []snmp.Router{}

	Snmp := snmp.Snmp_init(router, uint16(snmp_port), SNMP_COMMUNITY)
	defer Snmp.Conn.Close()

	local_hostname := snmp.Snmp_get_local_hostname(Snmp)

	remote_ip_oid := "1.3.6.1.4.1.9.9.23.1.2.1.1.4"
	res, err := Snmp.BulkWalkAll(remote_ip_oid)
	if err != nil {
		log.Fatalf("Error while retrieving oid '%s' : .Reason : %v", remote_ip_oid, err)
	}

	indexes := snmp.Extract_index(res)
	for _, index := range indexes {
		router_info := snmp.Router{
			Local_hostname: local_hostname,
		}

		router_info.Local_interface = snmp.Snmp_get_local_interface(Snmp, index)
		router_info.Remote_hostname = snmp.Snmp_get_remote_hostname(Snmp, index)
		router_info.Remote_interface = snmp.Snmp_get_remote_interface(Snmp, index)

		routers = append(routers, router_info)
	}

	return routers
}

func build_map(routers []snmp.Router, Api api.Api, current_map api.Map) {

	for _, router := range routers {

		local_hostid := Api.Host_get_id(router.Local_hostname)
		if len(local_hostid) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.Local_hostname)
			break
		}
		remote_hostid := Api.Host_get_id(router.Remote_hostname)
		if len(remote_hostid) == 0 {
			fmt.Printf("Host : '%s' does not exist on the zabbix server.", router.Remote_hostname)
			break
		}

		local_router_hostid := local_hostid[0].Hostid
		remote_router_hostid := remote_hostid[0].Hostid

		local_router_exist := Api.Map_selement_exist(current_map, local_router_hostid)
		if local_router_exist == "" {
			local_selement := Api.Map_build_selement(router.Local_hostname, local_router_hostid)
			current_map = Api.Map_add_selement(current_map, local_selement)
		} else {
			local_router_hostid = local_router_exist
		}

		remote_router_exist := Api.Map_selement_exist(current_map, remote_router_hostid)
		if remote_router_exist == "" {
			remote_selement := Api.Map_build_selement(router.Remote_hostname, remote_router_hostid)
			current_map = Api.Map_add_selement(current_map, remote_selement)
		} else {
			remote_router_hostid = remote_router_exist
		}

		local_trigger := Api.Trigger_get_id(local_router_hostid, router.Local_interface)
		remote_trigger := Api.Trigger_get_id(remote_router_hostid, router.Remote_interface)

		if !Api.Map_link_exist(current_map, router.Local_interface, router.Remote_interface, local_router_hostid, remote_router_hostid) {
			link := Api.Map_build_link(local_router_hostid,
				remote_router_hostid,
				router.Local_interface,
				router.Remote_interface,
				local_trigger[0].Triggerid,
				remote_trigger[0].Triggerid)

			current_map = Api.Map_add_link(current_map, link)
		}

		Api.Map_update(current_map)
	}
}

func main() {
	var Api, MAP_NAME, router_ip = init_app()

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
	// 	snmp function
	// ------------------------------------------------------------
	routeurs := get_router_connections(router_ip)

	// ------------------------------------------------------------
	// 	Build_map function
	// ------------------------------------------------------------
	build_map(routeurs, Api, current_map)

}
