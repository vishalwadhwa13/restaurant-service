package db

import (
	"context"
	"fmt"
	"github.com/vishalwadhwa13/restaurant-service/db/utils"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"log"
)

func (s *RestaurantServer) GetRestaurant(ctx context.Context, req *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s, ST_AsText(%s), %s FROM %s WHERE %s = ?;",
		FieldResID, FieldName,
		FieldRating, FieldCuisines,
		FieldOpeningTime, FieldClosingTime,
		FieldLocation, FieldCostForTwo,
		TableRestaurant, FieldResID)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return &pb.GetRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.")}, err
	}

	rows, err := stmt.QueryContext(ctx, req.GetResId())
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		return &pb.GetRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.")}, err
	}

	var (
		resID       uint64
		name        string
		rating      float64
		cuisines    string
		costForTwo  float64
		openingTime string
		closingTime string
		location    string
	)

	if !rows.Next() {
		return &pb.GetRestaurantResponse{
			Status: BuildStatus(pb.Status_FAIL, "Restaurant does not exist."),
		}, nil
	}

	err = rows.Scan(&resID, &name,
		&rating, &cuisines,
		&openingTime, &closingTime,
		&location, &costForTwo)
	if err != nil {
		log.Fatalln(err)
		return &pb.GetRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.")}, err
	}

	locn, err := utils.ParseLocation(location)
	if err != nil {
		log.Fatalln(err)
		return &pb.GetRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.")}, err
	}

	r := &pb.Restaurant{
		ResId: resID, Name: name,
		Rating: rating, Cuisines: utils.ParseCuisines(cuisines),
		OpeningTime: openingTime, ClosingTime: closingTime,
		Location: locn, CostForTwo: costForTwo,
	}

	// Iteration error
	// ignore here since the only result has been obtained
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
		// return err
	}

	log.Println("GetRestaurant executed. ", r)
	return &pb.GetRestaurantResponse{
		Status:     BuildStatus(pb.Status_SUCCESS, ""),
		Restaurant: r,
	}, nil
}
