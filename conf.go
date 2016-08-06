package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ConfFile is the path to the configuration file
var ConfFile = flag.String("conf", "conf.yaml", "Path to the conf file")

// Conf stores the configuration
var Conf Options

// Options is the structure of the config file
type Options struct {
	ServerPort        string `yaml:"ServerPort"`
	ImaginaryHostPort string `yaml:"ImaginaryHostPort"`
	CacheRoot         string `yaml:"CacheRoot"`
}

func init() {
	flag.Parse()
	b, err := ioutil.ReadFile(*ConfFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
