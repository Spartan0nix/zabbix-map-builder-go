package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	ZABBIX_URL  = "http://localhost:4444/api_jsonrpc.php"
	ZABBIX_USER = "Admin"
	ZABBIX_PWD  = "zabbix"
)

var mappingFilePath string

func init() {
	pwd, _ := os.Getwd()
	mappingFilePath = filepath.Join(pwd, "..", "examples", "mapping.json")
}

// generateMapName is used to generate a random name for each map created during test.
func generateMapName() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(rand.Intn(9999))

	return fmt.Sprintf("test-map-builder-%d", value)
}

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()
	if cmd == nil {
		t.Fatalf("expected a *cobra.Command.\nReturned a nil pointer")
	}
}

func TestExecute(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecute")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}
	// Add the required environment variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_URL=%s", ZABBIX_URL))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))
	// Read the output of the command
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command in the subprocess
	err := cmd.Run()

	if err != nil {
		exit := err.(*exec.ExitError)
		t.Fatalf("expected exit code 0.\nCode returned : %d\nError returned : %s", exit.ExitCode(), string(exit.Stderr))
	}
}

func TestExecuteFailMissingEnvironmentVariable(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteFailMissingEnvironmentVariable")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}
	// Add the required environment variables (missing ZABBIX_URL on purpose)
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))
	// --------------------------------------------
	// Should only be enabled for debug purpose
	// --------------------------------------------
	// Read the output of the command
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// --------------------------------------------
	// Run the command in the subprocess
	err := cmd.Run()

	if err == nil {
		t.Fatalf("expected an error to be returned, an nil pointer was returned instead")
	}

	exit := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nError returned : %s", exit.ExitCode(), string(exit.Stderr))
	}
}

func TestExecuteFail(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteFail")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}
	// Add the required environment variables
	cmd.Env = append(cmd.Env, "ZABBIX_URL=http://localhost:6666/api_jsonrpc.php")
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))
	// --------------------------------------------
	// Should only be enabled for debug purpose
	// --------------------------------------------
	// Read the output of the command
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// --------------------------------------------
	// Run the command in the subprocess
	err := cmd.Run()

	if err == nil {
		t.Fatalf("expected an error to be returned, an nil pointer was returned instead")
	}

	exit := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nError returned : %s", exit.ExitCode(), string(exit.Stderr))
	}
}

func TestCheckRequiredFlag(t *testing.T) {
	err := checkRequiredFlag("name", mappingFilePath)
	if err != "" {
		t.Fatalf("expected no error to be returned.\nError returned : %s", err)
	}
}

func TestCheckRequiredFlagMissingName(t *testing.T) {
	expectedError := "'name' flag is required and cannot be empty"

	err := checkRequiredFlag("", mappingFilePath)
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func TestCheckRequiredFlagMissingFile(t *testing.T) {
	expectedError := "'file' flag is required and cannot be empty"

	err := checkRequiredFlag("name", "")
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func TestCheckRequiredFlagUnexistingFile(t *testing.T) {
	expectedError := "error while reading file 'file-does-not-exist'."

	err := checkRequiredFlag("name", "file-does-not-exist")
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}
