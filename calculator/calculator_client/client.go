package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("could not connect to the server", err.Error())
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalln("error closing connection", err.Error())
		}
	}()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	/*
		doSum(c)
		doPrimeDecomposition(c)
		doComputeAverage(c)
	*/
	doFindMaximum(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	req := &calculatorpb.CalculatorRequest{
		FirstNumber:  10,
		SecondNumber: 3,
	}

	res, err := c.Sum(ctx, req)
	if err != nil {
		log.Fatalln("error retrieving sum result", err.Error())
	}

	fmt.Println("sum result:", res.SumResult)
}

func doPrimeDecomposition(c calculatorpb.CalculatorServiceClient) {

	log.Println("Starting to receive RPC prime decomposition")

	req := &calculatorpb.PrimeNumberRequest{
		Number: 12390392840,
	}

	streamRes, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalln("error retrieving prime decomposition result", err.Error())
	}

	for {
		msg, err := streamRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("error while reading stream", err.Error())
		}
		log.Println("Response", msg.GetResult())
	}
}

func doComputeAverage(c calculatorpb.CalculatorServiceClient) {

	log.Println("Starting to send stream to the server")

	reqs := []*calculatorpb.AverageRequest{
		{Number: 1},
		{Number: 2},
		{Number: 3},
		{Number: 4},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalln("error opening stream")
	}

	for _, req := range reqs {
		log.Println("Sending number", req.GetNumber())
		if err := stream.Send(req); err != nil {
			log.Fatalln("error sending stream to the server")
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error getting response from server")
	}

	log.Println("Average:", res.GetAverage())
}

func doFindMaximum(c calculatorpb.CalculatorServiceClient) {

	log.Println("Starting to send and receive data from the server")

	reqs := []*calculatorpb.MaximumRequest{
		{Number: 1},
		{Number: 5},
		{Number: 3},
		{Number: 6},
		{Number: 2},
		{Number: 20},
	}

	// Create stream from the FindMaximum function

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := c.FindMaximum(ctx)


	// Sending
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending:", req.GetNumber())
			if err = stream.Send(req); err != nil {
				log.Fatalln("error sending stream to the server")
			}
			time.Sleep(2*time.Millisecond)
		}

		if err = stream.CloseSend(); err != nil {
			log.Fatalln("error closing the stream")
		}
	}()

	// Receiving
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				cancel()
				break
			}
			if err != nil {
				log.Fatalln("error receiving data from server stream", err.Error())
			}

			fmt.Println("Received new max:", res.GetMax())
		}
	}()

	select {
	case <-time.After(15*time.Millisecond):
		fmt.Println("closing after 10 milliseconds")
		cancel()
	case <-ctx.Done():
		fmt.Println("finished streaming closed successfully")
	}
}
