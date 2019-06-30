compile_protos:
	protoc --go_out=plugins=grpc:. -I. ./protos/*.proto

build:
	mkdir -p out
	go build -o out/restaurant_service

docker_build:
	docker build -t restaurant-service .