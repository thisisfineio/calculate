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
}

var (
	NoValidProviderProvidedErr = errors.New("calculate: No valid provider was given to NewCompute()")
)

func NewCompute(provider, region string) (Compute, error) {
	switch provider {
	case AwsProvider:
		return NewEC2(region), nil
	}
	return nil, NoValidProviderProvidedErr
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

func NewAwsCompute(e *EC2) Compute {
	return e
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
