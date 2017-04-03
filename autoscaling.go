package awsutil

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	// "github.com/y0ssar1an/q"
	"strings"
)

// AutoScalingUtil wraps the AWS AutoScaling SDK API and provides additional utilities
type AutoScalingUtil struct {
	asgAPI autoscalingiface.AutoScalingAPI
	ec2Api ec2iface.EC2API
}

// NewAutoScalingUtilFromRegion creates a AutoScalingUtil from a region
func NewAutoScalingUtilFromRegion(region string) (*AutoScalingUtil, error) {

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	asgAPI := newAutoScalingAPI(session, region)
	ec2Api := newEC2API(session, region)

	return NewAutoScalingUtil(asgAPI, ec2Api)

}

// NewAutoScalingUtil creates a AutoScalingUtil given the underlying AWS api implementations (or mocks)
func NewAutoScalingUtil(asgAPI autoscalingiface.AutoScalingAPI, ec2Api ec2iface.EC2API) (*AutoScalingUtil, error) {

	cnfUtil := &AutoScalingUtil{
		asgAPI: asgAPI,
		ec2Api: ec2Api,
	}

	return cnfUtil, nil

}

// InAutoScaling checks whether the given instance id is part of an
// AutoScaling group. If it is, it returns a nil error and the asgName.
func (cfnu AutoScalingUtil) InAutoScaling(instanceId string) (in bool, asgName string, err error) {

	params := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId), // Required
		},
	}
	dsrOutput, err := cfnu.asgAPI.DescribeAutoScalingInstances(params)
	if err != nil {

		// if we get a "does not exist" then just absorb the error and don't treat it as an error condition
		// since it's expected if the instance ID is not present in any autoscaling groups
		if strings.Contains(err.Error(), "does not exist") {
			return false, "", nil
		}

		// it's an unexpected error, let's return it to the caller to bring to attention
		return false, "", err

	}
	if len(dsrOutput.AutoScalingInstances) == 0 {
		return false, "", nil
	}

	asgResource := dsrOutput.AutoScalingInstances[0]

	if asgResource.AutoScalingGroupName == nil {
		return false, "", fmt.Errorf("asgResource is missing AutoScalingGroupName: %+v", *asgResource)
	}

	return true, *asgResource.AutoScalingGroupName, nil
}

// ASG_ARN_fromName returns the ARN of an autoscaling group given its name
func (cfnu AutoScalingUtil) ASG_ARN_fromName(asgName string) (asgARN string, err error) {

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName), // Required
		},
	}
	dsrOutput, err := cfnu.asgAPI.DescribeAutoScalingGroups(params)
	if err != nil {

		// if we get a "does not exist" then just absorb the error and don't treat it as an error condition
		// since it's expected if the instance ID is not present in any autoscaling groups
		if strings.Contains(err.Error(), "does not exist") {
			return "", nil
		}

		// it's an unexpected error, let's return it to the caller to bring to attention
		return "", err

	}
	if len(dsrOutput.AutoScalingGroups) == 0 {
		return "", nil
	}

	asgGroup := dsrOutput.AutoScalingGroups[0]

	if asgGroup.AutoScalingGroupARN == nil {
		return "", fmt.Errorf("asgGroup is missing AutoScalingGroupARN: %+v", *asgGroup)
	}

	return *asgGroup.AutoScalingGroupARN, nil
}
