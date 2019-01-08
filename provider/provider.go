package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	provideriface "github.com/alphagov/paas-go/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pivotal-cf/brokerapi"
	ec2API "github.com/richardTowers/paas-drone-agent-broker/ec2"
	template2 "text/template"
)

type DroneAgentConfig struct {
	RPCSecret      string `json:"server_secret"`
	RPCServer      string `json:"server_address"`
	RunnerCapacity int    `json:"runner_capacity"`
	LogsDebug      bool   `json:"debug_logs"`
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
	var agentConfig DroneAgentConfig

	err = json.Unmarshal(provisionData.Details.RawParameters, &agentConfig)

	template, err := template2.ParseFiles("provider/userdata.txt")
	if err != nil {
		return "", "", false, err
	}
	var userData bytes.Buffer
	err = template.Execute(&userData, agentConfig)
	if err != nil {
		return "", "", false, err
	}
	b64UserData := base64.StdEncoding.EncodeToString(userData.Bytes())

	runInstancesInput := ec2.RunInstancesInput{
		ImageId:          aws.String("ami-0016c65679adc75f5"),
		SecurityGroupIds: aws.StringSlice([]string{"sg-0a1b0216ef7084cc0"}),
		InstanceType:     aws.String("t2.small"),
		UserData:         &b64UserData,
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
	}
	provisionResponse, err := s.Client.RunEC2(runInstancesInput)

	if err != nil {
		return "", "", false, err
	}

	instanceID := provisionResponse.Instances[0].InstanceId

	_, err = s.Client.TagEC2(instanceID, []*ec2.Tag{&ec2.Tag{
		Key:   aws.String("service_instance_ref"),
		Value: aws.String(provisionData.Service.ID),
	},
		&ec2.Tag{
			Key:   aws.String("org_guid"),
			Value: aws.String(provisionData.Details.OrganizationGUID),
		},
	})
	if err != nil {
		terminateInstanceInput := ec2.TerminateInstancesInput{
			InstanceIds: []*string{instanceID},
		}
		s.Client.TerminateEC2(terminateInstanceInput)
		return "", aws.StringValue(instanceID), true, errors.New("Tagging failed, terminating instance")
	}

	return "", aws.StringValue(instanceID), true, err
}

func (s *DroneAgentProvider) Deprovision(ctx context.Context, deprovisionData provideriface.DeprovisionData) (
	operationData string, isAsync bool, err error) {
	terminateInstanceInput := ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{"i-0491de23c4fb4a1c9"}),
	}
	_, err = s.Client.TerminateEC2(terminateInstanceInput)
	if err != nil {
		return "", false, err
	}
	return "", false, err
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
	return "", false, brokerapi.ErrPlanChangeNotSupported
}

func (s *DroneAgentProvider) LastOperation(ctx context.Context, lastOperationData provideriface.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error) {
	return "", "", errors.New("not implemented")
}
