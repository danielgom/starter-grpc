package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/advanced_features/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"math"
	"net"
	"time"
)

type server struct{}

func (s server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	n := req.GetNumber()
	if n < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintln("received a negative number:", n))
	}

	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(n))}, nil
}

func (s server) SaluteWithDeadline(ctx context.Context, req *calculatorpb.SaluteWithDeadlineRequest) (*calculatorpb.SaluteWithDeadlineResponse, error) {

	// Testing if the client has canceled the request
	for x := 0; x < 3; x++ {

		if ctx.Err() == context.Canceled {
			return nil, status.Error(codes.Canceled, "client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetSalute().GetFirstName()

	return &calculatorpb.SaluteWithDeadlineResponse{Response: "Hello" + firstName}, nil

}

func main() {

	tlsEnabled := false

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln("failed to listen,", err.Error())
	}

	var s *grpc.Server

	if tlsEnabled {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

		if sslErr != nil {
			log.Fatalln("Failed loading the certificates", sslErr.Error())
			return
		}

		s = grpc.NewServer(grpc.Creds(creds))

	} else {
		s = grpc.NewServer()
	}

	calculatorpb.RegisterAdvancedServiceServer(s, &server{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalln("error listening to the server,", err.Error())
	}
}
