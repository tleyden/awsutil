package awsutil_test

import (
	"testing"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/aws-sdk-mock/mockcloudformation"
	"github.com/tleyden/aws-sdk-mock/mockec2"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func TestStopEC2Instances(t *testing.T) {

	mockCfn := NewMockCloudformationAPI()
	mockEc2 := mockec2.NewEC2APIMock()

	// Mock cloudformation returns stack with some ec2 instances and some non-ec2 instances
	mockCfn.On("DescribeStackResources", mock.Anything).Return(
		&cloudformation.DescribeStackResourcesOutput{
			StackResources: []*cloudformation.StackResource{
				&cloudformation.StackResource{
					ResourceType: awsutil.StringPointer(awsutil.AWS_EC2_INSTANCE),
				},
				&cloudformation.StackResource{
					ResourceType: awsutil.StringPointer(awsutil.AWS_EC2_HOST),
				},

			},
		},
		nil,
	).Once()

	// Expect a call to ec2 StopInstances
	mockEc2.On("StopInstances", mock.Anything).Return(
		&ec2.StopInstancesOutput{},
		nil,
	).Once()

	// Create CloudformationUtil which is the struct being tested
	cfnUtil, err := awsutil.NewCloudformationUtil(mockCfn, mockEc2)
	assert.NoError(t, err, "Error creating NewCloudformationUtil")

	// Tell it to stop all EC2 instances in the fake Cloudformation Stack
	cfnUtil.StopEC2Instances("fake_stack")

	// assert that all expectations met
	mockCfn.AssertExpectations(t)
	mockEc2.AssertExpectations(t)


}

func TestStopEc2InstanceStackResource(t *testing.T)  {

	mockInstanceId := "i-mock"

	mockCfn := NewMockCloudformationAPI()
	mockEc2 := mockec2.NewEC2APIMock()

	// The mock ec2 API is expecting to get this as the parameter to
	// the ec2Api.StopInstances invocation
	expectedStopInstancesInput := ec2.StopInstancesInput{
		InstanceIds: []*string{
			&mockInstanceId,
		},
	}

	mockEc2.On("StopInstances", &expectedStopInstancesInput).Return(
		&ec2.StopInstancesOutput{},
		nil,
	).Once()

	cfnUtil, err := awsutil.NewCloudformationUtil(mockCfn, mockEc2)
	assert.NoError(t, err, "Error calling NewCloudformationUtil")

	stackResource := cloudformation.StackResource{
		ResourceType: awsutil.StringPointer(awsutil.AWS_EC2_INSTANCE),
		PhysicalResourceId: &mockInstanceId,
	}

	err = cfnUtil.StopEc2InstanceStackResource(stackResource)
	assert.NoError(t, err, "Error calling StopEc2InstanceStackResource")

	mockEc2.AssertExpectations(t)


}

// Creates a cloudformation API
func NewMockCloudformationAPI() *mockcloudformation.CloudFormationAPIMock {

	return mockcloudformation.NewCloudFormationAPIMock()

}

