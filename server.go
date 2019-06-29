package main

import (
	"github.com/joho/godotenv"
	"github.com/vishalwadhwa13/restaurant-service/db"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var dbUser string
var dbPassword string
var dbName string
var dbHost string
var dbPort string

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	dbUser = os.Getenv("DB_USER")
	dbHostTmp := os.Getenv("ENV")
	if dbHostTmp == "docker" {
		dbHost = "host.docker.internal"
	} else {
		dbHost = ""
	}
	dbPassword = os.Getenv("DB_PASSWORD")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")

	return nil
}

func main() {
	var err error
	restServer := &db.RestaurantServer{}

	err = initEnv()
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = restServer.DBInit(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer restServer.DBClose()

	// Open doesn't open a connection. Validate DSN data:
	err = restServer.DBPing()
	if err != nil {
		log.Fatalln(err)
		return
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}

	log.Println("Starting server.")
	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(grpcServer, restServer)
	grpcServer.Serve(lis)
}
