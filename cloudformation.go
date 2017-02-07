package awsutil

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	// "github.com/y0ssar1an/q"
	"log"
	"strings"
)

// CloudformationUtil wraps the AWS Cloudformation SDK API and provides additional utilities
type CloudformationUtil struct {
	cfnApi cloudformationiface.CloudFormationAPI
	ec2Api ec2iface.EC2API
}

// NewCloudformationUtilFromRegion creates a CloudformationUtil from a region
func NewCloudformationUtilFromRegion(region string) (*CloudformationUtil, error) {

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	cfnApi := newCloudformationAPI(session, region)
	ec2Api := newEC2API(session, region)

	return NewCloudformationUtil(cfnApi, ec2Api)

}

// NewCloudformationUtil creates a CloudformationUtil given the underlying AWS api implementations (or mocks)
func NewCloudformationUtil(cfnApi cloudformationiface.CloudFormationAPI, ec2Api ec2iface.EC2API) (*CloudformationUtil, error) {

	cnfUtil := &CloudformationUtil{
		cfnApi: cfnApi,
		ec2Api: ec2Api,
	}

	return cnfUtil, nil

}

// StopEC2Instances stops all EC2 Instances in a cloudformation stack
func (cfnu CloudformationUtil) StopEC2Instances(stackname string) error {

	f := func(s *cloudformation.StackResource, c CloudformationUtil) error {
		err := cfnu.StopEc2InstanceForStackResource(*s)
		return err
	}
	err := cfnu.startOrStop(stackname, f)
	return err

}

// StartEC2Instances starts/restarts all EC2 Instances in a cloudformation stack
func (cfnu CloudformationUtil) StartEC2Instances(stackname string) error {
	f := func(s *cloudformation.StackResource, c CloudformationUtil) error {
		err := cfnu.StartEc2InstanceForStackResource(*s)
		return err
	}
	err := cfnu.startOrStop(stackname, f)
	return err
}

// StopEc2InstanceForStackResource stops the EC2 Instance for the given Stack Resource
func (cfnu CloudformationUtil) StopEc2InstanceForStackResource(stackResource cloudformation.StackResource) error {

	if !IsStackResourceEc2Instance(stackResource) {
		return fmt.Errorf("Stack Resource [%+v] is not an EC2 instance.", stackResource)
	}

	stopInstancesInput := ec2.StopInstancesInput{
		InstanceIds: []*string{
			stackResource.PhysicalResourceId,
		},
	}

	_, err := cfnu.ec2Api.StopInstances(&stopInstancesInput)
	if err != nil {
		return err
	}

	return nil

}

// StartEc2InstanceForStackResource starts the EC2 Instance in given Stack Resource
func (cfnu CloudformationUtil) StartEc2InstanceForStackResource(stackResource cloudformation.StackResource) error {
	if !IsStackResourceEc2Instance(stackResource) {
		return fmt.Errorf("Stack Resource [%+v] is not an EC2 instance.", stackResource)
	}

	startInstancesInput := ec2.StartInstancesInput{
		InstanceIds: []*string{
			stackResource.PhysicalResourceId,
		},
	}

	_, err := cfnu.ec2Api.StartInstances(&startInstancesInput)
	if err != nil {
		return err
	}

	return nil

}

// InCloudformation checks whether the given instance id is part of a
// Cloudformation.  If it is, it returns a boolean val set to true, and a
// a *cloudformation.StackResource.  Otherwise it returns false/nil
func (cfnu CloudformationUtil) InCloudformation(instanceId string) (bool, *cloudformation.StackResource, error) {

	params := &cloudformation.DescribeStackResourcesInput{
		PhysicalResourceId: StringPointer(instanceId),
	}
	dsrOutput, err := cfnu.cfnApi.DescribeStackResources(params)
	if err != nil {

		// if we get a "does not exist" then just absorb the error and don't treat it as an error condition
		// since it's expected if the instance ID is not present in any cloudformation stacks
		if strings.Contains(err.Error(), "does not exist") {
			return false, nil, nil
		}

		// it's an unexpected error, let's return it to the caller to bring to attention
		return false, nil, err

	}
	if len(dsrOutput.StackResources) == 0 {
		return false, nil, nil
	}
	if len(dsrOutput.StackResources) > 1 {
		log.Printf(
			"Warning: got %d stackresources for instanceid: %v. Expected 0 or 1",
			len(dsrOutput.StackResources),
			instanceId,
		)
	}

	stackResource := dsrOutput.StackResources[0]

	if stackResource.LogicalResourceId == nil {
		return false, nil, nil
	}

	if stackResource.StackId == nil {
		return false, nil, nil
	}

	if stackResource.StackName == nil {
		return false, nil, nil
	}


	// Defensive check -- should never happen, because AWS should always return us a
	// stack resource that matches our filter.
	if *stackResource.PhysicalResourceId != instanceId {
		return false, stackResource, fmt.Errorf("Stack resource physicalresource id [%v] != expected [%v]",
			*stackResource.PhysicalResourceId,
			instanceId,
		)
	}

	return true, stackResource, nil
}


// IsStackResourceEc2Instance checks whether the StackResource parameter an EC2 instance
func IsStackResourceEc2Instance(stackResource cloudformation.StackResource) bool {
	return *stackResource.ResourceType == AWS_EC2_INSTANCE
}

func (cfnu CloudformationUtil) startOrStop(stackname string, f func(*cloudformation.StackResource, CloudformationUtil) error) error {

	params := &cloudformation.DescribeStackResourcesInput{
		StackName: aws.String(stackname),
	}

	describeStackResourcesOut, err := cfnu.cfnApi.DescribeStackResources(params)
	if err != nil {
		return err
	}

	for _, stackResource := range describeStackResourcesOut.StackResources {

		// if it's not an EC2 instance, ignore it
		if !IsStackResourceEc2Instance(*stackResource) {
			continue
		}

		// otherwise stop it
		err := f(stackResource, cfnu)
		if err != nil {
			return err
		}

	}

	return nil
}