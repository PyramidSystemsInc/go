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
// and returns the encryption key id.
func CreateEncryptionKey(awsSession *session.Session) (key string) {
	kmsClient := kms.New(awsSession)

	result, err := kmsClient.CreateKey(&kms.CreateKeyInput{
		Tags: []*kms.Tag{},
	})

	if err != nil {
		fmt.Println("Got error creating key: ", err)
		os.Exit(1)
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
