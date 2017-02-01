package main

import (
	"log"
	"os"
	"strings"

	"flag"

	"github.com/mitchellh/cli"
	"github.com/tleyden/awsutil"
)

func main() {

	c := cli.NewCLI("awsutil", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"cloudformation stop-instances": func() (cli.Command, error) {
			return &cloudformationStopInstancesCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

type cloudformationStopInstancesCommand struct {
	args []string
}

func (c *cloudformationStopInstancesCommand) Run(args []string) int {

	// Parse args + get the name of the stack
	cmdFlags := flag.NewFlagSet("cloudformation", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Help() }
	stackname := cmdFlags.String("stackname", "", "")
	region := cmdFlags.String("region", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if *region == "" {
		log.Fatalf("You must pass an AWS region")
		return 1
	}

	if *stackname == "" {
		log.Fatalf("You must pass a CloudFormation stack name")
		return 1
	}

	// create a session
	cfnUtil, err := awsutil.NewCloudformationUtilFromRegion(*region)
	if err != nil {
		log.Fatalf("Error creating cnf: %v", err)
		return 1
	}

	// Stop all instances in stack
	err = cfnUtil.StopEC2Instances(*stackname)
	if err != nil {
		log.Fatalf("Error stopping instances for stack: %v.  Err: %v", *stackname, err)
		return 1
	}

	log.Printf("Stopped all instances in stack: %v", *stackname)
	return 0

}

func (c *cloudformationStopInstancesCommand) Synopsis() string {
	return "Stop all instances in given cloudformation stack"
}

func (c *cloudformationStopInstancesCommand) Help() string {

	helpText := `
Usage: cloudformation stop-instances [options]

  Stops all instances in the given cloudformation stack

Options:

  -stackname=stackname      The name of the cloudformation stack
  -region=region            The AWS region

 `
	return strings.TrimSpace(helpText)

}
