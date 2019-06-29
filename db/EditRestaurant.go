package db

import (
	"context"
	"fmt"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"log"
	"strings"
)

func (s *RestaurantServer) EditRestaurant(ctx context.Context, req *pb.EditRestaurantRequest) (*pb.EditRestaurantResponse, error) {
	hasRestaurant, err := s.hasRestaurant(req.GetRestaurant().GetResId())
	if err != nil {
		log.Fatalln(err)
		return &pb.EditRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to edit restaurant.")}, err
	}

	if !hasRestaurant {
		log.Println("Restaurant does not exist.")
		return &pb.EditRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"No such restaurant exists.")}, nil
	}

	query := fmt.Sprintf("REPLACE INTO %s VALUES(?, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?);", TableRestaurant)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return &pb.EditRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to edit restaurant.")}, err
	}

	r := req.GetRestaurant()
	res, err := stmt.ExecContext(
		ctx,
		r.GetResId(),
		r.GetName(),
		r.GetRating(),
		strings.Join(r.GetCuisines()[:], ","),
		r.GetOpeningTime(),
		r.GetClosingTime(),
		fmt.Sprintf("POINT(%f %f)", r.GetLocation().GetLat(), r.GetLocation().GetLong()),
		r.GetCostForTwo(),
	)
	if err != nil {
		log.Fatalln(err)
		return &pb.EditRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to edit restaurant.")}, err
	}

	log.Println("EditRestaurant executed.", res)

	return &pb.EditRestaurantResponse{
		Status: BuildStatus(
			pb.Status_SUCCESS,
			"Restaurant edited successfully.",
		),
	}, nil
}
