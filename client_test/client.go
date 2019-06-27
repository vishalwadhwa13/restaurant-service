package main

import (
	"context"
	// "github.com/golang/protobuf/ptypes/empty"

	pb "github.com/vishalwadhwa13/restaurant-service/restaurant-service"
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
	// stream, err := client.GetAllRestaurant(context.Background(), &empty.Empty{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for {
	// 	rest, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	// 	}
	// 	log.Println(rest)
	// }

	// tmparr := [2]string{"thai", "continental"}
	// add rest
	// client.AddRestaurant(context.Background(), &pb.Restaurant{
	// 	Name:        "bbq nation",
	// 	Rating:      4.5,
	// 	Cuisines:    tmparr[:],
	// 	CostForTwo:  1200,
	// 	OpeningTime: "08:00",
	// 	ClosingTime: "22:00",
	// 	Coordinates: &pb.Restaurant_Point{Lat: 72.123, Long: 78.1321},
	// })

	// edit rest
	// client.EditRestaurant(context.Background(), &pb.Restaurant{
	// 	ResId:       3,
	// 	Name:        "bbq nation",
	// 	Rating:      4.5,
	// 	Cuisines:    tmparr[:1],
	// 	CostForTwo:  1200,
	// 	OpeningTime: "08:00",
	// 	ClosingTime: "20:00",
	// 	Coordinates: &pb.Restaurant_Point{Lat: 72.123, Long: 78.1321},
	// })

	// delete rest
	client.DeleteRestaurant(context.Background(), &pb.Restaurant{ResId: 4})
}
