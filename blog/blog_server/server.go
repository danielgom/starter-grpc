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

	// Shutdown gracefully
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs
		log.Println("Performing shutdown...")
		s.Stop()
		log.Println("Closing MongoDB connection...")
		database.Close()
		log.Println("Server closed!")
	}()

	log.Println("Server successfully started on", time.Now().Sub(st))
	if err := s.Serve(lis); err != nil {
		log.Fatalln("error listening to the server, ", err.Error())
	}
}
