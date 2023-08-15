package cmd

import (
	"fmt"
	"math/rand"
	"os"
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
	rand.NewSource(time.Now().UnixNano())
	value := rand.Intn(rand.Intn(9999))

	return fmt.Sprintf("test-map-builder-%d", value)
}

// resetOsConf is used to reset the previous os configuration when running certain test (stdout file, os arguments, etc.)
func resetOsConf(args []string, stdout *os.File) {
	os.Args = args
	os.Stdout = stdout
}

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()
	if cmd == nil {
		t.Fatalf("expected a *cobra.Command.\nReturned a nil pointer")
	}
}
