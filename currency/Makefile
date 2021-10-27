.PHONY: protos

protos:
	protoc -I=./protos --go_opt=paths=source_relative --go_out=plugins=grpc:./protos/currency ./protos/currency.proto
