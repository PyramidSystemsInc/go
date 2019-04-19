package s3

import (
	"fmt"
	"os"
	"strings"

	"github.com/PyramidSystemsInc/go/aws/util"
	"github.com/PyramidSystemsInc/go/errors"
	"github.com/PyramidSystemsInc/go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// The allowed values for the `access` parameter can be found here: https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl
func MakeBucket(bucketName string, access string, region string, awsSession *session.Session) {
	s3Client := s3.New(awsSession)
	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		ACL:    aws.String(access),
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
		ObjectLockEnabledForBucket: aws.Bool(false),
	})
	errors.QuitIfError(err)
}

func DeleteBucket(bucketNameOrArn string, awsSession *session.Session) {
	bucketName := getBucketName(bucketNameOrArn)
	s3Client := s3.New(awsSession)
	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	errors.QuitIfError(err)
}

func EmptyBucket(bucketNameOrArn string, awsSession *session.Session) {
	bucketName := getBucketName(bucketNameOrArn)
	s3Client := s3.New(awsSession)
	bucketObjects, err := s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	errors.QuitIfError(err)
	bucketContents := bucketObjects.Contents
	objectIdentifiers := make([]*s3.ObjectIdentifier, 0)
	for _, file := range bucketContents {
		objectIdentifiers = append(objectIdentifiers, &s3.ObjectIdentifier{
			Key: file.Key,
		})
	}
	_, err = s3Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: objectIdentifiers,
			Quiet:   aws.Bool(true),
		},
	})
}

func EnableWebsiteHosting(bucketName string, awsSession *session.Session) {
	s3Client := s3.New(awsSession)
	documentName := "index.html"
	_, err := s3Client.PutBucketWebsite(&s3.PutBucketWebsiteInput{
		Bucket: aws.String(bucketName),
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			ErrorDocument: &s3.ErrorDocument{
				Key: aws.String(documentName),
			},
			IndexDocument: &s3.IndexDocument{
				Suffix: aws.String(documentName),
			},
		},
	})
	errors.QuitIfError(err)
}

func TagBucket(bucketName string, key string, value string, awsSession *session.Session) {
	s3Client := s3.New(awsSession)
	_, err := s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				&s3.Tag{
					Key:   aws.String(key),
					Value: aws.String(value),
				},
			},
		},
	})
	errors.QuitIfError(err)
}

func getBucketName(arnOrName string) string {
	if util.IsArn(arnOrName) {
		return arnOrName[strings.LastIndex(arnOrName, ":::")+3 : len(arnOrName)]
	} else {
		return arnOrName
	}
}

// EncryptBucket turns on encryption on the S3 bucket
func EncryptBucket(bucket, key string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)

	defEnc := &s3.ServerSideEncryptionByDefault{KMSMasterKeyID: aws.String(key), SSEAlgorithm: aws.String(s3.ServerSideEncryptionAwsKms)}
	rule := &s3.ServerSideEncryptionRule{ApplyServerSideEncryptionByDefault: defEnc}
	rules := []*s3.ServerSideEncryptionRule{rule}
	serverConfig := &s3.ServerSideEncryptionConfiguration{Rules: rules}
	input := &s3.PutBucketEncryptionInput{Bucket: aws.String(bucket), ServerSideEncryptionConfiguration: serverConfig}
	_, err := svc.PutBucketEncryption(input)
	if err != nil {
		fmt.Println("Got an error adding default KMS encryption to bucket", bucket)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger.Info("Bucket " + bucket + " now has KMS encryption by default")
}

// EnableVersioning turnson version on the S3 bucket
func EnableVersioning(bucket string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)
	input := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &s3.VersioningConfiguration{
			MFADelete: aws.String("Disabled"),
			Status:    aws.String("Enabled"),
		},
	}

	_, err := svc.PutBucketVersioning(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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
}
