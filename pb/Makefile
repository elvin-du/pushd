all: create

.PHONY: create

gogofaster:
	go get github.com/gogo/protobuf/protoc-gen-gogofaster

src=gate.proto session.proto

create: gogofaster
	protoc --proto_path=$(GOPATH)/src/github.com/gogo/protobuf/:. --gogofaster_out=plugins=grpc:. $(src)

