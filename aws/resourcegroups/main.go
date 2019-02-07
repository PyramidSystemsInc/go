package resourcegroups

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/resourcegroups"
  "github.com/PyramidSystemsInc/go/aws/dynamodb"
  "github.com/PyramidSystemsInc/go/aws/ecs"
  "github.com/PyramidSystemsInc/go/aws/elbv2"
  "github.com/PyramidSystemsInc/go/aws/lambda"
  "github.com/PyramidSystemsInc/go/aws/s3"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/logger"
  "github.com/PyramidSystemsInc/go/str"
)

func Create(groupName string, tagKey string, tagValue string, awsSession *session.Session) {
  resourceGroupsClient := resourcegroups.New(awsSession)
  _, err := resourceGroupsClient.CreateGroup(&resourcegroups.CreateGroupInput{
    Name: aws.String(groupName),
    ResourceQuery: &resourcegroups.ResourceQuery{
      Query: aws.String(str.Concat("{\"ResourceTypeFilters\":[\"AWS::AllSupported\"],\"TagFilters\":[{\"Key\":\"", tagKey, "\", \"Values\":[\"", tagValue, "\"]}]}")),
      Type: aws.String("TAG_FILTERS_1_0"),
    },
  })
  errors.QuitIfError(err)
}

func DeleteAllResources(groupName string, awsSession *session.Session) {
  resourceGroupsClient := resourcegroups.New(awsSession)
  resourcesReport, err := resourceGroupsClient.ListGroupResources(&resourcegroups.ListGroupResourcesInput{
    GroupName: aws.String(groupName),
  })
  errors.QuitIfError(err)
  groupResources := resourcesReport.ResourceIdentifiers
  for _, resource := range groupResources {
    deleteResource(resource, awsSession)
  }
  DeleteGroup(groupName, awsSession)
}

func DeleteGroup(groupName string, awsSession *session.Session) {
  resourceGroupsClient := resourcegroups.New(awsSession)
  _, err := resourceGroupsClient.DeleteGroup(&resourcegroups.DeleteGroupInput{
    GroupName: aws.String(groupName),
  })
  errors.QuitIfError(err)
}

func deleteResource(resource *resourcegroups.ResourceIdentifier, awsSession *session.Session) {
  arn := *resource.ResourceArn
  switch(*resource.ResourceType) {
    case "AWS::DynamoDB::Table":
      dynamodb.DeleteTable(arn, awsSession)
    case "AWS::ECS::Cluster":
      ecs.StopAllTasksInCluster(arn, awsSession)
      ecs.DeleteCluster(arn, awsSession)
    case "AWS::ECS::TaskDefinition":
      ecs.DeregisterTaskDefinition(arn, awsSession)
    case "AWS::ElasticLoadBalancingV2::LoadBalancer":
      elbv2.Delete(arn, awsSession)
    case "AWS::Lambda::Function":
      lambda.Delete(arn, awsSession)
    case "AWS::S3::Bucket":
      s3.EmptyBucket(arn, awsSession)
      s3.DeleteBucket(arn, awsSession)
    default:
      logger.Err(str.Concat("There is a resource of type ", *resource.ResourceType, " which the github.com/PyramidSystemsInc/go/aws/resourcegroups package does not know how to handle"))
  }
}
