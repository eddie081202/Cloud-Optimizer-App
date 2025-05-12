package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"google.golang.org/grpc"
	
	"github.com/eddie081202/Cloud-Optimizer-App/backend/internal/aws"
	pb "github.com/eddie081202/Cloud-Optimizer-App/backend/proto"
)

type server struct {
	pb.UnimplementedOptimizerServiceServer
	ec2Client *ec2.Client
}

func newServer() (*server, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	return &server{
		ec2Client: ec2Client,
	}, nil
}

func (s *server) ListResources(ctx context.Context, req *pb.ListResourcesRequest) (*pb.ListResourcesResponse, error) {
	var resources []*pb.Resource
	var err error

	switch req.CloudProvider {
	case "aws":
		if req.ResourceType == "" || req.ResourceType == "ec2" {
			resources, err = aws.ListEC2Instances(ctx, s.ec2Client)
			if err != nil {
				return nil, err
			}
		}
	// TODO: Add other cloud providers
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", req.CloudProvider)
	}

	return &pb.ListResourcesResponse{
		Resources: resources,
	}, nil
}

func (s *server) GetResourceMetrics(ctx context.Context, req *GetResourceMetricsRequest) (*GetResourceMetricsResponse, error) {
	// TODO: Implement metrics retrieval logic
	return &GetResourceMetricsResponse{}, nil
}

func (s *server) GetRecommendations(ctx context.Context, req *GetRecommendationsRequest) (*GetRecommendationsResponse, error) {
	// TODO: Implement recommendations logic
	return &GetRecommendationsResponse{}, nil
}

func (s *server) ApplyOptimization(ctx context.Context, req *ApplyOptimizationRequest) (*ApplyOptimizationResponse, error) {
	// TODO: Implement optimization application logic
	return &ApplyOptimizationResponse{}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv, err := newServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOptimizerServiceServer(s, srv)
	
	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
