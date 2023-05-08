package api

import (
	"fmt"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/utils"
)

// initService is used to return a new ZabbixService after executing connectivity test.
func initService(url string) (*zabbixgosdk.ZabbixService, error) {
	client := zabbixgosdk.NewZabbixService()

	client.Auth.Client.Url = url
	client.Map.Client.Url = url
	client.Trigger.Client.Url = url

	if err := client.Auth.Client.CheckConnectivity(); err != nil {
		return nil, err
	}

	return client, nil
}

// authenticate is used to retrieve an Api token for the HostGroup service.
func authenticate(client *zabbixgosdk.ZabbixService, user string, password string) error {
	u := &zabbixgosdk.ApiUser{
		User: user,
		Pwd:  password,
	}

	res, err := client.Auth.GetCredentials(u.User, u.Pwd)
	if err != nil {
		return err
	}

	if len(res.Result) == 0 {
		return fmt.Errorf("no token were returned during the authentification phase")
	}

	var token string
	err = client.Auth.Client.ConvertResponse(*res, &token)
	if err != nil {
		return err
	}

	client.Map.Client.Token = token
	client.Trigger.Client.Token = token

	return nil
}

// InitApi is used to initialize the default Zabbix service to interact with the API.
// A connectivity test is also run during this step.
func InitApi(url string, user string, password string) (*zabbixgosdk.ZabbixService, error) {
	client, err := initService(url)
	if err != nil {
		return nil, err
	}

	err = authenticate(client, user, password)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Logout is used to release the API token retrieve during the intialization of the API client.
func Logout(client *zabbixgosdk.ZabbixService) error {
	err := client.Logout()
	if err != nil {
		return err
	}

	return nil
}

// GetHostsId is used to retrive the id of the given hosts.
// The hosts name must be set as key in the map.
// Map value for each key will be replace by the id of the host retrieve from the Zabbix server.
func GetHostsId(client *zabbixgosdk.ZabbixService, hosts map[string]string) (map[string]string, error) {
	hostsName := utils.GetMapKey(hosts)

	h, err := client.Host.Get(&zabbixgosdk.HostGetParameters{
		Output: []string{
			"hostid",
			"host",
		},
		Filter: map[string][]string{
			"host": hostsName,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, host := range h {
		if _, exist := hosts[host.Host]; !exist {
			return nil, fmt.Errorf("unknown key '%s'", host.Host)
		}

		hosts[host.Host] = host.HostId
	}

	return hosts, nil
}

// GetImagesId is used to retrive the id of the given images.
// The images name must be set as key in the map.
// Map value for each key will be replace by the id of the host retrieve from the Zabbix server.
func GetImagesId(client *zabbixgosdk.ZabbixService, images map[string]string) (map[string]string, error) {
	imagesName := utils.GetMapKey(images)

	i, err := client.Image.Get(&zabbixgosdk.ImageGetParameters{
		Output: []string{
			"imageid",
			"name",
		},
		Filter: map[string][]string{
			"name": imagesName,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, image := range i {
		if _, exist := images[image.Name]; !exist {
			return nil, fmt.Errorf("missing key for image '%s'", image.Name)
		}

		images[image.Name] = image.ImageId
	}

	return images, nil
}
