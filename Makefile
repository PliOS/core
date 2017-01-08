# Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
#
# Use of this source code is governed by a MIT-style license that can be found
# in the LICENSE file.

TARGET ?= x86_64-unknown-linux-musl
ARCH ?= x86_64

BUSYBOX_URL ?= https://busybox.net/downloads/binaries/1.26.1-defconfig-multiarch/busybox-$(ARCH)

QEMU_VMLINUZ ?=  build/vmlinuz-4.8.0-32-generic
QEMU_LINUX_ARGS ?= root=/dev/sda vga=0x344 devtmpfs.mount=0
QEMU_DRIVE ?= file=build/sysroot.img,format=raw,if=ide

QEMU_ARGS := -kernel $(QEMU_VMLINUZ) -append "$(QEMU_LINUX_ARGS)" -drive $(QEMU_DRIVE)

.PHONY: init optimize sysroot run

all: init busybox optimize sysroot

init:
	@mkdir -p build/bin
	@cd init && cargo build --target $(TARGET)
	@cp init/target/$(TARGET)/debug/plios_init build/bin/init

busybox: build/bin/busybox

build/bin/busybox:
	@cd build/bin && wget -O busybox $(BUSYBOX_URL)
	@chmod a+x build/bin/busybox

optimize: init busybox
	@strip --strip-debug build/bin/init
	@strip --strip-debug build/bin/busybox

sysroot: optimize
	@./scripts/create_rootfs.sh /media/plios_sysroot build/

run:
	@qemu-system-x86_64 $(QEMU_ARGS)
