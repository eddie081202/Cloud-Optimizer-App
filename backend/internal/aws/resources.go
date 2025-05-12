package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "github.com/eddie081202/Cloud-Optimizer-App/backend/proto"
)

// ListEC2Instances retrieves all EC2 instances and converts them to our Resource type
func ListEC2Instances(ctx context.Context, client *ec2.Client) ([]*pb.Resource, error) {
	input := &ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	var resources []*pb.Resource
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resource := &pb.Resource{
				Id:            *instance.InstanceId,
				Name:          getInstanceName(instance),
				Type:          "ec2",
				CloudProvider: "aws",
				Region:        *instance.Placement.AvailabilityZone,
				Tags:          convertTags(instance.Tags),
				CreatedAt:     timestamppb.New(*instance.LaunchTime),
			}
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// getInstanceName extracts the Name tag from instance tags
func getInstanceName(instance types.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return *instance.InstanceId
}

// convertTags converts AWS tags to our proto map format
func convertTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[*tag.Key] = *tag.Value
	}
	return result
} 