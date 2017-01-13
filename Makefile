# Copyright (c) 2017 The pliOS Authors. All rights reserved.
#
# Use of this source code is governed by a MIT-style license that can be found
# in the LICENSE file.

ARCH ?= x86_64

BUSYBOX_URL ?= https://busybox.net/downloads/binaries/1.26.1-defconfig-multiarch/busybox-$(ARCH)

QEMU_VMLINUZ ?=  build/vmlinuz
QEMU_LINUX_ARGS ?= root=/dev/sda vga=0x344 devtmpfs.mount=0
QEMU_DRIVE ?= file=build/sysroot.img,format=raw,if=ide

QEMU_ARGS := -kernel $(QEMU_VMLINUZ) -append "$(QEMU_LINUX_ARGS)" -drive $(QEMU_DRIVE)

export GOPATH := $(shell pwd)/build/gopath
export GOARCH := amd64
export GOOS := linux
export CGO_ENABLED := 0

GOPATH_LOC := $(GOPATH)/src/github.com/PliOS/core/

.PHONY: sysroot init service_manager gopath run fmt

all: sysroot

sysroot: init service_manager busybox
	@./scripts/create_rootfs.sh /media/plios_sysroot build/

init: gopath
	@./scripts/pprint.sh "Downloading dependencies for" "init"
	@go get github.com/PliOS/core/init
	@./scripts/pprint.sh "Building" "init"
	@go install github.com/PliOS/core/init

gopath: $(GOPATH_LOC) $(GOPATH_LOC)init

$(GOPATH_LOC):
	@mkdir -p $(GOPATH_LOC)
	@mkdir -p build/gopath/bin
	@mkdir -p build/gopath/pkg

$(GOPATH_LOC)init:
	ln -s $(shell pwd)/init $(GOPATH_LOC)

busybox: build/bin/busybox

build/bin/busybox:
	@cd build/bin && wget -O busybox $(BUSYBOX_URL)
	@chmod a+x build/bin/busybox

run:
	@qemu-system-x86_64 $(QEMU_ARGS)

fmt: gopath
	@./scripts/pprint.sh "Formatting" "init"
	@cd $(GOPATH_LOC)/init && go fmt

