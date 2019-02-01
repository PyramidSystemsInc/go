package elbv2

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/elbv2"
  "github.com/PyramidSystemsInc/go/aws/ec2"
  "github.com/PyramidSystemsInc/go/errors"
)

func CreateLoadBalancer(name string, awsSession *session.Session) (string, string, string) {
  elbv2Client := elbv2.New(awsSession)
  loadBalancer, err := elbv2Client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
    Name: aws.String(name),
    Subnets: ec2.ListAllSubnetIds(awsSession),
  })
  errors.QuitIfError(err)
  loadBalancerArn := loadBalancer.LoadBalancers[0].LoadBalancerArn
  loadBalancerUrl := loadBalancer.LoadBalancers[0].DNSName
  listenerArn := createDefaultListener(loadBalancerArn, elbv2Client)
  return *loadBalancerArn, *listenerArn, *loadBalancerUrl
}

func LoadBalancerExists(name string, awsSession *session.Session) bool {
  elbv2Client := elbv2.New(awsSession)
  result, err := elbv2Client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
    Names: []*string{
      aws.String(name),
    },
  })
  return err == nil && len(result.LoadBalancers) > 0
}

func createDefaultListener(loadBalancerArn *string, elbv2Client *elbv2.ELBV2) *string {
  listener, err := elbv2Client.CreateListener(&elbv2.CreateListenerInput{
    DefaultActions: []*elbv2.Action{
      {
        Order: aws.Int64(1),
        RedirectConfig: &elbv2.RedirectActionConfig{
          Host: aws.String("#{host}"),
          Path: aws.String("/api"),
          Port: aws.String("80"),
          Protocol: aws.String("HTTP"),
          Query: aws.String("#{query}"),
          StatusCode: aws.String("HTTP_301"),
        },
        Type: aws.String("redirect"),
      },
    },
    LoadBalancerArn: loadBalancerArn,
    Port: aws.Int64(80),
    Protocol: aws.String("HTTP"),
  })
  errors.QuitIfError(err)
  return listener.Listeners[0].ListenerArn
}
