package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}


func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Println("Greet function was invoked with", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName

	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Println("GreetManyTimes function was invoked with", req)

	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++{
		result := "Hello " + firstName + " " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 *time.Millisecond)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Println("LongGreet function was invoked with a streaming request")

	result := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalln("error while reading client stream", err.Error())
			return err
		}

		result += " hello " + req.GetGreeting().GetFirstName() + " " + req.GetGreeting().GetLastName()
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with a streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalln("error while reading client stream", err.Error())
			return err
		}

		result := " Hello " + req.GetGreeting().GetFirstName() + " " + req.GetGreeting().GetLastName()

		res := &greetpb.GreetEveryoneResponse{Result: result}
		if err := stream.Send(res); err != nil {
			log.Fatalln("error while sending stream to the client")
		}
	}
}

func main() {

	// Default port for grpc 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %s", err.Error())
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %s", err.Error())
	}
}
