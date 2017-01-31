package awsutil_test

import (
	"testing"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/aws-sdk-mock/mockcloudformation"
	"github.com/tleyden/aws-sdk-mock/mockec2"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func TestStopEC2Instances(t *testing.T) {

	mockCfn := NewMockCloudformationAPI()
	mockEc2 := mockec2.NewEC2APIMock()

	// TODO: return some ec2 instances and some non-ec2 instances
	mockCfn.On("DescribeStackResources", mock.Anything).Return(
		&cloudformation.DescribeStackResourcesOutput{},
		nil,
	).Once()

	cfnUtil := awsutil.NewCloudformationUtil(mockCfn, mockEc2)

	// We should expect calls to Stop the EC2 instance, exactly once

	cfnUtil.StopEC2Instances("fake_stack")

	// assert that all expectations met
	mockCfn.AssertExpectations(t)

}

func TestStopEc2InstanceStackResource(t *testing.T)  {

	mockCfn := NewMockCloudformationAPI()
	mockEc2 := mockec2.NewEC2APIMock()

	

	cfnUtil := awsutil.NewCloudformationUtil(mockCfn, mockEc2)

	stackResource := cloudformation.StackResource{
		ResourceType: stringPointer(awsutil.AWS_EC2_INSTANCE),
	}
	err := cfnUtil.StopEc2InstanceStackResource(stackResource)
	assert.NoError(t, err, "Error calling StopEc2InstanceStackResource")



}

// Creates a cloudformation API
func NewMockCloudformationAPI() *mockcloudformation.CloudFormationAPIMock {

	return mockcloudformation.NewCloudFormationAPIMock()


}

func stringPointer(s string) *string {
	return &s
}