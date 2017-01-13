// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/sys/unix"
	"sync"
)

type ServiceManager struct {
	config          *Config
	reaper          *GrimReaper
	runningServices map[string]int
	servicePids     map[int]string
	serviceLock     sync.RWMutex
}

func NewServiceManager(config *Config, reaper *GrimReaper) *ServiceManager {
	serviceManager := new(ServiceManager)
	serviceManager.config = config
	serviceManager.reaper = reaper

	serviceManager.runningServices = map[string]int{}
	serviceManager.servicePids = map[int]string{}

	return serviceManager
}

func (self *ServiceManager) Start(name string) {
	self.serviceLock.Lock()

	if _, running := self.runningServices[name]; !running {
		self.runningServices[name] = RunCommand(
			self.config.Services[name].Program,
			self.config.Services[name].Arguments,
		)

		self.servicePids[self.runningServices[name]] = name

		log.WithFields(log.Fields{
			"name": name,
			"pid":  self.runningServices[name],
		}).Debugf("Started service")
	}

	self.serviceLock.Unlock()
}

func (self *ServiceManager) Stop(name string) int {
	self.serviceLock.Lock()

	if pid, running := self.runningServices[name]; running {
		unix.Kill(pid, unix.SIGTERM)
		delete(self.runningServices, name)

		log.WithFields(log.Fields{
			"name": name,
			"pid":  pid,
		}).Debugf("Stopped service")

		return pid
	}

	self.serviceLock.Unlock()

	return 0
}

func (self *ServiceManager) Restart(name string) {
	self.serviceLock.Lock()

	if pid, running := self.runningServices[name]; running {
		unix.Kill(pid, unix.SIGTERM)

		log.WithFields(log.Fields{
			"name": name,
			"pid":  pid,
		}).Debugf("Restarting service")
	}

	self.serviceLock.Unlock()
}

func (self *ServiceManager) respawn(service string, oldPid int) {
	newPid := RunCommand(
		self.config.Services[service].Program,
		self.config.Services[service].Arguments,
	)

	self.runningServices[service] = newPid
	self.servicePids[newPid] = service

	delete(self.servicePids, oldPid)

	log.WithFields(log.Fields{
		"name":   service,
		"oldPid": oldPid,
		"newPid": newPid,
	}).Debugf("Respawned service")
}

func (self *ServiceManager) Run() {
	deadProcesses := self.reaper.WaitWildcard()

	for deadProcess := range deadProcesses {
		pid := deadProcess.Pid

		log.WithFields(log.Fields{
			"pid": pid,
		}).Debugf("Process died")

		self.serviceLock.Lock()

		if service, running := self.servicePids[pid]; running {
			if _, active := self.runningServices[service]; active {
				self.respawn(service, pid)
			} else {
				delete(self.servicePids, pid)

				log.WithFields(log.Fields{
					"name": service,
					"pid":  pid,
				}).Debugf("Service died")
			}
		}

		self.serviceLock.Unlock()
	}
}
