package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

var configStr = `{"TestSvc" : {
	"Image" : "foo/bar",
	"Ports" : {
		"80" : "80",
		"6379" : ""
	},
	"Env" : {
		"HOST" : "webapp.local"
	},
	"Deps" : [
		"Hipache"
	]
}}`

var expected ServiceMap = ServiceMap{"TestSvc": Service{
	Image: "foo/bar",
	Ports: map[string]string{
		"80":   "80",
		"6379": "",
	},
	Env: map[string]string{
		"HOST": "webapp.local",
	},
	Deps: []string{
		"Hipache",
	},
}}

func Test_ServiceMap_parseConfigFile(t *testing.T) {
	// given
	b := []byte(configStr)
	cfgPath := "/tmp/dockit-test-config.json"
	ioutil.WriteFile(cfgPath, b, 0666)

	defer os.Remove(cfgPath)

	expected, err := parseConfigData(b)

	// when
	actual, err := parseConfigFile(cfgPath)

	// then
	if err != nil {
		t.Errorf("Parse Error: %s", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("parseConfigFile returned something unexpected: %+v", actual)
	}
}

func Test_ServiceMap_parseConfigData(t *testing.T) {
	// given
	b := []byte(configStr)

	// when
	actual, err := parseConfigData(b)

	// then
	if err != nil {
		t.Errorf("Parse Error: %s", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("parseConfigData returned something unexpected: %+v", actual)
	}

}
