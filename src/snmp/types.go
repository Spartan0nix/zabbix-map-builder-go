package snmp

type Router_info struct {
	Snmp_index       string
	Local_hostname   string
	Local_interface  string
	Remote_hostname  string
	Remote_interface string
}

type Router struct {
	Local_hostname   string
	Local_interface  string
	Remote_hostname  string
	Remote_interface string
}
