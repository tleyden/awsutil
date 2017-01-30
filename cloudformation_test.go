package awsutil_test

import (
	"testing"
	"github.com/tleyden/awsutil"
	"github.com/tleyden/aws-sdk-mock/mockservice/mockcloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	// "github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"
)

func TestStopEC2Instances(t *testing.T) {

	mockCfn := NewMockCloudformationAPI()

	/*mockCfn.On("DescribeStackResources", &cloudformation.DescribeStackResourcesInput{StackName: aws.String("fake_stack")}).Return(
		cloudformation.DescribeStackResourcesOutput{},
		nil,
	).Once()*/

	mockCfn.On("DescribeStackResources", mock.Anything).Return(
		&cloudformation.DescribeStackResourcesOutput{},
		nil,
	).Once()

	cfnUtil := awsutil.NewCloudformationUtil(mockCfn)

	cfnUtil.StopEC2Instances("fake_stack")
}

// Creates a cloudformation API
func NewMockCloudformationAPI() *mockcloudformation.CloudFormationAPIMock {

	return mockcloudformation.NewCloudFormationAPIMock()


}