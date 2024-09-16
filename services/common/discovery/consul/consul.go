package consul

import (
	"context"
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(address, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = address

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &Registry{
		client: client,
	}, nil
}

func (r *Registry) RegisterService(ctx context.Context, serviceName, serviceID, serviceAddress string, servicePort int, tags []string) error {
	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePort,
		Tags:    tags,
		Check: &consul.AgentServiceCheck{
			TTL:                            "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	err := r.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	return nil
}

func (r *Registry) DeregisterService(ctx context.Context, serviceID string) error {
	err := r.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	return nil
}

func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*consul.ServiceEntry, error) {
	services, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return services, nil
}

func (r *Registry) HealthCheck(serviceID, serviceName string) error {
	checkID := "service:" + serviceID
	return r.client.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
}
