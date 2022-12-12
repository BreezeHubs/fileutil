package main

import (
	"fileutil/yaml"
	"testing"
)

func Test_yaml(t *testing.T) {
	config, _ := yaml.LoadConfig("./config.yaml")
	getString := config.GetString("devices[0].nodes[0].index")
	getInt := config.GetInt("ipport")
	t.Log(getString, getInt)
}
