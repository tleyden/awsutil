package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	c := cli.NewCLI("awsutil", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"cloudformation stop-instances": func() (cli.Command, error) {
			return &cloudformationStopInstancesCommand{}, nil
		},
		"cloudformation start-instances": func() (cli.Command, error) {
			return &cloudformationStartInstancesCommand{}, nil
		},
		"cloudformation stack-for-instance": func() (cli.Command, error) {
			return &cfnStackForInstanceCmd{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
