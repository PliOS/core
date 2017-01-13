// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	SetupProcessEnvironment()
	MountApiFilesystems()
	CreateApiSymlinks()

	log.Infof("Initalized system")

	config := ReadConfig()

	log.Infof("Read config file")

	grimReaper := NewGrimReaper()
	serviceManager := NewServiceManager(config, grimReaper)
	triggerRunner := NewTriggerRunner(config, grimReaper, serviceManager)

	triggerRunner.RunTrigger("init")

	go grimReaper.Run()
	go serviceManager.Run()
	go ProcessSignals(triggerRunner)

	triggerRunner.Run()
}
