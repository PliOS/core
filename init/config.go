// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	Environment map[string]string
	Triggers    map[string][]string
	Services    map[string]Service
}

type Service struct {
	Program      string
	Arguments    []string
	User         string
	Groups       []string
	Capabilities []string
}

func ReadConfig() *Config {
	data, err := ioutil.ReadFile("init.rc.toml")

	if err != nil {
		log.Fatalf("Fatal error - loading config file: %s", err)
	}

	var config *Config

	if _, err := toml.Decode(string(data), &config); err != nil {
		log.Fatalf("Fatal error - loading config file: %s", err)
	}

	return config
}
