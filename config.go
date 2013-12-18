package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var _ = log.Print // for debugging, remove

type Service struct {
	Image string
	Ports map[string]string
	Deps  []string
	Env   map[string]string
}

type ServiceMap map[string]Service

func parseConfigFile(file string) (ServiceMap, error) {
	var cfg = make(ServiceMap)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	return parseConfigData(b)
}

func parseConfigData(b []byte) (ServiceMap, error) {
	var cfg = make(ServiceMap)

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
