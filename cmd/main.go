package main

import (
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"flag"
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
	args              []string
}

func (c *cloudformationStopInstancesCommand) Run(args []string) int {

	// Parse args + get the name of the stack
	cmdFlags := flag.NewFlagSet("cloudformation", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Help() }
	stackname := cmdFlags.String("stackname", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Stop all instances in stack
	err := awsutil.StopInstancesInStack(*stackname)
	if err != nil {
		log.Printf("Error stopping instances for stack: %v.  Err: %v", *stackname, err)
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

  Starts the Consul agent and runs until an interrupt is received. The
  agent represents a single node in a cluster.

Options:

  -stack=stackname          The name of the cloudformation stack

 `
	return strings.TrimSpace(helpText)

}
