
################################################################################
# Variables                                                                    #
################################################################################

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

OUT_DIR := ./dist

dist/$(GOARCH)/consul2slurm:
	CGO_ENABLED=$(CGO) GOARCH=$(GOARCH) go build -o $(OUT_DIR)/$(GOARCH)/consul2slurm ./cmd/consul2slurm
	mkdir -p ./package/$(GOARCH)/usr/local/bin
	ln -f $(OUT_DIR)/$(GOARCH)/consul2slurm ./package/$(GOARCH)/usr/local/bin