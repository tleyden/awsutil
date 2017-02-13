package awsutil

import (
	"testing"
	"log"
	"time"
)

func TestFetchRunningEc2InstanceMetrics(t *testing.T) {

	cwUtil, err := NewCloudwatchMetricsUtilFromRegion("us-east-1")
	if err != nil {
		t.Fatalf("Failed to create cwUtil: %v", err)
	}

	input :=  FetchRunningEc2InstanceMetricsInput {
		RecentTimeSpan: time.Duration(time.Second * 3600),
		Namespace: "AWS/EC2",
		MetricName: "CPUUtilization",
	}


	metrics, err := cwUtil.FetchRunningEc2InstanceMetrics(input)
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	log.Printf("got metrics: %+v", metrics)


}
