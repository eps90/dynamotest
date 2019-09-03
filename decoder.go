package dynamotest

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// MigrationDecoder defines an interface for unmarshalling raw migrations into DynamoDB's CreateTableInput
type MigrationDecoder interface {
	Decode(input []byte) (*dynamodb.CreateTableInput, error)
}

// JsonMigrationDecoder decodes JSON migrations into DynamoDB's CreateTableInput
type JsonMigrationDecoder struct {
}

func (*JsonMigrationDecoder) Decode(input []byte) (*dynamodb.CreateTableInput, error) {
	var createTable dynamodb.CreateTableInput
	err := json.Unmarshal(input, &createTable)
	if err != nil {
		return nil, errors.Wrap(err, "migrate: cannot parse migration file")
	}

	return &createTable, nil
}

// TableWriteRequests is a collection of dynamodb.WriteRequest grouped by table
type TableWriteRequests map[string][]*dynamodb.WriteRequest

type fixture struct {
	TableName string                   `json:"table"`
	Items     []map[string]interface{} `json:"items"`
}

// FixturesDecoder defines an interface of collection of fixtures contents which writes to TableWriteRequests
type FixturesDecoder interface {
	Decode(input [][]byte) (TableWriteRequests, error)
}

type JsonFixturesDecoder struct {
}

func NewJsonFixturesDecoder() *JsonFixturesDecoder {
	return &JsonFixturesDecoder{}
}

func (*JsonFixturesDecoder) Decode(input [][]byte) (TableWriteRequests, error) {
	writeRequests := make(TableWriteRequests)
	for _, fixtureContents := range input {
		var fx fixture
		err := json.Unmarshal(fixtureContents, &fx)
		if err != nil {
			return nil, errors.Wrap(err, "fixtures: cannot parse fixture")
		}

		for _, fixtureItems := range fx.Items {
			m, err := dynamodbattribute.MarshalMap(fixtureItems)
			if err != nil {
				return nil, errors.Wrap(err, "fixtures: cannot marshal map")
			}
			writeRequest := &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: m,
				},
			}
			writeRequests[fx.TableName] = append(writeRequests[fx.TableName], writeRequest)
		}
	}

	return writeRequests, nil
}
