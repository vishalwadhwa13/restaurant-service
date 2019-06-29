package db

import (
	"context"
	"fmt"
	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"log"
)

func (s *RestaurantServer) DeleteRestaurant(ctx context.Context, r *pb.DeleteRestaurantRequest) (*pb.DeleteRestaurantResponse, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", TableRestaurant, FieldResID)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return &pb.DeleteRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL, "Unable to delete restaurant.")}, err
	}

	res, err := stmt.ExecContext(ctx, r.GetResId())
	if err != nil {
		log.Fatalln(err)
		return &pb.DeleteRestaurantResponse{
			Status: BuildStatus(
				pb.Status_FAIL, "Unable to delete restaurant.")}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &pb.DeleteRestaurantResponse{
			Status: BuildStatus(
				pb.Status_SUCCESS,
				"")}, err
	}

	log.Println("DeleteRestaurant executed.", res)

	if rowsAffected == 0 {
		return &pb.DeleteRestaurantResponse{Status: BuildStatus(
			pb.Status_FAIL,
			"Restaurant does not exist.",
		)}, nil
	}

	return &pb.DeleteRestaurantResponse{Status: BuildStatus(
		pb.Status_SUCCESS,
		"Restaurant deleted successfully.",
	)}, nil
}
