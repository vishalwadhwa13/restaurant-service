package db

import (
	"fmt"
	"github.com/vishalwadhwa13/restaurant-service/db/utils"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"log"
)

func (s *RestaurantServer) GetAllRestaurant(req *pb.GetAllRestaurantRequest, stream pb.RestaurantService_GetAllRestaurantServer) error {
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s, ST_AsText(%s), %s FROM %s;",
		FieldResID, FieldName,
		FieldRating, FieldCuisines,
		FieldOpeningTime, FieldClosingTime,
		FieldLocation, FieldCostForTwo,
		TableRestaurant)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
			pb.Status_FAIL,
			"Unable to get restaurant.",
		)})
		log.Fatalln(err)
		return err
	}

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
			pb.Status_FAIL,
			"Unable to get restaurant.",
		)})
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
		location    string
	)
	for rows.Next() {
		err := rows.Scan(&resID, &name,
			&rating, &cuisines,
			&openingTime, &closingTime,
			&location, &costForTwo)
		if err != nil {
			stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.",
			)})
			log.Fatalln(err)
			return err
		}

		locn, err := utils.ParseLocation(location)
		if err != nil {
			stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.",
			)})
			log.Fatalln(err)
			return err
		}

		r := &pb.Restaurant{
			ResId: resID, Name: name,
			Rating: rating, Cuisines: utils.ParseCuisines(cuisines),
			OpeningTime: openingTime, ClosingTime: closingTime,
			Location: locn, CostForTwo: costForTwo,
		}

		validResp := &pb.GetRestaurantResponse{
			Status: BuildStatus(pb.Status_SUCCESS, ""), Restaurant: r}
		if err := stream.Send(validResp); err != nil {
			stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
				pb.Status_FAIL,
				"Unable to get restaurant.",
			)})
			log.Fatalln(err)
			return err
		}
		log.Println("Sent rest. ", r)
	}

	if err := rows.Err(); err != nil {
		stream.Send(&pb.GetRestaurantResponse{Status: BuildStatus(
			pb.Status_FAIL,
			"Unable to get restaurant.",
		)})
		log.Fatalln(err)
		return err
	}
	return nil
}
