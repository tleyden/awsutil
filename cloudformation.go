package awsutil

import (

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Wraps AWS Cloudformation SDK API and provides additional utilities
type CloudformationUtil struct {
	cfnApi cloudformationiface.CloudFormationAPI
	ec2Api ec2iface.EC2API
}

// Create a new ClouformationUtil
func NewCloudformationUtil(cfnApi cloudformationiface.CloudFormationAPI, ec2Api ec2iface.EC2API) *CloudformationUtil {
	cfnUtil := &CloudformationUtil{
		cfnApi: cfnApi,
		ec2Api: ec2Api,
	}
	return cfnUtil
}

// Stop all EC2 Instances in a cloudformation stack
func (cfnu CloudformationUtil) StopEC2Instances(stackname string) error {

	params := &cloudformation.DescribeStackResourcesInput{
		StackName:         aws.String(stackname),
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
		err := cfnu.StopEc2InstanceStackResource(*stackResource)
		if err != nil {
			return err
		}

	}

	return nil

}

// Stop the EC2 Instance Stack Resource
func (cfnu CloudformationUtil) StopEc2InstanceStackResource(stackResource cloudformation.StackResource) error {

	if !IsStackResourceEc2Instance(stackResource) {
		return fmt.Errorf("Stack Resource [%+v] is not an EC2 instance.", stackResource)
	}

	stopInstancesInput := ec2.StopInstancesInput{
		InstanceIds: []*string{
			stackResource.LogicalResourceId,
		},
	}
	stopInstancesOutput, err := cfnu.ec2Api.StopInstances(&stopInstancesInput)
	if err != nil {
		return err
	}

	fmt.Printf("stopInstancesOutput: %v", stopInstancesOutput)

	return nil

}


// Creates a cloudformation API
func NewCloudformationAPI(session *session.Session, region string) *cloudformation.CloudFormation {

	cloudformationService := cloudformation.New(session,
		&aws.Config{
			Region: aws.String(region),
		},
	)
	return cloudformationService

}

// Is the StackResource parameter an EC2 instance?
func IsStackResourceEc2Instance(stackResource cloudformation.StackResource) bool {
	return *stackResource.ResourceType == AWS_EC2_INSTANCE
}

