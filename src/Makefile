TEST?=$$(go list ./... | grep -v 'vendor')
###### chang variables below according to your own modules ###
NAMESPACE=kusionstack
NAME=objectstorage
KUSION_HOME=~/.kcl/kpm
VERSION=0.2.0
BINARY=../bin/kusion-module-${NAME}_${VERSION}

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
GOARCH_LOCAL := amd64
else
GOARCH_LOCAL := $(LOCAL_ARCH)
endif
export GOOS_LOCAL := $(shell uname|tr 'A-Z' 'a-z')
export OS_ARCH ?= $(GOARCH_LOCAL)

default: install

build-darwin:
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY} .

install: build-darwin
	mkdir -p ${KUSION_HOME}/modules/${NAMESPACE}/${NAME}/${VERSION}/${GOOS_LOCAL}/${OS_ARCH}
	cp ${BINARY} ${KUSION_HOME}/modules/${NAMESPACE}/${NAME}/${VERSION}/${GOOS_LOCAL}/${OS_ARCH}

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY} ./

install-linux: build-linux
	mkdir -p ${KUSION_HOME}/objectstorage_0.2.0
	cp ${BINARY} ${KUSION_HOME}/objectstorage_0.2.0

release: 
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY}_darwin_arm64 ./${NAME}
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY}_darwin_amd64 ./${NAME}
	GOOS=linux GOARCH=arm64 go build -o ${BINARY}_linux_arm64 ./${NAME}
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}_linux_amd64 ./${NAME}
	GOOS=windows GOARCH=amd64 go build -o ${BINARY}_windows_amd64 ./${NAME}
	GOOS=windows GOARCH=386 go build -o ${BINARY}_windows_386 ./${NAME}

test:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 5m
