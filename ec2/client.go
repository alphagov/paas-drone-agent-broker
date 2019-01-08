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
	//CreateBucket(name string) error
	//DeleteBucket(name string) error
	//AddUserToBucket(username, bucketName string) (BucketCredentials, error)
	//RemoveUserFromBucket(username, bucketName string) error
}

type BucketCredentials struct {
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

type EC2Client struct {
	Timeout      time.Duration
	EC2          ec2iface.EC2API
}

func NewEC2Client(bucketPrefix, region string) *EC2Client {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client := ec2.New(sess)
	return &EC2Client{
		Timeout:      30 * time.Second,
		EC2:          ec2Client,
	}
}