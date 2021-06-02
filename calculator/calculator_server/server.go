package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {

	log.Println("Calculator function invoked with", req)

	return &calculatorpb.CalculatorResponse{
		SumResult: req.GetFirstNumber() + req.GetSecondNumber(),
	}, nil
}

func (s *server) PrimeDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {

	log.Println("Calculator function invoked with", req)

	n := req.GetNumber()
	div := int64(2)

	for n > 1 {
		if n%div == 0 {

			if err := stream.Send(&calculatorpb.PrimeNumberResponse{
				Result: div,
			}); err != nil {
				return err
			}
			n /= div
			continue
		}
		div++
		fmt.Println("Divisor", div)
	}
	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	log.Println("Average function invoked with a streaming request")

	sum := int64(0)
	count := 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.AverageResponse{Average: float64(sum) / float64(count)})
		}

		if err != nil {
			log.Fatalln("error receiving stream from the client")
		}
		sum += req.GetNumber()
		count++
	}
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {

	log.Println("FindMaximum function invoked with a streaming request")

	prevMax := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalln("error while reading client stream", err.Error())
			return err
		}
		current := req.GetNumber()
		if current > prevMax {
			if err := stream.Send(&calculatorpb.MaximumResponse{Max: req.GetNumber()}); err != nil {
				log.Println("error while sending data to the client", err.Error())
				return err
			}
			prevMax = current
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("failed to listen", err.Error())
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("failed to serve", err.Error())
	}
}
