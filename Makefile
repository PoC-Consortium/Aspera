.DEFAULT_GOAL := start

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	mkdir -p $(ROOT_DIR)/internal/p2p
	protoc -I $(ROOT_DIR) --go_out=plugins=grpc:$(ROOT_DIR)/internal/ api/protobuf-spec/p2p.proto
	go build -ldflags="-s -w" main.go
	go build -ldflags="-s -w" publish.go
	# mv $(ROOT_DIR)/internal/order/api/protobuf-spec/order.pb.go $(ROOT_DIR)/internal/order/
	# rm -r $(ROOT_DIR)/internal/order/api
start:
	go run main.go -path=$(ROOT_DIR)/var
clibs:
	cd c; \
	$(CC) $(CFLAGS) -c -o shabal64.o shabal64.s; \
	$(CC) $(CFLAGS) -c -o mshabal_sse4.o mshabal_sse4.c; \
	$(CC) $(CFLAGS) -mavx2 -c -o mshabal256_avx2.o mshabal256_avx2.c; \
	$(CC) $(CFLAGS) -shared -o libburstmath.a burstmath.c shabal64.o mshabal_sse4.o mshabal256_avx2.o -lpthread -std=gnu99;
gen-parse-transaction-tests:
	go run scripts/gen_parse_transaction_tests.go
