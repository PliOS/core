# Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
#
# Use of this source code is governed by a MIT-style license that can be found
# in the LICENSE file.

ARCH ?= x86_64

BUSYBOX_URL ?= https://busybox.net/downloads/binaries/1.26.1-defconfig-multiarch/busybox-$(ARCH)

QEMU_VMLINUZ ?=  build/vmlinuz-4.8.0-32-generic
QEMU_LINUX_ARGS ?= root=/dev/sda vga=0x344 devtmpfs.mount=0
QEMU_DRIVE ?= file=build/sysroot.img,format=raw,if=ide

QEMU_ARGS := -kernel $(QEMU_VMLINUZ) -append "$(QEMU_LINUX_ARGS)" -drive $(QEMU_DRIVE)

export GOPATH := $(shell pwd)/build/gopath
export GOARCH := amd64
export GOOS := linux
export CGO_ENABLED := 0

.PHONY: init optimize sysroot run

all: sysroot

sysroot: init busybox
	@./scripts/create_rootfs.sh /media/plios_sysroot build/

init: gopath
	@go get github.com/PliOS/core/init
	@go install github.com/PliOS/core/init

gopath: build/gopath/

build/gopath/:
	@mkdir -p build/gopath/src/github.com/PliOS/core
	@mkdir -p build/gopath/bin
	@mkdir -p build/gopath/pkg
	@ln -s $(shell pwd)/init $(shell pwd)/build/gopath/src/github.com/PliOS/core

busybox: build/bin/busybox

build/bin/busybox:
	@cd build/bin && wget -O busybox $(BUSYBOX_URL)
	@chmod a+x build/bin/busybox

run:
	@qemu-system-x86_64 $(QEMU_ARGS)
