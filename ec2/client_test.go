package ec2_test

import (
	"fmt"
	"time"

	"github.com/richardTowers/paas-drone-agent-broker/ec2"

	awsEC2 "github.com/aws/aws-sdk-go/service/ec2"
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
			err := ec2Client.RunEC2(awsEC2.RunInstancesInput{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if creating the instance fails", func() {
			expectedError := fmt.Errorf("some-error")
			ec2API.RunInstancesReturns(nil, expectedError)
			err := ec2Client.RunEC2(awsEC2.RunInstancesInput{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedError))
		})
	})

})
