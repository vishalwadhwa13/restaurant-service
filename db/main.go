package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
)

const (
	FieldResID       = "ResId"
	FieldName        = "Name"
	FieldRating      = "Rating"
	FieldCuisines    = "Cuisines"
	FieldOpeningTime = "OpeningTime"
	FieldClosingTime = "ClosingTime"
	FieldLocation    = "Location"
	FieldCostForTwo  = "CostForTwo"
	TableRestaurant  = "Restaurant"
)

type RestaurantServer struct {
	db *sql.DB
}

func (s *RestaurantServer) DBInit(dbUser string, dbPassword string,
	dbHost string, dbPort string, dbName string) error {
	var err error
	s.db, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName))

	return err
}

func (s *RestaurantServer) DBClose() error {
	return s.db.Close()
}

func (s *RestaurantServer) DBPing() error {
	return s.db.Ping()
}

func (s *RestaurantServer) hasRestaurant(resID uint64) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?;", TableRestaurant, FieldResID)
	rows, err := s.db.Query(query, resID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func BuildStatus(code pb.Status_StatusCode, msg string) *pb.Status {
	return &pb.Status{Code: code, Message: msg}
}
