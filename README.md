[![Build Status](http://drone.couchbase.io/api/badges/tleyden/awsutil/status.svg)](http://drone.couchbase.io/tleyden/awsutil)

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
