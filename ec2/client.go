package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

//go:generate counterfeiter -o fakes/fake_ec2_client.go . Client
type Client interface {
	RunEC2(input ec2.RunInstancesInput) (*ec2.Reservation, error)
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
