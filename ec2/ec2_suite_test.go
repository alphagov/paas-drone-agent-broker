package ec2_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEC2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EC2 Suite")
}
