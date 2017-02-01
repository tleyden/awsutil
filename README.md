A collection of utilities (ok, one utility!) to automate tasks on AWS that aren't covered by the AWS CLI or SDK.

Installation:

```
$ go get -v -t github.com/tleyden/awsutil/...
```

Run:

```
$ awsutil cloudformation stop-instances --region us-east-1 --stackname yourstack
```

Will stop all EC2 instances in the `yourstack` Cloudformation Stack.