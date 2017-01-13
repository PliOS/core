// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
)

type ApiFilesystem struct {
	Source string
	Target string
	Fstype string
	Flags  uintptr
	Data   string
	Mode   os.FileMode
}

type ApiSymlink struct {
	Newname string
	Oldname string
}

var ApiFilesystems = [...]ApiFilesystem{
	{
		Source: "proc",
		Target: "/proc",
		Fstype: "proc",
		Flags:  unix.MS_NOSUID | unix.MS_NODEV,
		Data:   "",
		Mode:   0555,
	},
	{
		Source: "sysfs",
		Target: "/sys",
		Fstype: "sysfs",
		Flags:  unix.MS_NOSUID | unix.MS_NOEXEC | unix.MS_NODEV,
		Data:   "",
		Mode:   0555,
	},
	{
		Source: "tmpfs",
		Target: "/sys/fs/cgroup",
		Fstype: "tmpfs",
		Flags:  unix.MS_NOSUID | unix.MS_NOEXEC | unix.MS_NODEV,
		Data:   "mode=0755,size=1M",
		Mode:   0755,
	},
	{
		Source: "cgroup",
		Target: "/sys/fs/cgroup/systemd",
		Fstype: "cgroup",
		Flags:  unix.MS_NOSUID | unix.MS_NOEXEC | unix.MS_NODEV,
		Data:   "name=systemd,none",
		Mode:   0755,
	},
	{
		Source: "devtmpfs",
		Target: "/dev",
		Fstype: "devtmpfs",
		Flags:  unix.MS_NOSUID | unix.MS_STRICTATIME,
		Data:   "mode=0755,size=10M",
		Mode:   0755,
	},
	{
		Source: "devpts",
		Target: "/dev/pts",
		Fstype: "devpts",
		Flags:  unix.MS_NOSUID | unix.MS_STRICTATIME | unix.MS_NOEXEC,
		Data:   "ptmxmode=0666,mode=0620,gid=5,newinstance",
		Mode:   0620,
	},
	{
		Source: "tmpfs",
		Target: "/run",
		Fstype: "tmpfs",
		Flags:  unix.MS_NOSUID | unix.MS_STRICTATIME | unix.MS_NODEV,
		Data:   "mode=0755,size=20%",
		Mode:   0755,
	},
	{
		Source: "tmpfs",
		Target: "/run/shm",
		Fstype: "tmpfs",
		Flags:  unix.MS_NOSUID | unix.MS_STRICTATIME | unix.MS_NOEXEC | unix.MS_NODEV,
		Data:   "mode=01777,size=50%",
		Mode:   01777,
	},
}

var ApiSymlinks = [...]ApiSymlink{
	{"/dev/ptmx", "pts/ptmx"},
	{"/dev/fd", "/proc/self/fd"},
	{"/dev/core", "/proc/kcore"},
	{"/dev/stdin", "/proc/self/fd/0"},
	{"/dev/stdout", "/proc/self/fd/1"},
	{"/dev/stderr", "/proc/self/fd/2"},
	{"/dev/shm", "/run/shm"},
}

func SetupProcessEnvironment() {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_CAD_OFF); err != nil {
		log.Fatalf("Fatal error - reboot(): %s", err)
	}

	if _, err := unix.Setsid(); err != nil {
		log.Fatalf("Fatal error - setsid(): %s", err)
	}

	if err := unix.Chdir("/"); err != nil {
		log.Fatalf("Fatal error - chdir(): %s", err)
	}

	unix.Umask(0000)

	os.Setenv("PATH", "/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin")
	os.Setenv("LANG", "C")
}

func MountApiFilesystems() {
	for _, fs := range ApiFilesystems {
		log.WithFields(log.Fields{
			"source": fs.Source,
			"target": fs.Target,
			"fstype": fs.Fstype,
		}).Debugf("Mounting filesystem")

		if err := os.MkdirAll(fs.Target, fs.Mode); err != nil {
			log.Fatalf("Fatal error - creating directory %s: %s", fs.Target, err)
		}

		if err := unix.Mount(fs.Source, fs.Target, fs.Fstype, fs.Flags, fs.Data); err != nil {
			log.Fatalf("Fatal error - mounting %s: %s", fs.Target, err)
		}
	}
}

func CreateApiSymlinks() {
	for _, symlink := range ApiSymlinks {
		log.WithFields(log.Fields{
			"newname": symlink.Newname,
			"oldname": symlink.Oldname,
		}).Debugf("Creating symlink")

		if _, err := os.Stat(symlink.Newname); err == nil {
			continue
		}

		if err := os.Symlink(symlink.Oldname, symlink.Newname); err != nil {
			log.Fatalf("Fatal error - creating symlink %s to %s: %s", symlink.Newname, symlink.Oldname, err)
		}
	}
}
