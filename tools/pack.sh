#!/bin/sh

NAME="doh"
SRC_BIN="bin"
SRC_ETC="etc"
SRC_SRC="src/github.com/m13253/dns-over-https"
SRC_TOOLS="tools"

if [ "$1" != "" ]; then NAME=$1; fi
if [ ! -x "${SRC_BIN}/${NAME}-server" ]; then
	echo "ERROR: Binary file '${SRC_BIN}/${NAME}-server' not found."
	exit 2
fi

PACK="${NAME}-`${SRC_BIN}/${NAME}-server -version`-release"
DST_SRC="${PACK}/src"
DST_BIN="${DST_SRC}/bin"
DST_ETC="${DST_SRC}/etc"

printf "Generating ${PACK}.tar.gz ..."
	mkdir -p ${PACK}
	mkdir -p ${DST_SRC}
	mkdir -p ${DST_BIN}
	mkdir -p ${DST_ETC}

	cp -rf ${SRC_BIN}/*               ${DST_BIN}
	cp -rf ${SRC_TOOLS}/install.sh    ${PACK}
	cp -rf ${SRC_TOOLS}/${NAME}-*.sh  ${DST_BIN}

	cp -rf ${SRC_SRC}/doh-client/doh-client.conf ${DST_ETC}
	cp -rf ${SRC_SRC}/doh-server/doh-server.conf ${DST_ETC}
	cp -rf ${SRC_SRC}/doh-server/doh-hosts.conf  ${DST_ETC}

	rm -rf ${PACK}.tar.gz
	tar zcvf ${PACK}.tar.gz ${PACK}/* > /dev/null
	rm -rf ${PACK}
printf "\rGenerating ${PACK}.tar.gz ... [OK]\n"
