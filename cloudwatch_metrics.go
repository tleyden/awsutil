package awsutil

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type CloudwatchMetricsUtil struct {
	ec2Api ec2iface.EC2API
	cwApi  cloudwatchiface.CloudWatchAPI
}

func NewCloudwatchMetricsUtilFromRegion(region string) (*CloudwatchMetricsUtil, error) {

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	cwApi := newCloudwatchAPI(session, region)
	ec2Api := newEC2API(session, region)

	return NewCloudwatchMetricsUtil(ec2Api, cwApi)

}

func NewCloudwatchMetricsUtil(ec2Api ec2iface.EC2API, cwApi cloudwatchiface.CloudWatchAPI) (*CloudwatchMetricsUtil, error) {

	cwmUtil := &CloudwatchMetricsUtil{
		cwApi:  cwApi,
		ec2Api: ec2Api,
	}

	return cwmUtil, nil

}

func (cwmu CloudwatchMetricsUtil) GetEc2InstancesWithStates(states []string) ([]*ec2.Instance, error) {

	matchingInstances := []*ec2.Instance{}

	statesPointers := []*string {}
	for _, state := range states {
		statesPointers = append(statesPointers, aws.String(state))
	}

	diInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: statesPointers,
			},
		},
	}

	diOutput, err := cwmu.ec2Api.DescribeInstances(diInput)
	if err != nil {
		return matchingInstances, err
	}

	if diOutput.NextToken != nil {
		log.Printf("Warning: got non empty NextToken, but ignoring since cannot handle yet: %v.  Some results skipped.", diOutput.NextToken)
	}

	for _, r := range diOutput.Reservations {
		for _, i := range r.Instances {
			matchingInstances = append(matchingInstances, i)
		}
	}

	return matchingInstances, nil
}


// FetchRunningEc2InstanceMetrics fetches certain metrics for running EC2 instances
//
// Steps
//
// - For each of the running ec2 instances
//   - Collect metrics with time range based on Frequency param and --dimensions Name=InstanceId,Value=<cur_instance>
// - Return struct result
//
func (cwmu CloudwatchMetricsUtil) FetchRunningEc2InstanceMetrics(input FetchRunningEc2InstanceMetricsInput) ([]FetchRunningEc2InstanceMetricsOuput, error) {

	metricsOutputs := []FetchRunningEc2InstanceMetricsOuput{}

	ec2Instances, err := cwmu.GetEc2InstancesWithStates([]string{EC2_INSTANCE_STATE_RUNNING})
	if err != nil {
		return metricsOutputs, err
	}

	for _, ec2Instance := range ec2Instances {

		now := time.Now()
		oneHourAgo := now.Add(time.Hour * -1)
		metricStatsParams := &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(input.Namespace),
			MetricName: aws.String(input.MetricName),
			Dimensions: []*cloudwatch.Dimension{
				&cloudwatch.Dimension{
					Name: aws.String("InstanceId"),
					Value: ec2Instance.InstanceId,
				},
			},
			Period:     aws.Int64(3600),

			StartTime:  &oneHourAgo,
			EndTime:    &now,
			Statistics: []*string{aws.String("Average")},
		}

		metricStats, err := cwmu.cwApi.GetMetricStatistics(metricStatsParams)
		if err != nil {
			return metricsOutputs, err
		}

		metricOutput := FetchRunningEc2InstanceMetricsOuput{
			Namespace: input.Namespace,
			Ec2InstanceId: *ec2Instance.InstanceId,
			MetricName: input.MetricName,
			Metrics: metricStats,
		}

		metricsOutputs = append(metricsOutputs, metricOutput)


	}

	return metricsOutputs, nil


}
