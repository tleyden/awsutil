package main

import (
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
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

}

func (c *cloudformationStopInstancesCommand) Run(args []string) int {
	log.Printf("cloudformation.Run()")
	return 0
}

func (c *cloudformationStopInstancesCommand) Synopsis() string {
	return "Cloudformation utilities"
}

func (c *cloudformationStopInstancesCommand) Help() string {
	helpText := `
Usage: cloudformation [options]

  Starts the Consul agent and runs until an interrupt is received. The
  agent represents a single node in a cluster.

Options:



 `
	return strings.TrimSpace(helpText)
}
