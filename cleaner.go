package dynamotest

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

// TableCleaner defines an interfaces for cleaning contents of a table
type TableCleaner interface {
	CleanTable(tableName string) error
}

// WholeTableDynamoCleaner removes the table.
// While removing it ignores the fact that the table doesn't exist.
// TODO: Write some tests
type WholeTableDynamoCleaner struct {
	dynamoSvc *dynamodb.DynamoDB
}

// NewWholeTableDynamoCleaner creates new instance of WholeTableDynamoCleaner
func NewWholeTableDynamoCleaner(dynamoSvc *dynamodb.DynamoDB) *WholeTableDynamoCleaner {
	return &WholeTableDynamoCleaner{dynamoSvc: dynamoSvc}
}

func (c *WholeTableDynamoCleaner) CleanTable(tableName string) error {
	deleteOutput, err := c.dynamoSvc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	var tableDeleted = true
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() != dynamodb.ErrCodeResourceNotFoundException || !ok {
			return errors.Wrapf(err, "migrate: cannot delete table '%s'", tableName)
		} else {
			tableDeleted = false
		}
	}

	if tableDeleted {
		createInput := createInputFromDeleteInput(deleteOutput)
		_, err = c.dynamoSvc.CreateTable(createInput)
		if err != nil {
			return errors.Wrapf(err, "migrate: cannot recreate table '%s'", tableName)
		}
	}

	return nil
}

func createInputFromDeleteInput(d *dynamodb.DeleteTableOutput) *dynamodb.CreateTableInput {
	createInput := dynamodb.CreateTableInput{
		TableName:            d.TableDescription.TableName,
		AttributeDefinitions: d.TableDescription.AttributeDefinitions,
		KeySchema:            d.TableDescription.KeySchema,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  d.TableDescription.ProvisionedThroughput.ReadCapacityUnits,
			WriteCapacityUnits: d.TableDescription.ProvisionedThroughput.WriteCapacityUnits,
		},
		GlobalSecondaryIndexes: createGSIFromDeletedGSI(d.TableDescription.GlobalSecondaryIndexes),
		LocalSecondaryIndexes:  createLSIFromDeletedLSI(d.TableDescription.LocalSecondaryIndexes),
	}

	return &createInput
}

func createGSIFromDeletedGSI(gsis []*dynamodb.GlobalSecondaryIndexDescription) []*dynamodb.GlobalSecondaryIndex {
	var result []*dynamodb.GlobalSecondaryIndex
	for _, g := range gsis {
		var gsi dynamodb.GlobalSecondaryIndex
		gsi.ProvisionedThroughput.ReadCapacityUnits = g.ProvisionedThroughput.ReadCapacityUnits
		gsi.ProvisionedThroughput.WriteCapacityUnits = g.ProvisionedThroughput.WriteCapacityUnits
		gsi.Projection = g.Projection
		gsi.KeySchema = g.KeySchema
		gsi.IndexName = g.IndexName
		result = append(result, &gsi)
	}

	return result
}

func createLSIFromDeletedLSI(lsis []*dynamodb.LocalSecondaryIndexDescription) []*dynamodb.LocalSecondaryIndex {
	var result []*dynamodb.LocalSecondaryIndex
	for _, l := range lsis {
		var lsi dynamodb.LocalSecondaryIndex
		lsi.IndexName = l.IndexName
		lsi.KeySchema = l.KeySchema
		lsi.Projection = l.Projection
		result = append(result, &lsi)
	}

	return result
}
