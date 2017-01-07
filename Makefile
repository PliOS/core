TARGET ?= x86_64-unknown-linux-musl
ARCH ?= x86_64

.PHONY: init optimize sysroot run

all: init busybox optimize sysroot

init:
	@mkdir -p build/bin
	@cd init && cargo build --target $(TARGET)
	@cp init/target/$(TARGET)/debug/plios_init build/bin/init

busybox: build/bin/busybox

build/bin/busybox:
	@cd build/bin && wget -O busybox https://busybox.net/downloads/binaries/1.26.1-defconfig-multiarch/busybox-$(ARCH)
	@chmod a+x build/bin/busybox

optimize: init busybox
	@strip --strip-debug build/bin/init
	@strip --strip-debug build/bin/busybox

sysroot: optimize
	@dd if=/dev/zero of=build/sysroot.img bs=1M count=32
	@mkfs.ext4 -F build/sysroot.img
	@sudo mkdir -p /media/plios_sysroot
	@sudo mount -t ext4 -o loop build/sysroot.img /media/plios_sysroot
	@sudo mkdir -p /media/plios_sysroot/sbin
	@sudo mkdir -p /media/plios_sysroot/bin
	@sudo cp build/bin/init /media/plios_sysroot/sbin/init
	@sudo cp build/bin/busybox /media/plios_sysroot/bin/busybox
	@sudo ./scripts/create_busybox_symlinks.sh
	@sudo umount /media/plios_sysroot

run:
	@qemu-system-x86_64 -kernel build/vmlinuz-4.8.0-32-generic -append "root=/dev/sda rw vga=0x344 \
	                    devtmpfs.mount=0" -drive file=build/sysroot.img,format=raw,if=ide
