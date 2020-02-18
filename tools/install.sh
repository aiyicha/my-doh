#!/bin/bash

HOME="/home/enlink"
MODULE="doh"
SERVICE="${MODULE}.service"
DESC="Enlink DoH (DNS over HTTPs) Module."
CMDPATH=$(dirname $0)

if [ "$1" != "" ] 
 then 
   HOME=$1
fi

mkdir -p ${HOME}

chmod +x ${CMDPATH}/src/bin/*.sh
cp -rf ${CMDPATH}/src/* ${HOME}/

FILE="${HOME}/bin/${SERVICE}"
rm -rf ${FILE}

echo "[Unit]" >> ${FILE}
echo "Description=${DESC}" >> ${FILE}
echo "After=network.target remote-fs.target nss-lookup.target" >> ${FILE}
echo "" >> ${FILE}
echo "[Service]" >> ${FILE}
echo "Type=forking" >> ${FILE}
echo "ExecStart=${HOME}/bin/${MODULE}-start.sh" >> ${FILE}
echo "ExecStop=${HOME}/bin/${MODULE}-stop.sh" >> ${FILE}
echo "" >> ${FILE}
echo "[Install]" >> ${FILE}
echo "WantedBy=multi-user.target" >> ${FILE}

SYSTEMCTL="/usr/lib/systemd/system /lib/systemd/system" 
for folder in ${SYSTEMCTL}; do
	if [ -d ${folder} ]; then
		rm -rf "${folder}/${SERVICE}"
		mv ${FILE} "${folder}/${SERVICE}"
		break
	fi
done

systemctl daemon-reload

