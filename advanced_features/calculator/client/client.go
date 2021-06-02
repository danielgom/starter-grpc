package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/advanced_features/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {

	// CA file
	certFile := "ssl/ca.crt"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")

	if err != nil {
		log.Fatalln("Failed to load CA certificate")
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalln("error connecting to the server,", err.Error())
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln("error closing the server,", err.Error())
		}
	}()

	c := calculatorpb.NewAdvancedServiceClient(conn)

	doErrorUnary(c)
	//doDeadlineUnary(c)
}

func doErrorUnary(c calculatorpb.AdvancedServiceClient) {

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: 33})

	if err != nil {
		respErr, ok := status.FromError(err)

		if !ok {
			log.Fatalln("error calling number root function", err.Error())
		}
		// Actual error from gRPC (user error)
		fmt.Print("error message from server: ", respErr.Message())

		if respErr.Code() == codes.InvalidArgument {
			fmt.Println("please send only positive numbers")
		}
		return
	}

	fmt.Println("Result square root:", res.GetNumberRoot())
}

func doDeadlineUnary(c calculatorpb.AdvancedServiceClient) {

	req := &calculatorpb.SaluteWithDeadlineRequest{Salute: &calculatorpb.Salute{
		FirstName: "Daniel",
		LastName:  "G.",
	}}


	// Max 5 second deadline
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.SaluteWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)

		if !ok {
			log.Fatalln("error while calling the server function", err.Error())
			return
		}

		if statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout hit, deadline exceeded")
		} else{
			fmt.Println("unexpected error,", statusErr)
		}
		return
	}

	fmt.Println("Result:", res.GetResponse())
}
