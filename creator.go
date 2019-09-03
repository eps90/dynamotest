package dynamotest

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

// TableCreator defines an interface for structs able to create new DynamoDB table
type TableCreator interface {
	CreateTable(input *dynamodb.CreateTableInput) error
}

// DefaultTableCreator just creates a table with given CreateTableInput
// TODO: Write some tests
type DefaultTableCreator struct {
	dynamoSvc *dynamodb.DynamoDB
}

// NewDefaultTableCreator creates new instance of DefaultTableCreator
func NewDefaultTableCreator(dynamoSvc *dynamodb.DynamoDB) *DefaultTableCreator {
	return &DefaultTableCreator{dynamoSvc: dynamoSvc}
}

func (c *DefaultTableCreator) CreateTable(input *dynamodb.CreateTableInput) error {
	_, err := c.dynamoSvc.CreateTable(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != dynamodb.ErrCodeResourceInUseException || !ok {
			return errors.Wrapf(err, "migrate: cannot create table '%s'", *input.TableName)
		}
	}

	return nil
}
