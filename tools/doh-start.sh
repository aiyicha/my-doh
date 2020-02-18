#!/bin/bash

CMDPATH=$(dirname $0)

CMD="${CMDPATH}/doh-server"
CMD="${CMD} --conf ${CMDPATH}/../etc/doh-server.conf"
CMD="${CMD} --hosts ${CMDPATH}/../etc/doh-hosts.conf"
CMD="${CMD} --pid-file /var/run/doh-server.pid"

${CMD} &
