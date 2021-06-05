package main

import (
	"github.com/danielgom/starter-grpc/blog/blogpb"
	"github.com/danielgom/starter-grpc/blog/database"
	"github.com/danielgom/starter-grpc/blog/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Starting Blog server...")

	st := time.Now()

	database.Init()

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln("failed to listen,", err.Error())
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"

	cred, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

	if sslErr != nil {
		log.Fatalln("Failed loading the certificates", sslErr.Error())
		return
	}

	s := grpc.NewServer(grpc.Creds(cred))

	blogpb.RegisterBlogServiceServer(s, &services.BlogService{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	go func() {
		log.Println("Server successfully started on", time.Now().Sub(st))
		if err := s.Serve(lis); err != nil {
			log.Fatalln("error listening to the server, ", err.Error())
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Println("Stopping the server...")
	s.Stop()
	log.Println("Stopping the listener...")
	log.Println("Closing MongoDB connection...")
	log.Println("Server closed")
}
