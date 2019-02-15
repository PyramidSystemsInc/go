package elbv2

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/elbv2"
  "github.com/PyramidSystemsInc/go/aws/ec2"
  "github.com/PyramidSystemsInc/go/errors"
)

func Create(name string, awsSession *session.Session) (string, string, string) {
  elbv2Client := elbv2.New(awsSession)
  vpcId := "vpc-76cf681f"
  loadBalancer, err := elbv2Client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
    Name: aws.String(name),
    Subnets: ec2.ListAllSubnetIds(vpcId, awsSession),
  })
  errors.QuitIfError(err)
  loadBalancerArn := loadBalancer.LoadBalancers[0].LoadBalancerArn
  loadBalancerUrl := loadBalancer.LoadBalancers[0].DNSName
  listenerArn := createDefaultListener(loadBalancerArn, elbv2Client)
  return *loadBalancerArn, *listenerArn, *loadBalancerUrl
}

func Delete(arn string, awsSession *session.Session) {
  elbv2Client := elbv2.New(awsSession)
  _, err := elbv2Client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
    LoadBalancerArn: aws.String(arn),
  })
  errors.QuitIfError(err)
}

func Exists(nameOrArn string, awsSession *session.Session) bool {
  loadBalancer := getLoadBalancer(nameOrArn, awsSession)
  return loadBalancer != nil
}

func Tag(nameOrArn string, key string, value string, awsSession *session.Session) {
  loadBalancer := getLoadBalancer(nameOrArn, awsSession)
  if loadBalancer != nil {
    arn := getArn(loadBalancer)
    elbv2Client := elbv2.New(awsSession)
    _, err := elbv2Client.AddTags(&elbv2.AddTagsInput{
      ResourceArns: []*string{
        aws.String(arn),
      },
      Tags: []*elbv2.Tag{
        &elbv2.Tag{
          Key: aws.String(key),
          Value: aws.String(value),
        },
      },
    })
    errors.LogIfError(err)
  }
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

func getArn(loadBalancer *elbv2.LoadBalancer) string {
  if loadBalancer != nil {
    return *loadBalancer.LoadBalancerArn
  }
  return ""
}

func getLoadBalancer(nameOrArn string, awsSession *session.Session) *elbv2.LoadBalancer {
  elbv2Client := elbv2.New(awsSession)
  result, err := elbv2Client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
    Names: []*string{
      aws.String(nameOrArn),
    },
  })
  if loadBalancerFound(result, err) {
    return result.LoadBalancers[0]
  } else {
    result, err := elbv2Client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
      LoadBalancerArns: []*string {
        aws.String(nameOrArn),
      },
    })
    if loadBalancerFound(result, err) {
      return result.LoadBalancers[0]
    } else {
      return nil
    }
  }
}

func loadBalancerFound(result *elbv2.DescribeLoadBalancersOutput, err error) bool {
  return err == nil && len(result.LoadBalancers) > 0
}
