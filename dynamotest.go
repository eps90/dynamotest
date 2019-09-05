package dynamotest

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

type DynamoTester struct {
	dynamoDbSvc       *dynamodb.DynamoDB
	Migrator          *Migrator
	FixturesLoader    DefinitionsLoader
	FixturesDecoder   FixturesDecoder
	TableNameResolver TableNameResolver
	Cleaner           TableCleaner
}

func NewDefaultDynamoTester(dynamoSvc *dynamodb.DynamoDB, migrationsPath string, fixturesPath string) *DynamoTester {
	dynamoTester := DynamoTester{
		dynamoDbSvc:       dynamoSvc,
		Migrator:          NewDefaultMigrator(dynamoSvc, migrationsPath),
		FixturesLoader:    NewJSONFilesystemReader(fixturesPath),
		FixturesDecoder:   NewJSONFixturesDecoder(),
		TableNameResolver: NewMemoizedTableNameResolver(NewTimestampTableNameResolver(new(RealClock))),
		Cleaner:           NewWholeTableDynamoCleaner(dynamoSvc),
	}
	dynamoTester.Migrator.TableNameResolver = dynamoTester.TableNameResolver

	return &dynamoTester
}

func (t *DynamoTester) LoadFixtures(names ...string) error {
	contents, err := t.FixturesLoader.ReadDefinitions(names...)
	if err != nil {
		return errors.Wrap(err, "fixtures: cannot load fixture files")
	}

	writeRequests, err := t.FixturesDecoder.Decode(contents)
	if err != nil {
		return errors.Wrap(err, "fixtures: cannot parse fixture files")
	}

	finalWriteRequests := make(TableWriteRequests)
	for tableName, w := range writeRequests {
		migrateErr := t.Migrator.MigrateTables(tableName)
		if migrateErr != nil {
			return errors.Wrap(migrateErr, "fixtures: cannot migrate tables")
		}

		resolvedName := t.TableNameResolver.Resolve(tableName)
		finalWriteRequests[resolvedName] = w
		migrateErr = t.Cleaner.CleanTable(resolvedName)
		if migrateErr != nil {
			return errors.Wrap(migrateErr, "fixtures: cannot clean table")
		}
	}

	batchWriteInput := dynamodb.BatchWriteItemInput{
		RequestItems: finalWriteRequests,
	}

	_, err = t.dynamoDbSvc.BatchWriteItem(&batchWriteInput)
	if err != nil {
		return errors.Wrap(err, "fixtures: cannot write items")
	}

	return nil
}

func (t *DynamoTester) MustLoadFixtures(names ...string) {
	err := t.LoadFixtures(names...)
	if err != nil {
		panic("Cannot load fixtures: " + err.Error())
	}
}

func (t *DynamoTester) TableNameFor(tableName string) string {
	return t.TableNameResolver.Resolve(tableName)
}
