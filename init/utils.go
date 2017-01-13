// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

func RunCommand(program string, arguments []string) int {
	cmd := exec.Command(program, arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Fatal error - exec(%s, %v): %s", program, arguments, err)
	}

	return cmd.Process.Pid
}
