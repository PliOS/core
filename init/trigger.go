// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
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
)

func ProcessTriggers(config *Config, triggers chan string, execCommands chan ExecCommand, execFinished chan int, serviceActions chan ServiceAction) {
	selfpipe := make(chan string, 128)

	go SelfPipeTriggers(triggers, selfpipe)

	for trigger := range triggers {
		log.WithFields(log.Fields{
			"trigger": trigger,
		}).Infof("Received trigger")

		for _, action := range config.Triggers[trigger] {
			log.WithFields(log.Fields{
				"command": action,
			}).Debugf("Executing command")

			ProcessAction(action, selfpipe, execCommands, execFinished, serviceActions)
		}
	}
}

func SelfPipeTriggers(triggers chan string, selfpipe chan string) {
	for trigger := range selfpipe {
		triggers <- trigger
	}
}

func ProcessAction(commandString string, selfpipe chan string, execCommands chan ExecCommand, execFinished chan int, serviceActions chan ServiceAction) {
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

		command := ExecCommand{
			Program:   command[1],
			Arguments: command[2:],
		}

		execCommands <- command
		<-execFinished
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
	case "restart", "start", "stop":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		service := ServiceAction{
			Service: command[1],
			Action:  command[0],
		}

		serviceActions <- service
	case "trigger":
		if len(command) != 2 {
			log.Fatalf("Fatal error - invalid command: %s", commandString)
		}

		selfpipe <- command[1]
	default:
		log.Fatalf("Fatal error - invalid command: %s", commandString)
	}
}
