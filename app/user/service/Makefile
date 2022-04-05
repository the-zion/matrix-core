INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=../../../third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: generate
# generate
wire:
	go generate ./...