package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func newCloudformationAPI(session *session.Session, region string) *cloudformation.CloudFormation {

	cloudformationService := cloudformation.New(session,
		&aws.Config{
			Region: aws.String(region),
		},
	)
	return cloudformationService

}

func newEC2API(session *session.Session, region string) *ec2.EC2 {

	ec2 := ec2.New(session,
		&aws.Config{
			Region: aws.String(region),
		},
	)
	return ec2

}

func newCloudwatchAPI(session *session.Session, region string) *cloudwatch.CloudWatch {

	cw := cloudwatch.New(session,
		&aws.Config{
			Region: aws.String(region),
		},
	)
	return cw

}
