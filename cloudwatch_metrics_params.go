package awsutil

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
)

type FetchRunningEc2InstanceMetricsInput struct {
	RecentTimeSpan time.Duration
	Namespace      string
	MetricName     string
}

type FetchRunningEc2InstanceMetricsOuput struct {
	Namespace     string
	Ec2InstanceId string
	MetricName    string
	Metrics       *cloudwatch.GetMetricStatisticsOutput
}
