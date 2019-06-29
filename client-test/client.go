package main

import (
	"context"
	// "github.com/golang/protobuf/ptypes/empty"

	pb "github.com/vishalwadhwa13/restaurant-service/protos"
	"google.golang.org/grpc"

	// "io"
	"log"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRestaurantServiceClient(conn)

	// get all
	// stream, err := client.GetAllRestaurant(context.Background(), &pb.GetAllRestaurantRequest{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for {
	// 	resp, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	// 	}
	// 	log.Println(resp)
	// }

	// tmparr := [3]string{"thai", "continental", "chinese"}
	// add rest
	// res, err := client.AddRestaurant(context.Background(), &pb.AddRestaurantRequest{Restaurant: &pb.Restaurant{
	// 	Name:        "bbq nation 2.O",
	// 	Rating:      4.5,
	// 	Cuisines:    tmparr[:],
	// 	CostForTwo:  1800,
	// 	OpeningTime: "08:00",
	// 	ClosingTime: "22:00",
	// 	Location:    &pb.Restaurant_Location{Lat: 72.123, Long: 78.1321},
	// }})
	// log.Println(res, err)

	// edit rest
	// res, err := client.EditRestaurant(context.Background(), &pb.EditRestaurantRequest{Restaurant: &pb.Restaurant{
	// 	ResId:       12,
	// 	Name:        "bbq nation",
	// 	Rating:      4.5,
	// 	Cuisines:    tmparr[1:],
	// 	CostForTwo:  1200,
	// 	OpeningTime: "08:00",
	// 	ClosingTime: "20:00",
	// 	Location:    &pb.Restaurant_Location{Lat: 72.123, Long: 78.1321},
	// }})
	// log.Println(res.GetStatus().GetCode(), err)

	// delete rest
	// res, err := client.DeleteRestaurant(context.Background(), &pb.DeleteRestaurantRequest{ResId: 7})
	// log.Println(res, err)

	// get rest
	res, err := client.GetRestaurant(context.Background(), &pb.GetRestaurantRequest{ResId: 6})
	log.Println(res, err)
}
