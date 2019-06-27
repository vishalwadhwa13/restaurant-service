package utils

import (
	pb "github.com/vishalwadhwa13/restaurant-service/restaurant-service"
	"regexp"
	"strconv"
	"strings"
)

var pointParseRE = regexp.MustCompile(`POINT\((\S+) (\S+)\)`)

func ParseCuisines(c string) []string {
	return strings.Split(c, ",")
}

func ParseCoordinates(c string) (*pb.Restaurant_Point, error) {
	res := pointParseRE.FindStringSubmatch(c)

	lat, err := strconv.ParseFloat(res[1], 64)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseFloat(res[2], 64)
	if err != nil {
		return nil, err
	}

	return &pb.Restaurant_Point{Lat: lat, Long: long}, nil
}
