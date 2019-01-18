package ecs

import (
  "strings"
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ecs"
  "github.com/PyramidSystemsInc/go/aws/ec2"
  "github.com/PyramidSystemsInc/go/aws/ecr"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

type Container struct {
  EnvironmentVars  map[string]string
  Essential        bool
  ImageName        string
  Name             string
}

func LaunchFargateContainer(taskDefinitionName string, clusterName string, securityGroupName string, awsSession *session.Session) string {
  clusterArn := findCluster(clusterName, awsSession)
  if clusterArn == "" {
    createClusterIfDoesNotExist(clusterName, awsSession)
  }
  taskArn := runTask(taskDefinitionName, clusterName, securityGroupName, awsSession)
  publicIp := findPublicIpOfTask(clusterName, taskArn, awsSession)
  return publicIp
}

func RegisterFargateTaskDefinition(taskName string, awsSession *session.Session, containers []Container) string {
  ecsClient := ecs.New(awsSession)
  ecrUrl, err := ecr.GetUrl()
  errors.LogIfError(err)
  // TODO: Remove hardcoded ecsTaskExecutionRole ARN
  // TODO: Add CPU and Memory as parameters
  var containerDefinitions []*ecs.ContainerDefinition
  for _, container := range containers {
    var environmentVariables []*ecs.KeyValuePair
    for name, value := range container.EnvironmentVars {
      environmentVariables = append(environmentVariables, &ecs.KeyValuePair{
        Name: aws.String(name),
        Value: aws.String(value),
      })
    }
    containerDefinitions = append(containerDefinitions, &ecs.ContainerDefinition{
      Environment: environmentVariables,
      Essential: aws.Bool(container.Essential),
      Image: aws.String(str.Concat(ecrUrl, "/", container.ImageName)),
      Name: aws.String(container.Name),
    })
  }
  _, err = ecsClient.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
    ContainerDefinitions: containerDefinitions,
    Cpu: aws.String("2048"),
    ExecutionRoleArn: aws.String("arn:aws:iam::118104210923:role/ecsTaskExecutionRole"),
    Family: aws.String(taskName),
    RequiresCompatibilities: []*string{
      aws.String("FARGATE"),
    },
    Memory: aws.String("16384"),
    NetworkMode: aws.String("awsvpc"),
    TaskRoleArn: aws.String("jenkins_instance"),
  })
  errors.LogIfError(err)
  return ""
}

func createClusterIfDoesNotExist(clusterName string, awsSession *session.Session) {
  ecsClient := ecs.New(awsSession)
  _, err := ecsClient.CreateCluster(&ecs.CreateClusterInput{
    ClusterName: &clusterName,
  })
  errors.LogIfError(err)
}

func findCluster(clusterName string, awsSession *session.Session) string {
  ecsClient := ecs.New(awsSession)
  result, err := ecsClient.ListClusters(&ecs.ListClustersInput{})
  errors.LogIfError(err)
  for _, arn := range result.ClusterArns {
    if strings.HasSuffix(*arn, "/" + clusterName) {
      return *arn
    }
  }
  return ""
}

func findPublicIpOfTask(clusterName string, taskArn string, awsSession *session.Session) string {
  time.Sleep(7 * time.Second)
  networkInterfaceId := findNetworkInterfaceIdOfTask(clusterName, taskArn, awsSession)
  return ec2.FindPublicIpOfNetworkInterface(networkInterfaceId, awsSession)
}

func findNetworkInterfaceIdOfTask(clusterName string, taskArn string, awsSession *session.Session) string {
  ecsClient := ecs.New(awsSession)
  result, err := ecsClient.DescribeTasks(&ecs.DescribeTasksInput{
    Cluster: aws.String(clusterName),
    Tasks: []*string{
      aws.String(taskArn),
    },
  })
  errors.LogIfError(err)
  networkDetails := result.Tasks[0].Attachments[0].Details
  var networkInterfaceId string
  for _, networkDetail := range networkDetails {
    if *networkDetail.Name == "networkInterfaceId" {
      networkInterfaceId = *networkDetail.Value
    }
  }
  return networkInterfaceId
}

func runTask(taskDefinitionName string, clusterName string, securityGroupName string, awsSession *session.Session) string {
  ecsClient := ecs.New(awsSession)
  result, err := ecsClient.RunTask(&ecs.RunTaskInput{
    Cluster: &clusterName,
    LaunchType: aws.String("FARGATE"),
    NetworkConfiguration: &ecs.NetworkConfiguration{
      AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
        AssignPublicIp: aws.String("ENABLED"),
        SecurityGroups: []*string{
          ec2.GetSecurityGroupId(securityGroupName, awsSession),
        },
        Subnets: ec2.ListAllSubnetIds(awsSession),
      },
    },
    TaskDefinition: &taskDefinitionName,
  })
  errors.LogIfError(err)
  return *result.Tasks[0].TaskArn
}
