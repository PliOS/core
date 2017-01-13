// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type TriggerRunner struct {
	config         *Config
	reaper         *GrimReaper
	serviceManager *ServiceManager
	triggerChan    chan string
}

func NewTriggerRunner(config *Config, reaper *GrimReaper, serviceManager *ServiceManager) *TriggerRunner {
	triggerRunner := new(TriggerRunner)
	triggerRunner.config = config
	triggerRunner.reaper = reaper
	triggerRunner.serviceManager = serviceManager
	triggerRunner.triggerChan = make(chan string, 128)

	return triggerRunner
}

func (self *TriggerRunner) RunTrigger(trigger string) {
	self.triggerChan <- trigger
}

func (self *TriggerRunner) Run() {
	for trigger := range self.triggerChan {
		log.WithFields(log.Fields{
			"trigger": trigger,
		}).Infof("Received trigger")

		for _, action := range self.config.Triggers[trigger] {
			log.WithFields(log.Fields{
				"command": action,
			}).Debugf("Executing command")

			self.ProcessAction(action)
		}
	}
}

func (self *TriggerRunner) ProcessAction(commandString string) {
	command := strings.Split(commandString, " ")

	if len(command) < 1 {
		log.Fatalf("Fatal error - invalid command: %s", commandString)
	}

	switch command[0] {
	case "chmod":
		if len(command) != 3 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		mode, err := strconv.ParseInt(command[1], 8, 64)

		if err != nil {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[2]

		if err := os.Chmod(path, os.FileMode(mode)); err != nil {
			log.Fatalf("Fatal error - chmod(%s, %o): %s", path, mode, err)
		}
	case "chown":
		if len(command) != 4 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		owner, err := strconv.Atoi(command[1])

		if err != nil {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		group, err := strconv.Atoi(command[2])

		if err != nil {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[3]

		if err := os.Chown(path, owner, group); err != nil {
			log.Fatalf("Fatal error - chown(%s, %d, %d): %s", path, owner, group, err)
		}
	case "exec":
		if len(command) < 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		pid := RunCommand(command[1], command[2:])

		<-self.reaper.WaitPid(pid)
	case "mount":
		if len(command) < 4 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		fstype := command[1]
		source := command[2]
		target := command[3]
		data := ""

		var flags uintptr

		flags = 0

		for _, tflag := range command[4:] {
			switch tflag {
			case "active":
				flags |= unix.MS_ACTIVE
			case "async":
				flags |= unix.MS_ASYNC
			case "bind":
				flags |= unix.MS_BIND
			case "dirsync":
				flags |= unix.MS_DIRSYNC
			case "invalidate":
				flags |= unix.MS_INVALIDATE
			case "i_version":
				flags |= unix.MS_I_VERSION
			case "kernmount":
				flags |= unix.MS_KERNMOUNT
			case "mandlock":
				flags |= unix.MS_MANDLOCK
			case "mgc_msk":
				flags |= unix.MS_MGC_MSK
			case "mgc_val":
				flags |= unix.MS_MGC_VAL
			case "move":
				flags |= unix.MS_MOVE
			case "noatime":
				flags |= unix.MS_NOATIME
			case "nodev":
				flags |= unix.MS_NODEV
			case "nodiratime":
				flags |= unix.MS_NODIRATIME
			case "noexec":
				flags |= unix.MS_NOEXEC
			case "nosuid":
				flags |= unix.MS_NOSUID
			case "posixacl":
				flags |= unix.MS_POSIXACL
			case "private":
				flags |= unix.MS_PRIVATE
			case "ro":
				flags |= unix.MS_RDONLY
			case "rec":
				flags |= unix.MS_REC
			case "relatime":
				flags |= unix.MS_RELATIME
			case "remount":
				flags |= unix.MS_REMOUNT
			case "rmt_mask":
				flags |= unix.MS_RMT_MASK
			case "shared":
				flags |= unix.MS_SHARED
			case "silent":
				flags |= unix.MS_SILENT
			case "slave":
				flags |= unix.MS_SLAVE
			case "strictatime":
				flags |= unix.MS_STRICTATIME
			case "sync":
				flags |= unix.MS_SYNC
			case "synchronous":
				flags |= unix.MS_SYNCHRONOUS
			case "unbindable":
				flags |= unix.MS_UNBINDABLE
			default:
				data = tflag
			}
		}

		if err := unix.Mount(source, target, fstype, flags, data); err != nil {
			log.Fatalf("Fatal error - mount(%s, %s, %s, %x, %d): %s", source, target, fstype, flags, data, err)
		}
	case "umount":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[1]

		if err := unix.Unmount(path, 0); err != nil {
			log.Fatalf("Fatal error - unmount(%s): %s", path, err)
		}
	case "rm":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[1]

		if err := unix.Unlink(path); err != nil {
			log.Fatalf("Fatal error - unlink(%s): %s", path, err)
		}
	case "write":
		if len(command) != 3 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[1]
		data := command[2]

		if err := ioutil.WriteFile(path, []byte(data), 0644); err != nil {
			log.Fatalf("Fatal error - write(%s, %s): %s", path, data, err)
		}

	case "mkdir":
		if len(command) != 3 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[1]

		mode, err := strconv.ParseInt(command[2], 8, 64)

		if err != nil {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		if err := os.MkdirAll(path, os.FileMode(mode)); err != nil {
			log.Fatalf("Fatal error - mkdir(%s, %o): %s", path, mode, err)
		}
	case "rmdir":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		path := command[1]

		if err := unix.Rmdir(path); err != nil {
			log.Fatalf("Fatal error - rmdir(%s): %s", path, err)
		}
	case "start":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		self.serviceManager.Start(command[1])
	case "restart":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		self.serviceManager.Restart(command[1])
	case "stop":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		self.serviceManager.Stop(command[1])
	case "stopwaitkill":
		if len(command) != 3 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		timeout, err := time.ParseDuration(command[2])

		if err != nil {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		pid := self.serviceManager.Stop(command[1])

		select {
		case <-self.reaper.WaitPid(pid):
		case <-time.After(timeout):
			unix.Kill(pid, unix.SIGKILL)
		}
	case "trigger":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		self.triggerChan <- command[1]
	case "reboot":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		switch command[1] {
		case "halt":
			unix.Reboot(unix.LINUX_REBOOT_CMD_HALT)
		case "shutdown":
			unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
		case "restart":
			unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
		default:
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}
	default:
		log.Fatalf("Fatal error - invalid command: %s", commandString)
	}
}
