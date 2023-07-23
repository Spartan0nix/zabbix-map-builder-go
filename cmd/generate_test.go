package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

const (
	ROUTER    = "172.16.81.161"
	COMMUNITY = "router-1"
)

func TestExecuteGenerate(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Set the required arguments
		os.Args = append(os.Args, "generate")
		os.Args = append(os.Args, "--host", ROUTER, "--community", COMMUNITY, "--port", "1161")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteGenerate")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=1",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command in the subprocess
	err := cmd.Run()
	if err != nil {
		exit := err.(*exec.ExitError)
		if exit.ExitCode() != 0 {
			t.Fatalf("expected exit code 0.\nCode returned : %d\nStdout : %s\nStderr : %s", exit.ExitCode(), stdout.String(), stderr.String())
		}
	}
}

func TestExecuteGenerateFail(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "2" {
		// Set the required arguments
		os.Args = append(os.Args, "generate")
		os.Args = append(os.Args, "--host", "random-host", "--community", COMMUNITY, "--port", "1161")
		Execute()

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteGenerateFail")
	// Reset the subprocess environment variable
	cmd.Env = []string{
		"BE_CRASHER=2",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command in the subprocess
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected an error to be returned when generate process fails")
	}

	exit := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nStdout : %s\nStderr : %s", exit.ExitCode(), stdout.String(), stderr.String())
	}
}

func TestNewGenerateCmd(t *testing.T) {
	cmd := newGenerateCmd()
	var expectedType *cobra.Command

	if cmd == nil {
		t.Fatalf("expected *cobra.Command to be returned, a nil pointer was returned instead")
	}

	if reflect.TypeOf(cmd) != reflect.TypeOf(expectedType) {
		t.Fatalf("wrong type returned\nExpected *cobra.Command\nReturned : %s", reflect.TypeOf(cmd))
	}
}

func TestCheckGenerateRequiredFlag(t *testing.T) {
	err := checkGenerateRequiredFlag(ROUTER, COMMUNITY)
	if err != "" {
		t.Fatalf("expected no error to be returned.\nError returned : %s", err)
	}
}

func TestCheckGenerateRequiredFlagMissingHost(t *testing.T) {
	expectedError := "'host' flag is required and cannot be empty"

	err := checkGenerateRequiredFlag("", mappingFilePath)
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

func TestCheckGenerateRequiredFlagMissingCommunity(t *testing.T) {
	expectedError := "'community' flag is required and cannot be empty"

	err := checkGenerateRequiredFlag(ROUTER, "")
	if err == "" {
		t.Fatalf("expected an error to be returned (none returned).")
	}

	if err != expectedError {
		t.Fatalf("wrong message returned.\nExpected : %s\nReturned : %s", expectedError, err)
	}
}

// func BenchmarkCheckGenerateRequiredFlag(b *testing.B) {
// 	b.ReportAllocs()
// 	for i := 0; i < b.N; i++ {
// 		str := checkGenerateRequiredFlag("host", "community")
// 		if str != "" {
// 			b.Fatalf("an empty string should be returned\nValue returned : %s", str)
// 		}
// 	}
// }
