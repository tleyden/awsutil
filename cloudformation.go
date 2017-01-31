package awsutil

import (

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
)

// Wraps AWS Cloudformation SDK API and provides additional utilities
type CloudformationUtil struct {
	cfnApi cloudformationiface.CloudFormationAPI
}

// Create a new ClouformationUtil
func NewCloudformationUtil(cfnApi cloudformationiface.CloudFormationAPI) *CloudformationUtil {
	cfnUtil := &CloudformationUtil{
		cfnApi: cfnApi,
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
		err := StopEc2InstanceStackResource(*stackResource)
		if err != nil {
			return err
		}

	}

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
	return *stackResource.ResourceType == "AWS::EC2::Instance"
}

// Stop the EC2 Instance Stack Resource
func StopEc2InstanceStackResource(stackResource cloudformation.StackResource) error {

	if !IsStackResourceEc2Instance(stackResource) {
		return fmt.Errorf("Stack Resource [%+v] is not an EC2 instance.", stackResource)
	}

	// TODO

	return nil

}