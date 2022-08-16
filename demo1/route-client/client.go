package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Youngpig1998/grpc-demos/demo1/utils"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/Youngpig1998/grpc-demos/demo1/route"
	"google.golang.org/grpc"
)

func getFeature(client pb.RouteGuideClient) {
	feature, err := client.GetFeature(context.Background(), &pb.Point{
		Latitude:  310235000,
		Longitude: 121437403,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(feature)
}

func listFeatures(client pb.RouteGuideClient) {
	serverStream, err := client.ListFeatures(context.Background(), &pb.Rectangle{
		Lo: &pb.Point{Latitude: 313374060, Longitude: 121358540},
		Hi: &pb.Point{Latitude: 311034130, Longitude: 121598790},
	})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		feature, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(feature)
	}
}

func recordRoute(client pb.RouteGuideClient) {
	// dummy data
	points := []*pb.Point{
		{Latitude: 313374060, Longitude: 121358540},
		{Latitude: 311034130, Longitude: 121598790},
		{Latitude: 310235000, Longitude: 121437403},
	}

	clientStream, err := client.RecordRoute(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for _, point := range points {
		if err := clientStream.Send(point); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}
	summary, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(summary)
}

func recommend(client pb.RouteGuideClient) {
	stream, err := client.Recommend(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// this goroutine listen to the server stream
	go func() {
		for {
			feature, err2 := stream.Recv()
			if err2 == io.EOF {
				break
			}
			if err2 != nil {
				log.Fatalln(err2)
			}
			fmt.Println("Recommended: ", feature)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		request := pb.RecommendationRequest{Point: new(pb.Point)}
		var mode int32
		fmt.Print("Enter Recommendation Mode (0 for farthest, 1 for the nearest)")
		utils.ReadIntFromCommandLine(reader, &mode)
		fmt.Print("Enter Latitude: ")
		utils.ReadIntFromCommandLine(reader, &request.Point.Latitude)
		fmt.Print("Enter Longitude: ")
		utils.ReadIntFromCommandLine(reader, &request.Point.Longitude)
		request.Mode = pb.RecommendationMode(mode)

		if err := stream.Send(&request); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln("client cannot dial grpc server")
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)

	getFeature(client)

}
