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
)

// Wraps AWS Cloudformation SDK API and provides additional utilities
type CloudformationUtil struct {
	cfnApi cloudformationiface.CloudFormationAPI
	ec2Api ec2iface.EC2API
}

func NewCloudformationUtilFromRegion(region string) (*CloudformationUtil, error) {

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	cfnApi := newCloudformationAPI(session, region)
	ec2Api := newEC2API(session, region)

	return NewCloudformationUtil(cfnApi, ec2Api)

}

func NewCloudformationUtil(cfnApi cloudformationiface.CloudFormationAPI, ec2Api ec2iface.EC2API) (*CloudformationUtil, error) {

	cnfUtil := &CloudformationUtil{
		cfnApi: cfnApi,
		ec2Api: ec2Api,
	}

	return cnfUtil, nil

}

// Stop all EC2 Instances in a cloudformation stack
func (cfnu CloudformationUtil) StopEC2Instances(stackname string) error {

	f := func(s *cloudformation.StackResource, c CloudformationUtil) error {
		err := cfnu.StopEc2InstanceForStackResource(*s)
		return err
	}
	err := cfnu.startOrStop(stackname, f)
	return err

}

// Start/restart all EC2 Instances in a cloudformation stack
func (cfnu CloudformationUtil) StartEC2Instances(stackname string) error {
	f := func(s *cloudformation.StackResource, c CloudformationUtil) error {
		err := cfnu.StartEc2InstanceForStackResource(*s)
		return err
	}
	err := cfnu.startOrStop(stackname, f)
	return err
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

// Stop the EC2 Instance Stack Resource
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

// Stop the EC2 Instance Stack Resource
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

// Is the StackResource parameter an EC2 instance?
func IsStackResourceEc2Instance(stackResource cloudformation.StackResource) bool {
	return *stackResource.ResourceType == AWS_EC2_INSTANCE
}
