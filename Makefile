compile_protos:
	protoc --go_out=plugins=grpc:. -I. ./restaurant-service/*.proto

build:
	mkdir -p out
	go build -o out/restaurant-service

build_docker:
	docker build -t restaurant-service .