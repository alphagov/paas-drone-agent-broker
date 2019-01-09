package ec2_test

import (
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	awsEC2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/richardTowers/paas-drone-agent-broker/ec2"
	fakeClient "github.com/richardTowers/paas-drone-agent-broker/ec2/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		ec2API    *fakeClient.FakeEC2API
		ec2Client ec2.EC2Client
	)

	BeforeEach(func() {
		ec2API = &fakeClient.FakeEC2API{}
		ec2Client = ec2.EC2Client{
			Timeout: 2 * time.Second,
			EC2:     ec2API,
		}
	})

	Describe("RunEC2", func() {
		It("should run an ec2 instance", func() {
			_, err := ec2Client.RunEC2(awsEC2.RunInstancesInput{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if creating the instance fails", func() {
			expectedError := fmt.Errorf("some-error")
			ec2API.RunInstancesReturns(nil, expectedError)
			_, err := ec2Client.RunEC2(awsEC2.RunInstancesInput{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedError))
		})
	})
	Describe("TerminateEC2", func() {
		It("should terminate an ec2 instance", func() {
			_, err := ec2Client.TerminateEC2(awsEC2.TerminateInstancesInput{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if terminating the instance fails", func() {
			expectedError := fmt.Errorf("terminate-error")
			ec2API.TerminateInstancesReturns(nil, expectedError)
			_, err := ec2Client.TerminateEC2(awsEC2.TerminateInstancesInput{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedError))
		})
	})
	Describe("TagEC2", func() {
		It("should tag an ec2 instance", func() {
			_, err := ec2Client.TagEC2(aws.String("instanceID"), []*awsEC2.Tag{{}})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if tagging the instance fails", func() {
			expectedError := fmt.Errorf("tag-error")
			ec2API.CreateTagsReturns(nil, expectedError)
			_, err := ec2Client.TagEC2(aws.String("instanceID"), []*awsEC2.Tag{{}})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedError))
		})
	})
	Describe("IdentifyEC2", func() {
		It("should find at least one ec2 reservation", func() {
			ec2DescribeOutput := awsEC2.DescribeInstancesOutput{Reservations: []*awsEC2.Reservation{{Instances: []*awsEC2.Instance{{}}}}}
			ec2API.DescribeInstancesReturns(&ec2DescribeOutput, nil)

			result, err := ec2Client.IdentifyEC2("some-service-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(1))
		})

		It("should not return an error if there are no reservations", func() {
			ec2DescribeOutput := awsEC2.DescribeInstancesOutput{Reservations: []*awsEC2.Reservation{}}
			ec2API.DescribeInstancesReturns(&ec2DescribeOutput, nil)

			_, err := ec2Client.IdentifyEC2("some-service-id")
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
