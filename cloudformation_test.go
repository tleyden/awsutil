package awsutil_test

import (
	"testing"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/aws-sdk-mock"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	// "github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"
)

func TestStopEC2Instances(t *testing.T) {

	mockCfn := NewMockCloudformationAPI()

	// TODO: return some ec2 instances and some non-ec2 instances
	mockCfn.On("DescribeStackResources", mock.Anything).Return(
		&cloudformation.DescribeStackResourcesOutput{},
		nil,
	).Once()

	cfnUtil := awsutil.NewCloudformationUtil(mockCfn)

	// We should expect calls to Stop the EC2 instance, exactly once

	cfnUtil.StopEC2Instances("fake_stack")

	// assert that all expectations met
	mockCfn.AssertExpectations(t)

}

// Creates a cloudformation API
func NewMockCloudformationAPI() *mockcloudformation.CloudFormationAPIMock {

	return mockcloudformation.NewCloudFormationAPIMock()


}