package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/joho/godotenv"
	pb "github.com/vishalwadhwa13/restaurant-service/restaurant-service"
	"github.com/vishalwadhwa13/restaurant-service/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

var db *sql.DB

var dbUser string
var dbPassword string
var dbName string
var dbHost string
var dbPort string

const (
	FieldResID       = "ResId"
	FieldName        = "Name"
	FieldRating      = "Rating"
	FieldCuisines    = "Cuisines"
	FieldOpeningTime = "OpeningTime"
	FieldClosingTime = "ClosingTime"
	FieldCoordinates = "Coordinates"
	FieldCostForTwo  = "CostForTwo"
	TableRestaurant  = "Restaurant"
)

type RestaurantServer struct {
}

func (s *RestaurantServer) AddRestaurant(ctx context.Context, r *pb.Restaurant) (*empty.Empty, error) {
	query := fmt.Sprintf("INSERT INTO %s VALUES(NULL, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?);", TableRestaurant)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	res, err := stmt.ExecContext(
		ctx,
		r.GetName(),
		r.GetRating(),
		strings.Join(r.GetCuisines()[:], ","),
		r.GetOpeningTime(),
		r.GetClosingTime(),
		fmt.Sprintf("POINT(%f %f)", r.GetCoordinates().GetLat(), r.GetCoordinates().GetLong()),
		r.GetCostForTwo(),
	)

	log.Println("AddRestaurant executed.", res)
	return &empty.Empty{}, nil
}

func (s *RestaurantServer) EditRestaurant(ctx context.Context, r *pb.Restaurant) (*empty.Empty, error) {
	query := fmt.Sprintf("REPLACE INTO %s VALUES(?, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?);", TableRestaurant)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	res, err := stmt.ExecContext(
		ctx,
		r.GetResId(),
		r.GetName(),
		r.GetRating(),
		strings.Join(r.GetCuisines()[:], ","),
		r.GetOpeningTime(),
		r.GetClosingTime(),
		fmt.Sprintf("POINT(%f %f)", r.GetCoordinates().GetLat(), r.GetCoordinates().GetLong()),
		r.GetCostForTwo(),
	)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	log.Println("EditRestaurant executed.", res)
	return &empty.Empty{}, nil
}

func (s *RestaurantServer) DeleteRestaurant(ctx context.Context, r *pb.Restaurant) (*empty.Empty, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", TableRestaurant, FieldResID)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, r.GetResId())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	log.Println("DeleteRestaurant executed.", res)
	return &empty.Empty{}, nil
}

func (s *RestaurantServer) GetAllRestaurant(e *empty.Empty, stream pb.RestaurantService_GetAllRestaurantServer) error {
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s, ST_AsText(%s), %s FROM %s;",
		FieldResID, FieldName,
		FieldRating, FieldCuisines,
		FieldOpeningTime, FieldClosingTime,
		FieldCoordinates, FieldCostForTwo,
		TableRestaurant)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	var (
		resID       uint64
		name        string
		rating      float64
		cuisines    string
		costForTwo  float64
		openingTime string
		closingTime string
		coords      string
	)
	for rows.Next() {
		err := rows.Scan(&resID, &name,
			&rating, &cuisines,
			&openingTime, &closingTime,
			&coords, &costForTwo)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		crds, err := utils.ParseCoordinates(coords)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		r := &pb.Restaurant{
			ResId: resID, Name: name,
			Rating: rating, Cuisines: utils.ParseCuisines(cuisines),
			OpeningTime: openingTime, ClosingTime: closingTime,
			Coordinates: crds, CostForTwo: costForTwo,
		}

		if err := stream.Send(r); err != nil {
			log.Fatalln(err)
			return err
		}
		log.Println("Sent rest. ", r)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

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

	initEnv()
	if err != nil {
		log.Fatalln(err)
		return
	}

	db, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName))

	if err != nil {
		db = nil
		log.Fatalln(err)
		return
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
		return
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	log.Println("Starting server.")
	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(grpcServer, &RestaurantServer{})
	grpcServer.Serve(lis)

}
