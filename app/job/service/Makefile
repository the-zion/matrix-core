INTERNAL_EASYJSON_FILES=$(shell find cmd -name main.go)

.PHONY: easyjson
# generate internal proto
easyjson:
	easyjson -lower_camel_case $(INTERNAL_EASYJSON_FILES)