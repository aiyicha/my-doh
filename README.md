# doh
DNS over HTTPs module (open source ```m13253/dns-over-https```, version 1.4.3). The newest version is not compatible with ```chromium```.

## Build Dependencies
- GNU make tool 'make':   
The 'make' tool is installed by default in CentOS, if not, please run ```sudo yum install make``` to install it.
- Golang command 'go':   
See project http://192.168.100.70/tools/golang for more information.

## Build Packages
Run command
```shell script
make
```
to build tarball (like ```doh-1.4.3-release.tar.gz```).

## Install
Use following commands to install ```doh``` module (use ```1.4.3-release``` version as example):
```shell script
tar zxvf doh-1.4.3-release.tar.gz
cd doh-1.4.3-release
sh install.sh <path>
```
Argument ```<path>``` in the command is optional, default ```path``` is ```/home/enlink```.

## Run
Use following command to run ```doh``` module:
```shell script
service doh start
```

## Notes
Doh module provide HTTP service only, please learn to run a ```nginx``` server before ```doh``` server to provide HTTPs service.