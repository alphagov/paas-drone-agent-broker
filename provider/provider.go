package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	provideriface "github.com/alphagov/paas-go/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pivotal-cf/brokerapi"
	ec2API "github.com/alphagov/paas-drone-agent-broker/ec2"
	"log"
	"strings"
	template2 "text/template"
)

type DroneAgentConfig struct {
	RPCSecret      string `json:"server_secret"`
	RPCServer      string `json:"server_address"`
	RunnerCapacity int    `json:"runner_capacity"`
	LogsDebug      bool   `json:"debug_logs"`
}

type AWSConfig struct {
	AWSRegion       string `json:"aws_region"`
	SecurityGroupID string `json:"security_group_id"`
}

type DroneAgentProvider struct {
	Client          ec2API.Client
	Config          *AWSConfig
	SecurityGroupID string
}

func NewDroneAgentProvider(configJSON []byte) (provideriface.ServiceProvider, error) {
	config := &AWSConfig{
		AWSRegion:       "eu-west-2",
		SecurityGroupID: "",
	}
	err := json.Unmarshal(configJSON, &config)
	if err != nil {
		return nil, err
	}
	client := ec2API.NewEC2Client(config.AWSRegion)
	return &DroneAgentProvider{
		Client:          client,
		Config:          config,
		SecurityGroupID: config.SecurityGroupID,
	}, nil
}

func (s *DroneAgentProvider) RunInstance(provisionData provideriface.ProvisionData) (awsInstanceID string, error error) {
	var agentConfig DroneAgentConfig

	err := json.Unmarshal(provisionData.Details.RawParameters, &agentConfig)

	template, err := template2.ParseFiles("provider/userdata.txt")
	if err != nil {
		return "", err
	}
	var userData bytes.Buffer
	err = template.Execute(&userData, agentConfig)
	if err != nil {
		return "", err
	}
	b64UserData := base64.StdEncoding.EncodeToString(userData.Bytes())

	runInstancesInput := ec2.RunInstancesInput{
		ImageId:          aws.String("ami-0016c65679adc75f5"),
		SecurityGroupIds: aws.StringSlice([]string{s.Config.SecurityGroupID}),
		InstanceType:     aws.String(provisionData.Plan.Name),
		UserData:         &b64UserData,
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
	}
	provisionResponse, err := s.Client.RunEC2(runInstancesInput)

	if err != nil {
		return "", err
	}

	awsInstanceID = aws.StringValue(provisionResponse.Instances[0].InstanceId)

	_, err = s.Client.TagEC2(aws.String(awsInstanceID), []*ec2.Tag{{
		Key:   aws.String("service_instance_ref"),
		Value: aws.String(provisionData.InstanceID),
	},
		{
			Key:   aws.String("org_guid"),
			Value: aws.String(provisionData.Details.OrganizationGUID),
		},
		{
			Key:   aws.String("space_guid"),
			Value: aws.String(provisionData.Details.SpaceGUID),
		},
		{
			Key:   aws.String("service_type"),
			Value: aws.String("drone_agent"),
		},
	})
	return awsInstanceID, err
}

func (s *DroneAgentProvider) Provision(ctx context.Context, provisionData provideriface.ProvisionData) (
	dashboardURL, operationData string, isAsync bool, err error) {
	reservations, err := s.Client.IdentifyEC2(provisionData.InstanceID)
	if len(reservations) != 0 {
		return "", "", false, errors.New(fmt.Sprintf("An instance with ID %v already exists", provisionData.InstanceID))
	}
	awsInstanceID, err := s.RunInstance(provisionData)
	if err != nil {
		if awsInstanceID == "" {
			return "", "", false, errors.New(fmt.Sprint(err))
		}
		terminateInstanceInput := ec2.TerminateInstancesInput{
			InstanceIds: []*string{aws.String(awsInstanceID)},
		}
		_, err = s.Client.TerminateEC2(terminateInstanceInput)
		if err != nil {
			return "", awsInstanceID, true, errors.New("Tagging failed, then terminating the new instance failed.")
		}

		return "", awsInstanceID, true, errors.New("Tagging failed, terminating instance")
	}

	return "", awsInstanceID, true, err
}

func (s *DroneAgentProvider) Deprovision(ctx context.Context, deprovisionData provideriface.DeprovisionData) (
	operationData string, isAsync bool, err error) {
	serviceRef := deprovisionData.InstanceID
	reservations, err := s.Client.IdentifyEC2(serviceRef)
	var instanceIDs []string
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, aws.StringValue(instance.InstanceId))
		}
	}
	terminateInstanceInput := ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(instanceIDs),
	}
	_, err = s.Client.TerminateEC2(terminateInstanceInput)
	if err != nil {
		return "", false, errors.New(fmt.Sprintf("No instances with ID %v exist", serviceRef))
	}
	return strings.Join(instanceIDs, ","), true, err
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
	serviceRef := updateData.InstanceID
	reservations, err := s.Client.IdentifyEC2(serviceRef)
	var instancesTerminated []string
	var instancesCreated []string
	var errorInstancesCreated []string
	var errorInstancesTerminated []string
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceType != aws.String(updateData.Plan.Name) {
				provisionData := provideriface.ProvisionData{
					InstanceID: updateData.InstanceID,
					Details: brokerapi.ProvisionDetails{
						ServiceID:        updateData.Details.ServiceID,
						PlanID:           updateData.Details.PlanID,
						OrganizationGUID: updateData.Details.PreviousValues.OrgID,
						SpaceGUID:        updateData.Details.PreviousValues.SpaceID,
						RawContext:       updateData.Details.RawContext,
						RawParameters:    updateData.Details.RawParameters,
					},
					Service: updateData.Service,
					Plan:    updateData.Plan,
				}
				awsInstanceID, err := s.RunInstance(provisionData)
				instancesCreated = append(instancesCreated, awsInstanceID)
				terminateInstanceInput := ec2.TerminateInstancesInput{
					InstanceIds: []*string{instance.InstanceId},
				}
				_, err = s.Client.TerminateEC2(terminateInstanceInput)
				if err != nil {
					log.Printf("Termination of %v failed.", aws.StringValue(instance.InstanceId))
					errorInstancesTerminated = append(errorInstancesTerminated, aws.StringValue(instance.InstanceId))
				}
				instancesTerminated = append(instancesTerminated, aws.StringValue(instance.InstanceId))
			}
		}
	}
	if len(errorInstancesCreated) != 0 || len(errorInstancesTerminated) != 0 {
		outputData := fmt.Sprintf("Terminated: %v, Created: %v", strings.Join(instancesTerminated, ","), strings.Join(instancesCreated, ","))
		return "", false, errors.New(outputData)
	}
	outputData := fmt.Sprintf("Terminated: %v, Created: %v", strings.Join(instancesTerminated, ","), strings.Join(instancesCreated, ","))
	return outputData, false, nil
}

func (s *DroneAgentProvider) LastOperation(ctx context.Context, lastOperationData provideriface.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error) {
	return "", "", errors.New("not implemented")
}
