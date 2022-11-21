INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
INTERNAL_EASYJSON_DO_FILES=$(shell find internal -name do.go)
INTERNAL_EASYJSON_MODEL_FILES=$(shell find internal -name model.go)
.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=../../../third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: easyjson
# generate internal proto
easyjson:
	easyjson -lower_camel_case $(INTERNAL_EASYJSON_DO_FILES)
	easyjson -lower_camel_case $(INTERNAL_EASYJSON_MODEL_FILES)