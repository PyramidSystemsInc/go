package ec2

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

// GetAllVpcCidrBlocks - Returns all CIDR blocks in use by VPCs
func GetAllVpcCidrBlocks(awsSession *session.Session) []string {
  ec2Client := ec2.New(awsSession)
  result, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{})
  errors.LogIfError(err)
  if len(result.Vpcs) == 0 {
    errors.LogAndQuit("ERROR: VPC information was queried, but no VPCs were found")
  }
  var cidrBlocks []string
  for _, vpc := range result.Vpcs {
    cidrBlocks = append(cidrBlocks, *vpc.CidrBlock)
  }
  return cidrBlocks
}

// FindPublicIpOfNetworkInterface - Given a network interface ID, returns the public IP associated with it
func FindPublicIpOfNetworkInterface(networkInterfaceId string, awsSession *session.Session) string {
  ec2Client := ec2.New(awsSession)
  result, err := ec2Client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
    NetworkInterfaceIds: []*string{
      aws.String(networkInterfaceId),
    },
  })
  errors.LogIfError(err)
  if len(result.NetworkInterfaces) == 0 {
    notFoundError := str.Concat("A network interface with ID ", networkInterfaceId, " was not found")
    errors.LogAndQuit(notFoundError)
  }
  return *result.NetworkInterfaces[0].Association.PublicIp
}

// ListAllSubnetIds - Given a VPC ID, returns all the IDs of the subnets within it
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

// GetSecurityGroupId - Given the name of a security group, returns the ID of that security group
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
