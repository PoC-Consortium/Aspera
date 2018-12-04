.DEFAULT_GOAL := start

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: clibs build

build:
	mkdir -p $(ROOT_DIR)/pkg/api/p2p
	mkdir -p $(ROOT_DIR)/pkg/account/pb
	protoc -I $(ROOT_DIR)/api/protobuf-spec/p2p --go_out=plugins=grpc:$(ROOT_DIR)/pkg/api/p2p transaction.proto p2p.proto
	protoc -I $(ROOT_DIR)/api/protobuf-spec/account --go_out=plugins=grpc:$(ROOT_DIR)/pkg/account/pb account.proto
	qtc $(ROOT_DIR)/pkg/api/p2p/compat/template/block.qtpl
	go build -ldflags="-s -w" main.go
start:
	go run main.go -path=$(ROOT_DIR)/var
clibs:
	cd c; \
	$(CC) $(CFLAGS) -c -o shabal64.o shabal64.s; \
	$(CC) $(CFLAGS) -c -o mshabal_sse4.o mshabal_sse4.c; \
	$(CC) $(CFLAGS) -mavx2 -c -o mshabal256_avx2.o mshabal256_avx2.c; \
	$(CC) $(CFLAGS) -shared -o libburstmath.a burstmath.c shabal64.o mshabal_sse4.o mshabal256_avx2.o -lpthread -std=gnu99;
