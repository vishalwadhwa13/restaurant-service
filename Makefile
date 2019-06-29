compile_protos:
	protoc --go_out=plugins=grpc:. -I. ./protos/*.proto

build:
	mkdir -p out
	go build -o out/restaurant-service

build_docker:
	docker build -t restaurant-service .