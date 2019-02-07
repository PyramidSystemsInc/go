package dynamodb

import (
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/PyramidSystemsInc/go/aws/util"
  "github.com/PyramidSystemsInc/go/errors"
)

func DeleteTable(arnOrName string, awsSession *session.Session) {
  tableName := getTableName(arnOrName)
  dynamoDbClient := dynamodb.New(awsSession)
  _, err := dynamoDbClient.DeleteTable(&dynamodb.DeleteTableInput{
    TableName: aws.String(tableName),
  })
  errors.QuitIfError(err)
}

func getTableName(arnOrName string) string {
  if util.IsArn(arnOrName) {
    return arnOrName[strings.LastIndex(arnOrName, "/") + 1:len(arnOrName)]
  } else {
    return arnOrName
  }
}
