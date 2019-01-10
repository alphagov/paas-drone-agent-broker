package ec2

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

//go:generate counterfeiter -o fakes/fake_ec2_client.go . Client
type Client interface {
	RunEC2(input ec2.RunInstancesInput) (*ec2.Reservation, error)
	TerminateEC2(input ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)
	TagEC2(instanceID *string, tags []*ec2.Tag) (output *ec2.CreateTagsOutput, err error)
	IdentifyEC2(instanceID string) (output []*ec2.Reservation, err error)
	GetAmi(owner string, name string) (ami *ec2.Image, err error)
}

type EC2Client struct {
	Timeout time.Duration
	EC2     ec2iface.EC2API
}

//function to generate new EC2 instance
func NewEC2Client(region string) *EC2Client {
	config := aws.Config{Region: aws.String(region)}
	sess := session.Must(session.NewSession(&config))
	ec2Client := ec2.New(sess)
	return &EC2Client{
		Timeout: 30 * time.Second,
		EC2:     ec2Client,
	}
}

//function to return error if EC2 instance creation generates one
func (s *EC2Client) RunEC2(input ec2.RunInstancesInput) (*ec2.Reservation, error) {
	result, err := s.EC2.RunInstances(&input)
	return result, err
}

func (s *EC2Client) TerminateEC2(input ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	result, err := s.EC2.TerminateInstances(&input)
	return result, err
}

func (s *EC2Client) TagEC2(instanceID *string, tags []*ec2.Tag) (output *ec2.CreateTagsOutput, err error) {
	createTagsInput := ec2.CreateTagsInput{
		Resources: []*string{instanceID},
		Tags:      tags,
	}
	result, err := s.EC2.CreateTags(&createTagsInput)
	return result, err
}

func (s *EC2Client) IdentifyEC2(serviceRef string) (output []*ec2.Reservation, err error) {
	describeInstanceInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:service_instance_ref"),
				Values: []*string{aws.String(serviceRef)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	describeInstancesOutput, err := s.EC2.DescribeInstances(&describeInstanceInput)
	return describeInstancesOutput.Reservations, err
}

func (s *EC2Client) GetAmi(owner string, name string) (ami *ec2.Image, err error) {
	imageFilter := ec2.Filter{
		Name:   aws.String("name"),
		Values: aws.StringSlice([]string{name}),
	}
	describeImageInput := ec2.DescribeImagesInput{
		Owners:  aws.StringSlice([]string{owner}),
		Filters: []*ec2.Filter{&imageFilter},
	}

	image, err := s.EC2.DescribeImages(&describeImageInput)

	return image.Images[0], err
}
