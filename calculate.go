package calculate

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	AwsProvider = "AWS"
)

type DescribeInstancesInput struct {
	AwsInput *ec2.DescribeInstancesInput
}

type DescribeInstancesOutput struct {
	AwsOutput *ec2.DescribeInstancesOutput
}

type Compute interface {
	Provider() string
	DescribeInstances(*DescribeInstancesInput) (*DescribeInstancesOutput, error)
	CreateImage(*CreateImageInput) (*CreateImageOutput, error)
	DescribeImages(*DescribeImagesOutput) (*DescribeImagesOutput, error)
}

var (
	NoValidProviderErr = errors.New("calculate: No valid provider was given to NewCompute()")
)

func NewCompute(provider, region string) (Compute, error) {
	switch provider {
	case AwsProvider:
		return Compute(NewEC2(region)), nil
	}
	return nil, NoValidProviderErr
}

func NewEC2WithConfig(c *aws.Config) *EC2 {
	return &EC2{ec2.New(session.New(c))}
}

type EC2 struct {
	service *ec2.EC2
}

func NewEC2(region string) *EC2 {
	return NewEC2WithConfig(&aws.Config{Region: aws.String(region)})
}

func (e *EC2) DescribeInstances(input *DescribeInstancesInput) (*DescribeInstancesOutput, error) {
	if input == nil {
		input = &DescribeInstancesInput{}
	}
	output, err := e.DescribeEC2Instances(input.AwsInput)
	return &DescribeInstancesOutput{AwsOutput: output}, err
}

func (e *EC2) DescribeEC2Instances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	instances, err := e.service.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	for instances.NextToken != nil {
		input.NextToken = instances.NextToken
		i, err := e.service.DescribeInstances(input)
		if err != nil {
			return nil, err
		}
		instances.Reservations = append(instances.Reservations, i.Reservations...)
		instances.NextToken = i.NextToken
	}
	return instances, nil
}

func (e *EC2) Provider() string {
	return AwsProvider
}

func (e *EC2) SetRegion(s string) {
	e.service.Config.Region = aws.String(s)
}

type CreateImageInput struct {
	AwsInput *ec2.CreateImageInput
}

type CreateImageOutput struct {
	AwsOutput *ec2.CreateImageOutput
}

func (e *EC2) CreateImage(input *CreateImageInput) (*CreateImageOutput, error) {
	snapshot, err := e.service.CreateImage(input.AwsInput)
	return &CreateImageOutput{AwsOutput: snapshot}, err
}

type DescribeImagesInput struct {
	AwsInput *ec2.DescribeImagesInput
}

type DescribeImagesOutput struct {
	AwsOutput *ec2.DescribeImagesOutput
}

func (e *EC2) DescribeImages(input *DescribeImagesInput) (*DescribeImagesOutput, error) {
	output, err := e.service.DescribeImages(input.AwsInput)
	return &DescribeImagesOutput{AwsOutput:output}, err
}

type CreateSnapshotInput struct {}

type CreateSnapshotOutput struct {}

func (e *EC2) CreateSnapshot(input *CreateSnapshotInput) (*CreateSnapshotOutput, error) {
	return nil, nil
}