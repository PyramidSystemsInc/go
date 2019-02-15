package ec2

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/PyramidSystemsInc/go/errors"
)

func FindPublicIpOfNetworkInterface(networkInterfaceId string, awsSession *session.Session) string {
  ec2Client := ec2.New(awsSession)
  result, err := ec2Client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
    NetworkInterfaceIds: []*string{
      aws.String(networkInterfaceId),
    },
  })
  errors.LogIfError(err)
  return *result.NetworkInterfaces[0].Association.PublicIp
}

func ListAllSubnetIds(vpcId string, awsSession *session.Session) []*string {
  ec2Client := ec2.New(awsSession)
  result, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{})
  errors.LogIfError(err)
  subnets := make([]*string, 0)
  for _, subnet := range result.Subnets {
    if *subnet.VpcId == vpcId {
      subnets = append(subnets, subnet.SubnetId)
    }
  }
  return subnets
}

func GetSecurityGroupId(securityGroupName string, awsSession *session.Session) *string {
  ec2Client := ec2.New(awsSession)
  result, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
    GroupNames: []*string{
      aws.String(securityGroupName),
    },
  })
  errors.LogIfError(err)
  if len(result.SecurityGroups) == 1 {
    return result.SecurityGroups[0].GroupId
  } else {
    notFound := ""
    return &notFound
  }
}
