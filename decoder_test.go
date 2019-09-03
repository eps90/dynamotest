package dynamotest_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/eps90/dynamotest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonMigrationDecoder(t *testing.T) {
	decoder := new(dynamotest.JsonMigrationDecoder)
	input := createSampleMigrationBytes()
	expectedOutput := createSampleCreateTableInput()
	actualOutput, err := decoder.Decode(input)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, actualOutput)
}

func TestJsonMigrationDecoderInvalidInput(t *testing.T) {
	decoder := new(dynamotest.JsonMigrationDecoder)
	input := createSampleInvalidMigrationBytes()
	_, err := decoder.Decode(input)

	require.Error(t, err)
}

func TestJsonFixturesDecoder(t *testing.T) {
	decoder := dynamotest.NewJsonFixturesDecoder()
	input := createSampleFixturesBytes()
	expected := createSampleWriteRequestMap()
	actual, err := decoder.Decode(input)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func createSampleMigrationBytes() []byte {
	return []byte(`
		{
		  "TableName": "tableName",
		  "AttributeDefinitions": [
			{
			  "AttributeName": "ID",
			  "AttributeType": "N"
			}
		  ],
		  "KeySchema": [
			{
			  "AttributeName": "ID",
			  "KeyType": "HASH"
			}
		  ],
		  "ProvisionedThroughput": {
			"ReadCapacityUnits": 5,
			"WriteCapacityUnits": 10
		  }
		}
	`)
}

func createSampleInvalidMigrationBytes() []byte {
	return []byte(`
		{
		  "TableName": "tableName",
		  "AttributeDefinitions": [
			{
			  "AttributeName": "ID",
			  "AttributeType": "N"
			}
		  ],
	`)
}

func createSampleCreateTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		TableName: aws.String("tableName"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			WriteCapacityUnits: aws.Int64(10),
			ReadCapacityUnits:  aws.Int64(5),
		},
	}
}

func createSampleFixturesBytes() [][]byte {
	return [][]byte{
		[]byte(`
			{
				"table": "tableName",
				"items": [
					{
					  "ID": 1,
					  "Name": "Abc"
					},
					{
					  "ID": 2,
					  "Name": "Bca"
					}
				]
			}`),
		[]byte(`
				{
				  "table": "tableName",
				  "items": [
					{
					  "ID": 5,
					  "Nested": {
						"Name": "BBB",
						"Price": 200,
						"CreatedAt": "2019-01-05T12:13:56Z"
					  }
					}
				  ]
				}`),
		[]byte(`
				{
				  "table": "otherTable",
				  "items": [
					{
					  "ID": 7,
					  "Name": "CCC"
					}
				  ]
				}`),
	}
}

func createSampleWriteRequestMap() dynamotest.TableWriteRequests {
	return dynamotest.TableWriteRequests{
		"tableName": {
			{
				PutRequest: &dynamodb.PutRequest{
					Item: marshalMap(map[string]interface{}{
						"ID":   1,
						"Name": "Abc",
					}),
				},
			},
			{
				PutRequest: &dynamodb.PutRequest{
					Item: marshalMap(map[string]interface{}{
						"ID":   2,
						"Name": "Bca",
					}),
				},
			},
			{
				PutRequest: &dynamodb.PutRequest{
					Item: marshalMap(map[string]interface{}{
						"ID": 5,
						"Nested": map[string]interface{}{
							"Name":      "BBB",
							"Price":     200,
							"CreatedAt": "2019-01-05T12:13:56Z",
						},
					}),
				},
			},
		},
		"otherTable": {
			{
				PutRequest: &dynamodb.PutRequest{
					Item: marshalMap(map[string]interface{}{
						"ID":   7,
						"Name": "CCC",
					}),
				},
			},
		},
	}
}

func marshalMap(in map[string]interface{}) map[string]*dynamodb.AttributeValue {
	m, _ := dynamodbattribute.MarshalMap(in)
	return m
}
