package _map

import (
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

const (
	width  = 500
	height = 500
	spacer = 100
)

func TestInitPosition(t *testing.T) {
	p, err := initPosition(width, height, spacer)
	if err != nil {
		t.Fatalf("error while executing initPosition function.\nReason : %v", err)
	}

	if p == nil {
		t.Fatalf("expected a hostParameters pointer to be returned. An nil value was returned")
	}

	if p.spacer != spacer {
		t.Fatalf("wrong value assigned to 'spacer' field.\nExpected : %d.\nReturned : %d", spacer, p.spacer)
	}

	if p.mapX != width {
		t.Fatalf("wrong value assigned to 'mapX' field.\nExpected : %d.\nReturned : %d", width, p.mapX)
	}

	if p.mapY != height {
		t.Fatalf("wrong value assigned to 'mapY' field.\nExpected : %d.\nReturned : %d", height, p.mapY)
	}

	if p.x != spacer {
		t.Fatalf("wrong value assigned to 'x' field.\nExpected : %d.\nReturned : %d", spacer, p.x)
	}

	if p.y != spacer {
		t.Fatalf("wrong value assigned to 'y' field.\nExpected : %d.\nReturned : %d", spacer, p.y)
	}
}

func BenchmarkInitPosition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initPosition(width, height, spacer)
	}
}

func TestUpdateHostPosition(t *testing.T) {
	// Update only the X position
	expectedX := spacer * 2
	expectedY := spacer

	p := hostPosition{
		spacer: spacer,
		mapX:   width,
		mapY:   height,
		x:      spacer,
		y:      spacer,
	}

	p.updateHostPosition()

	if p.x != expectedX {
		t.Fatalf("wrong 'x' value returned after position update.\nExpected : %d.\nReturned : %d", expectedX, p.x)
	}

	if p.y != expectedY {
		t.Fatalf("wrong 'y' value returned after position update.\nExpected : %d.\nReturned : %d", expectedY, p.y)
	}

	// Update the X and Y to return to the top left corner (x > width && y > height)
	expectedX = spacer
	expectedY = spacer

	p.x = width - 50
	p.y = height - 50

	p.updateHostPosition()

	if p.x != expectedX {
		t.Fatalf("wrong 'x' value returned after position update.\nExpected : %d.\nReturned : %d", expectedX, p.x)
	}

	if p.y != expectedY {
		t.Fatalf("wrong 'y' value returned after position update.\nExpected : %d.\nReturned : %d", expectedY, p.y)
	}
}

func BenchmarkUpdateHostPosition(b *testing.B) {
	p := hostPosition{
		spacer: spacer,
		mapX:   width,
		mapY:   height,
		x:      spacer,
		y:      spacer,
	}

	for i := 0; i < b.N; i++ {
		p.updateHostPosition()
	}
}

func TestCreateHostElement(t *testing.T) {
	element := createHostElement("2", "1", "11", "135", "135")
	elementHosts := element.Elements.([]zabbixgosdk.MapElementHost)

	if element.Id != "2" {
		t.Fatalf("wrong element id set.\nExpected : '2'\nReturned : %s", element.Id)
	}

	if len(elementHosts) != 1 {
		t.Fatalf("wrong number of host elements set.\nExpected : 1\nReturned : %d", len(elementHosts))
	}

	if elementHosts[0].Id != "1" {
		t.Fatalf("wrong id set for the host element.\nExpected : '1'\nReturned : %s", elementHosts[0].Id)
	}

	if element.ElementType != zabbixgosdk.MapHost {
		t.Fatalf("wrong type of elementType set.\nExpected : %s\nReturned : %v", zabbixgosdk.MapHost, element.ElementType)
	}

	if element.IconIdOff != "11" {
		t.Fatalf("wrong icon id set.\nExpected : 11\nReturned : %s", element.IconIdOff)
	}
}

func BenchmarkCreateHostElement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createHostElement("router-1", "router-1", "11", "135", "135")
	}
}

func TestAddHosts(t *testing.T) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap = addHosts(zbxMap, &hostParameters{
		id:    "2",
		name:  "1",
		image: "1",
		position: &hostPosition{
			spacer: 10,
			mapX:   100,
			mapY:   100,
			x:      10,
			y:      10,
		},
	})

	if len(zbxMap.Elements) != 1 {
		t.Fatalf("wrong number of elements set.\nExpected : 1\nReturned : %d", len(zbxMap.Elements))
	}
}

func BenchmarkAddHosts(b *testing.B) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	host := &hostParameters{
		id:    "router-1",
		name:  "router-1",
		image: "router-1",
		position: &hostPosition{
			spacer: 100,
			mapX:   500,
			mapY:   500,
			x:      100,
			y:      100,
		},
	}

	for i := 0; i < b.N; i++ {
		zbxMap = addHosts(zbxMap, host)
	}
}
