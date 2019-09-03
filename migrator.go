package dynamotest

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

type Migrator struct {
	definitions       map[string]*dynamodb.CreateTableInput
	MigrationsLoader  DefinitionsLoader
	MigrationsDecoder MigrationDecoder
	TableNameResolver TableNameResolver
	Creator           TableCreator
}

func NewDefaultMigrator(dynamoSvc *dynamodb.DynamoDB, migrationsPath string) *Migrator {
	return &Migrator{
		MigrationsLoader:  NewJsonFilesystemReader(migrationsPath),
		MigrationsDecoder: new(JsonMigrationDecoder),
		TableNameResolver: new(DefaultTableNameResolver),
		Creator:           NewDefaultTableCreator(dynamoSvc),
	}
}

func (m *Migrator) MigrateTables(tableNames ...string) error {
	err := m.loadDefinitions()
	if err != nil {
		return errors.Wrap(err, "migrate: cannot load migration definitions")
	}

	for _, tableName := range tableNames {
		if tableDefinition, ok := m.definitions[tableName]; ok {
			err := m.Creator.CreateTable(tableDefinition)
			if err != nil {
				return errors.Wrapf(err, "migrate: cannot crete table %s(%s)", tableName, *tableDefinition.TableName)
			}
		}
	}

	return nil
}

func (m *Migrator) loadDefinitions() error {
	if len(m.definitions) > 0 {
		return nil
	}

	m.definitions = make(map[string]*dynamodb.CreateTableInput)
	tablesDefinitions, err := m.MigrationsLoader.ReadDefinitions()

	if err != nil {
		return errors.Wrap(err, "migrate: cannot load migration files")
	}

	for _, d := range tablesDefinitions {
		createTableInput, err := m.MigrationsDecoder.Decode(d)
		if err != nil {
			return errors.Wrap(err, "migrate: cannot decode migration file")
		}

		tableName := createTableInput.TableName
		newTableName := m.TableNameResolver.Resolve(*createTableInput.TableName)
		createTableInput.TableName = aws.String(newTableName)

		m.definitions[*tableName] = createTableInput
	}

	return nil
}
