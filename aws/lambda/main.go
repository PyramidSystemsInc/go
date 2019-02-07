package lambda

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/lambda"
  "github.com/PyramidSystemsInc/go/errors"
)

func Delete(functionArnOrName string, awsSession *session.Session) {
  lambdaClient := lambda.New(awsSession)
  _, err := lambdaClient.DeleteFunction(&lambda.DeleteFunctionInput{
    FunctionName: aws.String(functionArnOrName),
  })
  errors.QuitIfError(err)
}
