package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestCheckCreateRequiredFlag(t *testing.T) {
	err := checkCreateRequiredFlag("name", mappingFilePath)
	if err != "" {
		t.Fatalf("expected no error to be returned.\nError returned : %s", err)
	}
}

func TestCheckCreateRequiredFlagMissingName(t *testing.T) {
	expectedError := "'name' flag is required and cannot be empty"

	err := checkCreateRequiredFlag("", mappingFilePath)
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func TestCheckCreateRequiredFlagMissingFile(t *testing.T) {
	expectedError := "'file' flag is required and cannot be empty"

	err := checkCreateRequiredFlag("name", "")
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func TestCheckCreateRequiredFlagUnexistingFile(t *testing.T) {
	expectedError := "error while reading file 'file-does-not-exist'."

	err := checkCreateRequiredFlag("name", "file-does-not-exist")
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func BenchmarkCheckCreateRequiredFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkCreateRequiredFlag("name", mappingFilePath)
	}
}

func TestExecuteCreate(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "create")
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B", "--width", "400", "--height", "400")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteCreate")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Add the required environment variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_URL=%s", ZABBIX_URL))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))

	// Run the command in the subprocess
	err := cmd.Run()
	if err != nil {
		exit := err.(*exec.ExitError)
		t.Fatalf("expected exit code 0.\nCode returned : %d\nStdout : %s\nStderr : %s", exit.ExitCode(), stdout.String(), stderr.String())
	}
}

func TestExecuteCreateFailMissingEnvironmentVariable(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "create")
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteCreateFailMissingEnvironmentVariable")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Add the required environment variables (missing ZABBIX_URL on purpose)
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))

	// Run the command in the subprocess
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected an error to be returned, an nil pointer was returned instead")
	}

	exit := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nStdout : %s\nStderr : %s", exit.ExitCode(), stdout.String(), stderr.String())
	}
}

func TestExecuteCreateFail(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "create")
		os.Args = append(os.Args, "--name", generateMapName())
		os.Args = append(os.Args, "--file", mappingFilePath)
		os.Args = append(os.Args, "--color", "7AC2E1", "--trigger-color", "EE445B")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteCreateFail")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Add the required environment variables
	cmd.Env = append(cmd.Env, "ZABBIX_URL=http://localhost:6666/api_jsonrpc.php")
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_USER=%s", ZABBIX_USER))
	cmd.Env = append(cmd.Env, fmt.Sprintf("ZABBIX_PWD=%s", ZABBIX_PWD))

	// Run the command in the subprocess
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected an error to be returned, an nil pointer was returned instead")
	}

	exit := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nStdout : %s\nStderr : %s", exit.ExitCode(), stdout.String(), stderr.String())
	}
}

func TestNewCreateCmd(t *testing.T) {
	cmd := newCreateCmd()
	var expectedType *cobra.Command

	if cmd == nil {
		t.Fatalf("expected *cobra.Command to be returned, a nil pointer was returned instead")
	}

	if reflect.TypeOf(cmd) != reflect.TypeOf(expectedType) {
		t.Fatalf("wrong type returned\nExpected *cobra.Command\nReturned : %s", reflect.TypeOf(cmd))
	}
}

func resetOsEnv() {
	os.Unsetenv("ZABBIX_URL")
	os.Unsetenv("ZABBIX_USER")
	os.Unsetenv("ZABBIX_PWD")
}

func BenchmarkExecuteCreate(b *testing.B) {
	oldOsArgs := os.Args
	oldStdout := os.Stdout

	defer resetOsConf(oldOsArgs, oldStdout)
	defer resetOsEnv()

	os.Stdout = nil
	os.Setenv("ZABBIX_URL", ZABBIX_URL)
	os.Setenv("ZABBIX_USER", ZABBIX_USER)
	os.Setenv("ZABBIX_PWD", ZABBIX_PWD)

	// Run the command in the subprocess
	for i := 0; i < b.N; i++ {
		os.Args = []string{
			os.Args[0],
			"create",
			"--name",
			generateMapName(),
			"--file",
			mappingFilePath,
			"--color",
			"7AC2E1",
			"--trigger-color",
			"EE445B",
		}

		Execute()
	}

}
