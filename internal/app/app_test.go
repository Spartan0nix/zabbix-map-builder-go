package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// generateMapName is used to generate a random name for each map created during test.
func generateMapName() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(rand.Intn(9999))

	return fmt.Sprintf("test-map-builder-%d", value)
}

func TestOutputToFile(t *testing.T) {
	file := "test-output-file.json"

	err := outputToFile(file, &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	})

	if err != nil {
		t.Fatalf("error while execution outputToFile function.\nReason : %v", err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("error while retrieving information about the file '%s'.\nReason : %v", file, err)
	}

	if info == nil {
		t.Fatalf("no info was retrieved about the file '%s'", file)
	}

	err = os.Remove(file)
	if err != nil {
		t.Fatalf("error while removing file '%s'.\nReason : %v", file, err)
	}
}

func TestOutputToFileEmptyName(t *testing.T) {
	err := outputToFile("", &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	})

	if err == nil {
		t.Fatal("an error should be returned when an empty file name is passed to the outputToFile function")
	}
}

func TestRunApp(t *testing.T) {
	opts := Options{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
	}

	err := RunApp(mappingFilePath, &opts)
	if err != nil {
		t.Fatalf("error while executing RunApp function.\nReason : %v", err)
	}
}

func TestRunAppDryRun(t *testing.T) {
	// Keep the previous stdout file
	oldStdout := os.Stdout
	// Create a new read (r) and write (w) pipe file
	r, w, _ := os.Pipe()
	// Switch to the new out file
	os.Stdout = w

	// Set the required arguments
	opts := Options{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
		DryRun:       true,
	}

	// Create a new channel
	outChannel := make(chan string)
	// Copy the stdout content in a goroutine
	go func() {
		// Write content from the file to a new buffer
		var buf bytes.Buffer
		io.Copy(&buf, r)
		// Write the buffer to the channel
		outChannel <- buf.String()
	}()

	err := RunApp(mappingFilePath, &opts)
	if err != nil {
		t.Fatalf("error while executing RunApp function.\nReason : %v", err)
	}

	// Close the pipe
	w.Close()
	// Restore the previous stdout
	os.Stdout = oldStdout
	// Sotre the content of the channel to a variable
	out := <-outChannel

	if out == "" {
		log.Fatalf("expected the content of the created map to be exposed in the current shell (an empty string was returned)")
	}
}

func TestRunAppOutFile(t *testing.T) {
	outFile := "test-map-builder-output"

	// Set the required arguments
	opts := Options{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
		OutFile:      outFile,
	}

	err := RunApp(mappingFilePath, &opts)
	if err != nil {
		t.Fatalf("error while executing RunApp function.\nReason : %v", err)
	}

	b, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("error while reading output file '%s'.\nReason : %v", outFile, err)
	}

	if string(b) == "" {
		t.Fatalf("file '%s' does not contains the request used to create the map", outFile)
	}

	if err = os.Remove(outFile); err != nil {
		t.Fatalf("error while removing output file '%s'.\nReason : %v", outFile, err)
	}
}

func TestRunAppFailReadInput(t *testing.T) {
	// Set the required arguments
	opts := Options{}

	err := RunApp("file-does-not-exist", &opts)
	if err == nil {
		t.Fatalf("an error should be returned when the host mapping file does not exist")
	}
}
func TestRunAppFailInitApi(t *testing.T) {
	// Set the required arguments
	opts := Options{
		ZabbixUrl:    "http://localhost:6666/api_jsonrpc.php",
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         "test-map-builder_should-not-exist",
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
	}

	err := RunApp(mappingFilePath, &opts)
	if err == nil {
		t.Fatalf("an error should be returned when the Zabbix API is unreachable")
	}
}
