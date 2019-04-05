package dynamodb

import (
	"log"
	"testing"
	"time"

	pacaws "github.com/PyramidSystemsInc/go/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestCreateTable(t *testing.T) {
	session := pacaws.CreateAwsSession("us-east-2")

	svc := dynamodb.New(session)
	name := "test-dynamodb-table"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("LockID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("LockID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(name),
	}

	_, err := svc.CreateTable(input)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	// waiting 10 seconds for table to be created otherwise we'll get an error when we attempt to destroy it

	time.Sleep(10 * time.Second)
	
	describe := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}

	_, err = svc.DescribeTable(describe)
	if err != nil {
		log.Fatal(err)
	}

	delete := &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	}

	_, err = svc.DeleteTable(delete)
	if err != nil {
		log.Fatal(err)
	}
}
