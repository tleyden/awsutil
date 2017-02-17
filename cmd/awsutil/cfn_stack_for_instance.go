package main

import (
	"flag"
	"log"
	"strings"

	"github.com/tleyden/awsutil"
)

type cfnStackForInstanceCmd struct {
	args []string
}

func (c *cfnStackForInstanceCmd) Run(args []string) int {
	// Parse args + get the name of the stack
	cmdFlags := flag.NewFlagSet("cloudformation", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Help() }
	instanceId := cmdFlags.String("instance-id", "", "")
	region := cmdFlags.String("region", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if *region == "" {
		log.Fatalf("You must pass an AWS region")
		return 1
	}

	if *instanceId == "" {
		log.Fatalf("You must pass an EC2 instance id")
		return 1
	}

	// create a session
	cfnUtil, err := awsutil.NewCloudformationUtilFromRegion(*region)
	if err != nil {
		log.Fatalf("Error creating cnf: %v", err)
		return 1
	}

	// Start/restart all instances in stack
	in, stackId, stackName, err := cfnUtil.InCloudformation(*instanceId)
	if err != nil {
		log.Fatalf("Error seeing if instance id in cfn stack: %v.  Err: %v.  Err type: %T", *instanceId, err, err)
		return 1
	}

	if in {
		log.Printf("Instance %s in cloudformation stack: %v (%v)", *instanceId, stackName, stackId)
	} else {
		log.Printf("Instance %s not in any cloudformation stacks.", *instanceId)
	}

	return 0

}

func (c *cfnStackForInstanceCmd) Synopsis() string {
	return "Find the cloudformation stack associated with this instance, if any"
}

func (c *cfnStackForInstanceCmd) Help() string {

	helpText := `
Usage: cloudformation stack-for-instance [options]

  Find the cloudformation stack associated with this instance, if any

Options:

  -instance-id=instance-id  The EC2 instance id
  -region=region            The AWS region

 `
	return strings.TrimSpace(helpText)

}
