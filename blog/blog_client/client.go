package main

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
)

func main() {

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

	c := blogpb.NewBlogServiceClient(conn)

	CreateBlog(c)
	ReadBlog(c)
	UpdateBlog(c)
	ListBlog(c)
}

func CreateBlog(c blogpb.BlogServiceClient) {

	blog := &blogpb.CreateBlogRequest{Blog: &blogpb.Blog{
		AuthorId: "nunu's",
		Title:    "Cave Icy melting",
		Content:  "Throw bananas?",
	}}

	res, err := c.CreateBlog(context.Background(), blog)

	if err != nil {
		log.Fatalln("error sending the CreateBlog request to the server", err.Error())
	}

	log.Println("Blog has been created", res.GetBlog())
}

func ReadBlog(c blogpb.BlogServiceClient) {

	req := &blogpb.ReadBlogRequest{BlogId: "60b71fa138be3779484e0699"}

	blog, err := c.ReadBlog(context.Background(), req)

	if err != nil {
		log.Fatalln("error happened while reading", err.Error())
	}

	fmt.Println("Blogsssss", blog.GetBlog())
}

func UpdateBlog(c blogpb.BlogServiceClient) {

	blog := &blogpb.UpdateBlogRequest{Blog: &blogpb.Blog{
		Id:       "60bbde40e2c36b84551bd348",
		AuthorId: "Shivana has fire",
		Title:    "Melting Ice",
		Content:  "Throw fire blast",
	}}

	res, err := c.UpdateBlog(context.Background(), blog)

	if err != nil {
		log.Fatalln("error sending the UpdateBlog request to the server", err.Error())
	}

	log.Println("Blog has been updated", res.GetBlog())
}

func DeleteBlog(c blogpb.BlogServiceClient) {

	req := &blogpb.ReadBlogRequest{BlogId: "60b71fa138be3779484e0699"}

	blog, err := c.ReadBlog(context.Background(), req)

	if err != nil {
		log.Fatalln("error happened while reading", err.Error())
	}

	fmt.Println("Blogsssss", blog.GetBlog())
}

func ListBlog(c blogpb.BlogServiceClient) {

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalln("error calling ListBlog from the server", err.Error())
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Something happened while streaming blogs", err.Error())
		}

		fmt.Println(res.GetBlog())
	}
}
