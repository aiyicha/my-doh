#!/bin/bash

PIDFILE="/var/run/doh-server.pid"
if [ ! -f ${PIDFILE} ]; then exit 0; fi

PID=`cat ${PIDFILE}`
if [ "${PID}" != "" ]; then
	if [ -d "/proc/${PID}" ]; then
		kill -9 ${PID}
	fi
fi

rm -rf ${PIDFILE}
