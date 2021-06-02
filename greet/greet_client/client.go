package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to the server %s", err.Error())
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection to the server %s", err.Error())
		}
	}()

	c := greetpb.NewGreetServiceClient(conn)

	/*
		doUnary(c)
		doServerStreaming(c)
		doClientStreaming(c)
	*/
	doBiDiStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do unary rpc")
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Daniel",
		LastName:  "Gómez",
	}}

	greetRes, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("error while calling greet from the client %s", err.Error())
	}
	log.Println("Response from greet,", greetRes.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a server streaming RPC")

	// Maybe not the best way to handle grpc streaming
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*11)
	defer cancel()

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daniel",
			LastName:  "Gómez",
		},
	}
	resStream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("error while calling greet many times RPC %s", err.Error())
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalln("error while reading stream", err.Error())
		}
		log.Println("Response from GreetManyTimes", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a client streaming RPC")

	reqs := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Daniel",
				LastName:  "Gomez",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "AAA",
				LastName:  "AA",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "BBB",
				LastName:  "BB",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "CCC",
				LastName:  "CC",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "DDD",
				LastName:  "DD",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3200*time.Millisecond)
	defer cancel()

	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalln("error while calling LongGreet", err.Error())
	}

	for _, req := range reqs {
		fmt.Println("Sending request", req.GetGreeting())
		if err := stream.Send(req); err != nil {
			log.Fatalln("error while sending stream from client", err.Error())
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error while receiving response from the server LongGreet", err.Error())
	}
	log.Println("Long greet response", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a client streaming RPC")

	reqs := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Daniel",
				LastName:  "Gomez",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Maria",
				LastName:  "Cervantes",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cologne",
				LastName:  "azure",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Cliat",
				LastName:  "romer",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Dide",
				LastName:  "Sorana",
			},
		},
	}

	// create stream by invoking the function
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalln("error while creating the stream", err.Error())
		return
	}

	// create a channel
	waitc := make(chan struct{})

	// send messages to the server (go routine)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending message", req)
			if err := stream.Send(req); err != nil {
				log.Fatalln("error while sending messages to the server")
			}
			time.Sleep(2*time.Millisecond)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalln("error closing the stream")
		}
	}()

	// receive messages from the server (go routine)

	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				//log.Fatalln("error receiving messages from server stream", err.Error())
				break
			}

			fmt.Println("Result:", res.GetResult())
		}
		close(waitc)

	}()

	// block until everything is done
	<-waitc
}
