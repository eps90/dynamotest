package dynamotest_test

import (
	"github.com/eps90/dynamotest"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDefaultTableNameResolver(t *testing.T) {
	resolver := new(dynamotest.DefaultTableNameResolver)
	inputTableName := "tableName"
	actual := resolver.Resolve(inputTableName)
	if actual != inputTableName {
		t.Errorf("Expected output to be %s, got %s instead", inputTableName, actual)
	}
}

func TestRandomTableNameResolver(t *testing.T) {
	resolver := new(dynamotest.RandomTableNameResolver)
	resolver.Seed = 15

	inputTableName := "tableName"
	expected := "tableName_FOGwh"
	actual := resolver.Resolve(inputTableName)
	if actual != expected {
		t.Errorf("Expected output to be %s, got %s instead", inputTableName, actual)
	}
}

func TestRandomTableNameResolverWithoutSeed(t *testing.T) {
	resolver := new(dynamotest.RandomTableNameResolver)

	inputTableName := "tableName"

	firstTry := resolver.Resolve(inputTableName)
	secondTry := resolver.Resolve(inputTableName)
	if firstTry == secondTry {
		t.Error("Expected tableName to different")
	}
}

func TestMemoizedTableNameResolver(t *testing.T) {
	resolver := dynamotest.NewMemoizedTableNameResolver(dynamotest.NewRandomTableNameResolver())
	inputTableName := "tableName"

	firstTry := resolver.Resolve(inputTableName)
	secondTry := resolver.Resolve(inputTableName)
	if firstTry != secondTry {
		t.Errorf("Expected output to be the same, got %s and %s", firstTry, secondTry)
	}
}

func TestTimestampBasedTableNameResolver(t *testing.T) {
	frozenTime := time.Date(2019, 4, 5, 12, 55, 13, 1, time.UTC)
	clock := dynamotest.FakeClock{FrozenTime: frozenTime}
	resolver := dynamotest.NewTimestampTableNameResolver(&clock)
	inputTableName := "tableName"

	expectedTableName := "tableName_1554468913000000001"
	actualTableName := resolver.Resolve(inputTableName)

	require.Equal(t, expectedTableName, actualTableName)
}
