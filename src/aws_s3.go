package main

import (
	"errors"
	"os"

	v1 "k8s.io/api/core/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/modules"
)

var ErrEmptyAWSProviderRegion = errors.New("empty aws provider region")

var (
	awsRegionEnv     = "AWS_REGION"
	awsSecurityGroup = "aws_security_group"
	awsDBInstance    = "aws_db_instance"
	awsS3Bucket      = "aws_s3_bucket"
)

var defaultAWSProviderCfg = module.ProviderConfig{
	Source:  "hashicorp/aws",
	Version: "5.0.1",
}

type awsSecurityGroupTraffic struct {
	CidrBlocks     []string `yaml:"cidr_blocks" json:"cidr_blocks"`
	Description    string   `yaml:"description" json:"description"`
	FromPort       int      `yaml:"from_port" json:"from_port"`
	IPv6CIDRBlocks []string `yaml:"ipv6_cidr_blocks" json:"ipv6_cidr_blocks"`
	PrefixListIDs  []string `yaml:"prefix_list_ids" json:"prefix_list_ids"`
	Protocol       string   `yaml:"protocol" json:"protocol"`
	SecurityGroups []string `yaml:"security_groups" json:"security_groups"`
	Self           bool     `yaml:"self" json:"self"`
	ToPort         int      `yaml:"to_port" json:"to_port"`
}

// GenerateAWSResources generates the AWS provided ObjectStorage database instance.
func (objectStorage *ObjectStorage) GenerateAWSResources(request *module.GeneratorRequest) ([]apiv1.Resource, *apiv1.Patcher, error) {
	var resources []apiv1.Resource

	// Set the AWS provider with the default provider config.
	awsProviderCfg := defaultAWSProviderCfg

	// Get the AWS Terraform provider region, which should not be empty.
	var region string
	if region = module.TerraformProviderRegion(awsProviderCfg); region == "" {
		region = os.Getenv(awsRegionEnv)
	}
	if region == "" {
		return nil, nil, ErrEmptyAWSProviderRegion
	}

	// Build random_password resource.
	// randomPasswordRes, randomPasswordID, err := objectStorage.GenerateRandomPasswordResource(request)
	// resources = append(resources, *randomPasswordRes)

	// // Build aws_security_group resource.
	// awsSecurityGroupRes, awsSecurityGroupID, err := objectStorage.generateAWSSecurityGroup(awsProviderCfg, region)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// resources = append(resources, *awsSecurityGroupRes)

	// Build aws_db_instance resource.
	// awsDBInstance, awsDBInstanceID, err := objectStorage.generateAWSDBInstance(awsProviderCfg, region, randomPasswordID, awsSecurityGroupID)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// resources = append(resources, *awsDBInstance)

	// Build aws_db_instance resource.
	// awsS3Bucket, awsS3BucketID, err := objectStorage.generateAWSS3Bucket(awsProviderCfg, region)
	awsS3Bucket, awsS3BucketID, err := objectStorage.generateAWSS3Bucket(awsProviderCfg, region)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *awsS3Bucket)

	hostAddress := modules.KusionPathDependency(awsS3BucketID, "address")
	// password := modules.KusionPathDependency(randomPasswordID, "result")

	envVars := []v1.EnvVar{
		{
			Name:  "KUSION_AWS_S3_BUCKET_ADDRESS",
			Value: modules.KusionPathDependency(hostAddress, "result"),
		},
	}
	patcher := &apiv1.Patcher{
		Environments: envVars,
	}

	// hostAddress := modules.KusionPathDependency(awsDBInstanceID, "address")
	// password := modules.KusionPathDependency(randomPasswordID, "result")

	// Build Kubernetes Secret with the hostAddress, username and password of the AWS provided ObjectStorage instance,
	// and inject the credentials as the environment variable patcher.
	// dbSecret, patcher, err := objectStorage.GenerateDBSecret(request, hostAddress, objectStorage.Username, password)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// resources = append(resources, *dbSecret)

	return resources, patcher, nil
}

// generateAWSS3 generates aws_s3_bucket resource for the AWS provided ObjectStorage database instance.
func (objectStorage *ObjectStorage) generateAWSS3Bucket(awsProviderCfg module.ProviderConfig, region string) (*apiv1.Resource, string, error) {
	// SecurityIPs should be in the format of IP address or Classes Inter-Domain
	// Routing (CIDR) mode.
	// for _, ip := range objectStorage.SecurityIPs {
	// 	if !IsIPAddress(ip) && !IsCIDR(ip) {
	// 		return nil, "", fmt.Errorf("illegal security ip format: %s", ip)
	// 	}
	// }

	resAttrs := map[string]interface{}{
		"bucket": objectStorage.Bucket,
	}

	id, err := module.TerraformResourceID(awsProviderCfg, awsS3Bucket, objectStorage.Bucket)
	if err != nil {
		return nil, "", err
	}

	awsProviderCfg.ProviderMeta = map[string]any{"region": region}
	resource, err := module.WrapTFResourceToKusionResource(awsProviderCfg, awsS3Bucket, id, resAttrs, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}
