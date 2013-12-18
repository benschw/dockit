package main

import (
	"fmt"
	"github.com/dotcloud/docker"
	"log"
	"os"
	"reflect"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.
var _ = os.Stdout // For debugging; delete when done.

var (
	address = "unix:///var/run/docker.sock"
)

func getCfg() ServiceMap {
	var configStr = `{
		"Hipache" : {
			"Image" : "stackbrew/hipache",
			"Ports" : {
				"80" : "80",
				"6379" : ""
			}
			
		}, "WebApp" : {
			"Image" : "benschw/go-webapp",
			"Deps" : [
				"Hipache"
			],
			"Env" : {
				"HOST" : "webapp.local"
			}
		}
	}`
	cfg, err := parseConfigData([]byte(configStr))

	if err != nil {
		panic(err)
	}
	return cfg
}

func getLib() Lib {
	cfg := getCfg()
	return NewLib(cfg, address, "/tmp/pid-tests-lib")
}

func Test_dockitLib_getPortBindings(t *testing.T) {
	// given
	expected := make(map[docker.Port][]docker.PortBinding)

	portBinding1 := []docker.PortBinding{}
	portBinding1 = []docker.PortBinding{docker.PortBinding{HostIp: "0.0.0.0", HostPort: "80"}}
	port1 := docker.NewPort("tcp", "80")
	expected[port1] = portBinding1

	portBinding2 := []docker.PortBinding{}
	port2 := docker.NewPort("tcp", "6379")
	expected[port2] = portBinding2

	lib := getLib()

	// when
	ports := getCfg()["Hipache"].Ports

	actual := lib.getPortBindings(ports)

	// then
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getPortBinding returned something unexpected")
	}

}

func Test_dockitLib_getEnv(t *testing.T) {
	// given
	expected := []string{"HOST=webapp.local"}

	lib := getLib()

	// when
	env := getCfg()["WebApp"].Env

	actual := lib.getEnv(env)

	// then
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getEnv returned something unexpected")
	}
}
