package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/workspace"
)

func main() {
	server.Start(&ObjectStorage{})
}

var (
	ErrEmptyInstanceTypeForCloudDB = errors.New("empty instance type for cloud managed mysql instance")
	ErrEmptyCloudProviderType      = errors.New("empty cloud provider type in mysql module config")
)

// ObjectStorage implements the Kusion Module generator interface.
//
// Note that as an example of a Kusion Module, ObjectStorage consists of two components, one of which
// is a 'Service', which is used to generate a Kubernetes Service resource, and the other is a
// 'RandomePassword', which is used to generate a Terraform random_password resource.
//
// Typically, these two resources are not particularly related, but here they are combined to primarily
// illustrate how to develop a Kusion Module.
type ObjectStorage struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	Bucket string `yaml:"bucket,omitempty" json:"bucket,omitempty"`
}

// Generate implements the generation logic of objectStorage module, including a Kubernetes Service and
// a Terraform random_password resource.
func (objectStorage *ObjectStorage) Generate(_ context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("failed to generate objectStorage module: %v", r)
		}
	}()

	// ObjectStorage module does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		log.Info("ObjectStorage module does not exist in AppConfiguration configs")
	}

	// Get the complete objectStorage module configs.
	if err := objectStorage.CompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		log.Debugf("failed to get complete objectStorage module configs: %v", err)
		return nil, err
	}

	// Validate the completed objectStorage module configs.
	if err := objectStorage.ValidateConfig(); err != nil {
		log.Debugf("failed to validate the objectStorage module configs: %v", err)
		return nil, err
	}

	var resources []apiv1.Resource
	var patcher *apiv1.Patcher

	// var providerType string
	switch strings.ToLower(objectStorage.Type) {
	// case "local":
	// 	resources, patcher, err = mysql.GenerateLocalResources(request)
	case "cloud":
		providerType, err := GetCloudProviderType(request.PlatformConfig)
		if err != nil {
			return nil, err
		}

		switch strings.ToLower(providerType) {
		case "aws":
			resources, patcher, err = objectStorage.GenerateAWSResources(request)
			if err != nil {
				return nil, err
			}
		// case "alicloud":
		// 	resources, patcher, err = mysql.GenerateAlicloudResources(request)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		default:
			return nil, fmt.Errorf("unsupported cloud provider type: %s", providerType)
		}
	default:
		return nil, fmt.Errorf("unsupported mysql type: %s", objectStorage.Type)
	}

	// Return the Kusion generator response.
	return &module.GeneratorResponse{
		Resources: resources,
		Patcher:   patcher,
	}, nil
}

func GetCloudProviderType(platformConfig apiv1.GenericConfig) (string, error) {
	if platformConfig == nil {
		return "", workspace.ErrEmptyModuleConfigBlock
	}

	if cloud, ok := platformConfig["cloud"]; ok {
		return cloud.(string), nil
	}

	return "", ErrEmptyCloudProviderType
}

// CompleteConfig completes the objectStorage module configs with both devModuleConfig and platformModuleConfig.
func (objectStorage *ObjectStorage) CompleteConfig(devConfig apiv1.Accessory, platformConfig apiv1.GenericConfig) error {
	// Retrieve the config items the developers are concerned about.
	if devConfig != nil {
		devCfgYamlStr, err := yaml.Marshal(devConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(devCfgYamlStr, objectStorage); err != nil {
			return err
		}
	}

	// Retrieve the config items the platform engineers care about.
	if platformConfig != nil {
		platformCfgYamlStr, err := yaml.Marshal(platformConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(platformCfgYamlStr, objectStorage); err != nil {
			return err
		}
	}

	return nil
}

// ValidateConfig validates the completed objectStorage configs are valid or not.
func (objectStorage *ObjectStorage) ValidateConfig() error {
	return nil
}
