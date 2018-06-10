
.PHONY: build proto

build:
	go build

proto:
	protoc -Iproto ./proto/sprints.proto --go_out=plugins=grpc:proto
	protoc -Iproto ./proto/vis.proto --go_out=plugins=grpc:proto
