#!/bin/bash

STOP="$(dirname $0)/doh-stop.sh"
START="$(dirname $0)/doh-start.sh"

sh ${STOP}
sh ${START}
