package _map

import (
	"fmt"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// hostPosition is used to keep track of the position of an host on the map
type hostPosition struct {
	spacer int
	mapX   int
	mapY   int
	x      int
	y      int
}

// hostParameters define the parameters required to add an host element to a map
type hostParameters struct {
	id       string
	name     string
	image    string
	position *hostPosition
}

// initPosition is used to initialize the default placement values based on the given map width and height
func initPosition(width int, height int, spacer int) (*hostPosition, error) {
	// Initialize the default position of the first host to the left upper corner
	return &hostPosition{
		spacer: spacer,
		mapX:   width,
		mapY:   height,
		x:      spacer,
		y:      spacer,
	}, nil
}

// updateHostPosition is used to update the position for the next host
func (p *hostPosition) updateHostPosition() {
	// Move the host to the right
	newX := p.x + p.spacer
	// If the new X position is less than the less of the map, set it
	if newX < p.mapX {
		p.x = newX
	} else {
		// Otherwise, return to the left of the map
		p.x = p.spacer

		// If the host was pushed back to the left, we need to go one row down (Y)
		newY := p.y + p.spacer
		// If Y is greater than 0, set the new value
		if newY < p.mapY {
			p.y = newY
		} else {
			// Otherwise, return to the first row
			p.y = p.spacer
		}
	}
}

// createHostElement is used to create a new MapElementHost.
// The given id is used to reference the host in the map links.
func createHostElement(id string, host string, image string, x string, y string) *zabbixgosdk.MapElement {
	return &zabbixgosdk.MapElement{
		Id: id,
		Elements: []zabbixgosdk.MapElementHost{
			{
				Id: host,
			},
		},
		ElementType: zabbixgosdk.MapHost,
		IconIdOff:   image,
		X:           x,
		Y:           y,
	}
}

// addHosts is used to add hosts (local and remote) for a given mapping if they do not already exist in the map.
func addHosts(zbxMap *zabbixgosdk.MapCreateParameters, params *hostParameters) *zabbixgosdk.MapCreateParameters {
	exist := elementExist(params.id, zbxMap.Elements)
	if !exist {
		element := createHostElement(params.id, params.name, params.image, fmt.Sprintf("%d", params.position.x), fmt.Sprintf("%d", params.position.y))
		zbxMap.Elements = append(zbxMap.Elements, element)

		// Update placement for the next host
		params.position.updateHostPosition()
	}

	return zbxMap
}
