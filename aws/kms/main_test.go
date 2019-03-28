package kms

import (
	"fmt"
	"testing"

	pacaws "github.com/PyramidSystemsInc/go/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
)

// TestCreateEncryptionKey creates an encryption key then tests the successful creation of an encryption key by
// calling the DescribeKey function using the created keys id. If an error is returned, the key wasn't created.
func TestCreateEncryptionKey(t *testing.T) {
	session := pacaws.CreateAwsSession("us-east-2")

	key := CreateEncryptionKey(session, "createdBy", "TestCreateEncryptionKey unit test")

	svc := kms.New(session)
	input := &kms.DescribeKeyInput{
		KeyId: aws.String(key),
	}

	_, err := svc.DescribeKey(input)

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
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	ScheduleEncryptionKeyDeletion(key, session)
}

// TestScheduleEncryptionKeyDeletion creates an encryption key and then attempts to schedule it for deletion.
func TestScheduleEncryptionKeyDeletion(t *testing.T) {
	session := pacaws.CreateAwsSession("us-east-2")
	key := CreateEncryptionKey(session, "createdBy", "TestScheduleEncryptionKeyDeletion unit test")
	ScheduleEncryptionKeyDeletion(key, session)
}
