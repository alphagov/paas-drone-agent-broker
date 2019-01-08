package provider

import (
	"context"
	"errors"
	provideriface "github.com/alphagov/paas-go/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pivotal-cf/brokerapi"
	ec2API "github.com/richardTowers/paas-drone-agent-broker/ec2"
)

type DroneAgentConfig struct {
	RPCSecret      string
	RPCServer      string
	RunnerCapacity int
	LogsDebug      bool
}
type DroneAgentProvider struct {
	Client ec2API.Client
	Config []byte
}

func NewDroneAgentProvider(config []byte) (provideriface.ServiceProvider, error) {
	client := ec2API.NewEC2Client("eu-west-2")
	return &DroneAgentProvider{
		Client: client,
		Config: config,
	}, nil
}

func (s *DroneAgentProvider) Provision(ctx context.Context, provisionData provideriface.ProvisionData) (
	dashboardURL, operationData string, isAsync bool, err error) {
	runInstancesInput := ec2.RunInstancesInput{
		ImageId:          aws.String("ami-0016c65679adc75f5"),
		SecurityGroupIds: aws.StringSlice([]string{"sg-0a1b0216ef7084cc0"}),
		InstanceType:     aws.String("t2.small"),
		UserData:         aws.String("IyEvYmluL3NoCgpkb2NrZXIgcnVuIFxcCiAgLS12b2x1bWU9L3Zhci9ydW4vZG9ja2VyLnNvY2s6L3Zhci9ydW4vZG9ja2VyLnNvY2sgXFwKICAtLXZvbHVtZT0vdmFyL2xpYi9kcm9uZTovZGF0YSBcXAogIC0tZW52PURST05FX1JQQ19TRVJWRVI9JGRyb25lX3JwY19zZXJ2ZXIgXFwKICAtLWVudj1EUk9ORV9SUENfU0VDUkVUPSRkcm9uZV9ycGNfc2VjcmV0IFxcCiAgLS1lbnY9RFJPTkVfUlVOTkVSX0NBUEFDSVRZPTIgXFwKICAtLWVudj1EUk9ORV9MT0dTX0RFQlVHPXRydWUgXFwKICAtLXJlc3RhcnQ9YWx3YXlzIFxcCiAgLS1kZXRhY2g9dHJ1ZSBcXAogIC0tbmFtZT1kcm9uZSBcXAogIGRyb25lL2FnZW50OjEuMC4wLXJjLjE="),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
	}
	reservation, err := s.Client.RunEC2(runInstancesInput)
	if err != nil {
		return "", "", false, err
	}

	return aws.StringValue(reservation.Instances[0].PublicIpAddress), "", true, err
}

func (s *DroneAgentProvider) Deprovision(ctx context.Context, deprovisionData provideriface.DeprovisionData) (
	operationData string, isAsync bool, err error) {
	return "", false, errors.New("not implemented")
}

func (s *DroneAgentProvider) Bind(ctx context.Context, bindData provideriface.BindData) (
	binding brokerapi.Binding, err error) {
	return brokerapi.Binding{}, errors.New("not implemented")
}

func (s *DroneAgentProvider) Unbind(ctx context.Context, unbindData provideriface.UnbindData) (
	spec brokerapi.UnbindSpec, err error) {
	return brokerapi.UnbindSpec{}, errors.New("not implemented")
}

func (s *DroneAgentProvider) Update(ctx context.Context, updateData provideriface.UpdateData) (
	operationData string, isAsync bool, err error) {
	return "", false, errors.New("not implemented")
}

func (s *DroneAgentProvider) LastOperation(ctx context.Context, lastOperationData provideriface.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error) {
	return "", "", errors.New("not implemented")
}
