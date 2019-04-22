package dynamodb

import (
  "strings"

  "github.com/PyramidSystemsInc/go/aws/util"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

// DeleteTable - Deletes an AWS DynamoDB table
func DeleteTable(arnOrName string, awsSession *session.Session) {
  tableName := getTableName(arnOrName)
  dynamoDbClient := dynamodb.New(awsSession)
  _, err := dynamoDbClient.DeleteTable(&dynamodb.DeleteTableInput{
    TableName: aws.String(tableName),
  })
  errors.QuitIfError(err)
}

// CreateTable - Creates a new AWS DynamoDB table
func CreateTable(input *dynamodb.CreateTableInput, awsSession *session.Session) {
  dynamoDbClient := dynamodb.New(awsSession)
  _, err := dynamoDbClient.CreateTable(input)

  errors.QuitIfError(err)
}

func getTableName(arnOrName string) string {
  if util.IsArn(arnOrName) {
    return arnOrName[strings.LastIndex(arnOrName, "/")+1 : len(arnOrName)]
  }

  return arnOrName
}

