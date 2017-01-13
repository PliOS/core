#!/bin/bash

# Copyright (c) 2017 The pliOS Authors. All rights reserved.
#
# Use of this source code is governed by a MIT-style license that can be found
# in the LICENSE file.

./scripts/pprint.sh "Building" "rootfs"

function create_image() {
    # Create a 32MB ext4 image
    dd if=/dev/zero of=$2/sysroot.img bs=1M count=32 status=none
    mkfs.ext4 -q -F $2/sysroot.img

    # Mount it
    sudo mkdir -p $1
    sudo mount -t ext4 -o loop $2/sysroot.img $1

    # Copy files
    sudo mkdir -p $1/{sbin,bin,proc,sys,dev,run}
    sudo cp build/gopath/bin/init $1/sbin/
    sudo cp rootfs/init.rc.toml $1
    sudo cp $2/bin/busybox $1/bin/busybox
    sudo $1/bin/busybox --install $1/bin

    # Unmount the image
    sudo umount $1
}

create_image $1 $2

./scripts/pprint.sh "Finished" "rootfs"
