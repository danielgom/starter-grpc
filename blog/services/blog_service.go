package services

import (
	"context"
	"fmt"
	"github.com/danielgom/starter-grpc/blog/blogpb"
	"github.com/danielgom/starter-grpc/blog/database"
	"github.com/danielgom/starter-grpc/blog/domain/blogs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type BlogService struct{}

func (b *BlogService) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := new(blogs.BlogItem)
	data.SetAuthorID(blog.GetAuthorId())
	data.SetContent(blog.GetContent())
	data.SetTitle(blog.GetTitle())

	res, err := database.Collection.InsertOne(ctx, data)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintln("Internal server error", err.Error()))
	}

	oID, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintln("Cannot convert to oID"))
	}

	return &blogpb.CreateBlogResponse{Blog: &blogpb.Blog{
		Id:       oID.Hex(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}}, nil
}

func (b *BlogService) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	ID := req.GetBlogId()
	fmt.Println(ID)

	oID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintln("Cannot parse ID", err.Error()))
	}

	data := new(blogs.BlogItem)
	filter := primitive.M{"_id": oID}

	res := database.Collection.FindOne(ctx, filter)

	if err = res.Decode(data); err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintln("Blog not found with ID", ID))
	}

	fmt.Println(data)

	return &blogpb.ReadBlogResponse{Blog: dataToBlogPb(data)}, nil
}

func (b *BlogService) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {

	blog := req.GetBlog()

	oID, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintln("Cannot parse ID", err.Error()))
	}

	data := new(blogs.BlogItem)
	filter := primitive.M{"_id": oID}

	res := database.Collection.FindOne(ctx, filter)

	if err = res.Decode(data); err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintln("Blog not found with ID", blog.GetId()))
	}

	data.SetAuthorID(blog.GetAuthorId())
	data.SetContent(blog.GetContent())
	data.SetTitle(blog.GetTitle())

	_, updateErr := database.Collection.ReplaceOne(ctx, filter, data)
	if updateErr != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintln("Cannot update object in MongoDB", updateErr.Error()))
	}

	return &blogpb.UpdateBlogResponse{Blog: dataToBlogPb(data)}, nil
}

func (b *BlogService) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {

	ID := req.GetBlogId()

	oID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintln("Cannot parse ID", err.Error()))
	}

	filter := primitive.M{"_id": oID}

	res, err := database.Collection.DeleteOne(ctx, filter)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintln("Cannot delete blog in MongoDB", err.Error()))
	}

	if res.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintln("Cannot find blog ", ID))
	}

	return &blogpb.DeleteBlogResponse{BlogId: ID}, nil
}

func (b *BlogService) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	cur, err := database.Collection.Find(context.Background(), primitive.D{})

	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintln("Internal server error:", err.Error()))
	}

	defer func() {
		if err = cur.Close(context.Background()); err != nil {
			log.Fatalln("error closing the listing", err.Error())
		}
	}()

	for cur.Next(context.Background()) {
		data := new(blogs.BlogItem)
		if err = cur.Decode(data); err != nil {
			return status.Error(codes.Internal, fmt.Sprintln("Error decoding data from MongoDB", err.Error()))
		}

		if err = stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)}); err != nil {
			return status.Error(codes.Internal, fmt.Sprintln("Error streaming data", err.Error()))
		}
	}

	if err = cur.Err(); err != nil {
		return status.Error(codes.Internal, fmt.Sprintln("Internal server error", err.Error()))
	}

	return nil
}

func dataToBlogPb(data *blogs.BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.Id().Hex(),
		AuthorId: data.GetAuthorID(),
		Title:    data.GetTitle(),
		Content:  data.GetContent(),
	}
}
