package db

import (
	"context"
	"fmt"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"log"
	"strings"
)

func (s *RestaurantServer) AddRestaurant(ctx context.Context, req *pb.AddRestaurantRequest) (*pb.AddRestaurantResponse, error) {
	query := fmt.Sprintf("INSERT INTO %s VALUES(NULL, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?);", TableRestaurant)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return &pb.AddRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to add restaurant.")}, err
	}

	r := req.GetRestaurant()
	res, err := stmt.ExecContext(
		ctx,
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
		return &pb.AddRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to add restaurant.")}, err
	}

	resID, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
		return &pb.AddRestaurantResponse{
			Status: BuildStatus(
				pb.Status_SUCCESS,
				"Restaurant added but could not obtain ResId.")}, err
	}

	log.Println("AddRestaurant executed.", res)
	return &pb.AddRestaurantResponse{
		Status: BuildStatus(
			pb.Status_SUCCESS,
			"Restaurant added successfully."),
		ResId: uint64(resID)}, nil
}
