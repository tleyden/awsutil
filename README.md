[![CircleCI](https://circleci.com/gh/tleyden/awsutil.svg?style=svg)](https://circleci.com/gh/tleyden/awsutil) [![GoDoc](https://godoc.org/github.com/tleyden/awsutil?status.png)](https://godoc.org/github.com/tleyden/awsutil) [![Coverage Status](https://coveralls.io/repos/github/tleyden/awsutil/badge.svg?branch=master)](https://coveralls.io/github/tleyden/awsutil?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/tleyden/awsutil)](https://goreportcard.com/report/github.com/tleyden/awsutil)

A collection of utilities (ok, one utility!) to automate tasks on AWS that aren't covered by the AWS CLI or SDK.

Installation:

```
$ go get -v -t github.com/tleyden/awsutil/...
```

Run:

```
$ awsutil cloudformation stop-instances --region us-east-1 --stackname yourstack
```

This will stop all EC2 instances in the `yourstack` Cloudformation Stack.
