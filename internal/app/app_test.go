package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

const (
	router          = "172.16.81.161"
	router2         = "172.16.81.162"
	port            = uint16(1161)
	community       = "router-1"
	community2      = "router-2"
	generateOutFile = "generated_test_file.json"
)

// generateMapName is used to generate a random name for each map created during test.
func generateMapName() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(rand.Intn(9999))

	return fmt.Sprintf("test-map-builder-%d", value)
}

func TestOutputToFile(t *testing.T) {
	file := "test-output-file.json"
	p := &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("error while marshaling data.\nReason : %v", err)
	}

	err = outputToFile(file, b)
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
	b := make([]byte, 0)

	err := outputToFile("", b)
	if err == nil {
		t.Fatal("an error should be returned when an empty file name is passed to the outputToFile function")
	}
}

func TestRunCreate(t *testing.T) {
	opts := MapOptions{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
		Width:        "400",
		Height:       "400",
		Spacer:       50,
	}

	err := RunCreate(mappingFilePath, &opts, nil)
	if err != nil {
		t.Fatalf("error while executing RunCreate function.\nReason : %v", err)
	}
}

func TestRunCreateDryRun(t *testing.T) {
	// Keep the previous stdout file
	oldStdout := os.Stdout
	// Create a new read (r) and write (w) pipe file
	r, w, _ := os.Pipe()
	// Switch to the new out file
	os.Stdout = w

	// Set the required arguments
	opts := MapOptions{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
		Height:       "800",
		Width:        "800",
		Spacer:       50,
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

	err := RunCreate(mappingFilePath, &opts, nil)
	if err != nil {
		t.Fatalf("error while executing RunCreate function.\nReason : %v", err)
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

func TestRunCreateOutFile(t *testing.T) {
	outFile := "test-map-builder-output"

	// Set the required arguments
	opts := MapOptions{
		ZabbixUrl:    ZABBIX_URL,
		ZabbixUser:   ZABBIX_USER,
		ZabbixPwd:    ZABBIX_PWD,
		Name:         generateMapName(),
		Color:        "7AC2E1",
		TriggerColor: "EE445B",
		Height:       "800",
		Width:        "800",
		Spacer:       50,
		OutFile:      outFile,
	}

	err := RunCreate(mappingFilePath, &opts, nil)
	if err != nil {
		t.Fatalf("error while executing RunCreate function.\nReason : %v", err)
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

func TestRunCreateFailReadInput(t *testing.T) {
	// Set the required arguments
	opts := MapOptions{}

	err := RunCreate("file-does-not-exist", &opts, nil)
	if err == nil {
		t.Fatalf("an error should be returned when the host mapping file does not exist")
	}
}
func TestRunCreateFailInitApi(t *testing.T) {
	// Set the required arguments
	opts := MapOptions{
		ZabbixUrl:  "http://localhost:6666/api_jsonrpc.php",
		ZabbixUser: ZABBIX_USER,
		ZabbixPwd:  ZABBIX_PWD,
	}

	err := RunCreate(mappingFilePath, &opts, nil)
	if err == nil {
		t.Fatalf("an error should be returned when the Zabbix API is unreachable")
	}
}

func TestRunGenerate(t *testing.T) {
	// Keep the previous stdout file
	oldStdout := os.Stdout
	// Create a new read (r) and write (w) pipe file
	r, w, _ := os.Pipe()
	// Switch to the new out file
	os.Stdout = w

	// Set the required arguments
	opts := GenerateOptions{
		Host:           router,
		Port:           port,
		Community:      community,
		TriggerPattern: "Interface #INTERFACE down",
		LocalImage:     "Switch_(64)",
		RemoteImage:    "Switch_(64)",
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

	err := RunGenerate(&opts, nil)
	if err != nil {
		t.Fatalf("error while executing RunGenerate function.\nReason : %v", err)
	}

	// Close the pipe
	w.Close()
	// Restore the previous stdout
	os.Stdout = oldStdout
	// Sotre the content of the channel to a variable
	out := <-outChannel

	if out == "" {
		log.Fatalf("expected the content of the generated mappings to be exposed in the current shell (an empty string was returned)")
	}
}

func TestRunGenerateOutFile(t *testing.T) {

	// Set the required arguments
	opts := GenerateOptions{
		Host:           router,
		Port:           port,
		Community:      community,
		TriggerPattern: "Interface #INTERFACE down",
		LocalImage:     "Switch_(64)",
		RemoteImage:    "Switch_(64)",
		OutFile:        generateOutFile,
	}

	err := RunGenerate(&opts, nil)
	if err != nil {
		t.Fatalf("error while executing RunGenerate function\nReason : %v", err)
	}

	_, err = os.Stat(generateOutFile)
	if err != nil {
		t.Fatalf("error while retrieving test output file ('%s') stats\nReason : %v", generateOutFile, err)
	}

	err = os.Remove(generateOutFile)
	if err != nil {
		t.Fatalf("error while removing test output file ('%s')\nReason : %v", generateOutFile, err)
	}
}

func TestRunGenerateFail(t *testing.T) {
	expectedError := fmt.Sprintf("no cdp data found on host '%s', check if cdp is up and running on the host", router2)

	// Set the required arguments
	opts := GenerateOptions{
		Host:           router2,
		Port:           port,
		Community:      community2,
		TriggerPattern: "Interface #INTERFACE down",
		LocalImage:     "Switch_(64)",
		RemoteImage:    "Switch_(64)",
	}

	err := RunGenerate(&opts, nil)
	if err == nil {
		t.Fatalf("expected an error to be returned when the host is unreachable")
	}

	if err.Error() != expectedError {
		t.Fatalf("wrong error format returned\nExpected : %s\nReturned : %s", expectedError, err.Error())
	}
}

func BenchmarkOutputToFile(b *testing.B) {
	data := []byte("random-test")

	for i := 0; i < b.N; i++ {
		fileName := fmt.Sprintf("benchmark-output-file-%d", i)

		outputToFile(fileName, data)
		os.Remove(fileName)
	}
}
