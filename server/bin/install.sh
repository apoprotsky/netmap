#!/bin/sh

# Prefix
PREFIX="/usr/local"

# Platform
PLATFORM=`uname`

# Location
if [ -L ${0} ]; then
    tmp=`readlink ${0}`
else
    tmp=${0}
fi
tmp=`dirname ${tmp}`
LOCATION=`cd ${tmp} && pwd`

# Build
cd ${LOCATION}/../
go get -d
go build

# Binary
install -o 0 -g 0 -m u=rwx,go=rx ${LOCATION}/../server ${PREFIX}/bin/netmap

# RC script
if [ ${PLATFORM} = "FreeBSD" ]; then
    install -o 0 -g 0 -m u=rw,go=r ${LOCATION}/../dist/freebsd/netmap ${PREFIX}/etc/rc.d/netmap
fi

# Configuration files
install -o 0 -g 0 -m u=rw,go=r ${LOCATION}/../dist/netmap.ini.sample ${PREFIX}/etc/netmap.ini.sample
[ ! -f ${PREFIX}/etc/netmap.ini ] && install -o 0 -g 0 -m u=rw,go=r ${LOCATION}/../dist/netmap.ini.sample ${PREFIX}/etc/netmap.ini
