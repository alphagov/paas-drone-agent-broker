package provider

import (
	"context"
	"errors"

	provideriface "github.com/alphagov/paas-go/provider"
	"github.com/pivotal-cf/brokerapi"
)

type DroneAgentProvider struct{}

func NewDroneAgentProvider(config []byte) (provideriface.ServiceProvider, error) {
	return &DroneAgentProvider{}, nil
}

func (s *DroneAgentProvider) Provision(ctx context.Context, provisionData provideriface.ProvisionData) (
	dashboardURL, operationData string, isAsync bool, err error) {
	return "", "", false, errors.New("not implemented")
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
