package kms

import (
  "fmt"
  "os"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/kms"
)

// CreateEncryptionKey creates a customer managed key in the AWS Key Management Service
// and returns the encryption key id. The session holds the region information, the k and v
// are key/value pairs used to tag the encryption key for later identification.
func CreateEncryptionKey(awsSession *session.Session, k string, v string) (key string) {
  kmsClient := kms.New(awsSession)

  result, err := kmsClient.CreateKey(&kms.CreateKeyInput{
    Tags: []*kms.Tag{
      {
        TagKey:   aws.String(k),
        TagValue: aws.String(v),
      },
    },
  })

  if err != nil {
    fmt.Println("error creating encryption key: ", err)
    os.Exit(1)
  }

  alias := "alias/pac/" + v

  _, err = kmsClient.CreateAlias(&kms.CreateAliasInput{
    AliasName:   aws.String(alias),
    TargetKeyId: aws.String(*result.KeyMetadata.KeyId),
  })

  if err != nil {
    fmt.Println("error creating encryption key alias: ", err)
  }

  return *result.KeyMetadata.KeyId
}

// ScheduleEncryptionKeyDeletion schedules encryption key for deletion in 7 days
// AWS does not allow for immediate deletion of encryption keys just-in-case encrypted
// resources are later found and need decrypting.
func ScheduleEncryptionKeyDeletion(key string, awsSession *session.Session) {
  svc := kms.New(awsSession)
  input := &kms.ScheduleKeyDeletionInput{
    KeyId:               aws.String(key),
    PendingWindowInDays: aws.Int64(7),
  }

  result, err := svc.ScheduleKeyDeletion(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case kms.ErrCodeNotFoundException:
        fmt.Println(kms.ErrCodeNotFoundException, aerr.Error())
      case kms.ErrCodeInvalidArnException:
        fmt.Println(kms.ErrCodeInvalidArnException, aerr.Error())
      case kms.ErrCodeDependencyTimeoutException:
        fmt.Println(kms.ErrCodeDependencyTimeoutException, aerr.Error())
      case kms.ErrCodeInternalException:
        fmt.Println(kms.ErrCodeInternalException, aerr.Error())
      case kms.ErrCodeInvalidStateException:
        fmt.Println(kms.ErrCodeInvalidStateException, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    return
  }

  fmt.Println(result)
}

// GetParameter returns the value stored in the systems manager paramter store at the given path
func GetParameter(awsSession *session.Session, k, v, path string) {
  //
}
