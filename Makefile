export GOPATH=$(shell pwd)

SRC_PATH="github.com/m13253/dns-over-https"
MODULE=doh

a all: release

release: server client
	@sh tools/pack.sh $(MODULE) release

client:
	@printf "Building $(MODULE)-client ..."
	@go install -i $(SRC_PATH)/$(MODULE)-client
	@printf "\rBuilding $(MODULE)-client ... [OK]\n"

server:
	@printf "Building $(MODULE)-server ..."
	@go install -i $(SRC_PATH)/$(MODULE)-server
	@printf "\rBuilding $(MODULE)-server ... [OK]\n"

clean:
	@go clean -cache
	@rm -rf *.tar.gz
	@rm -rf bin/*
	@rm -rf pkg/*
