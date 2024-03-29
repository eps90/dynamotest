package dynamotest_test

import (
	"testing"

	"github.com/eps90/dynamotest"
	"github.com/stretchr/testify/require"
)

func TestJsonFilesystemLoader(t *testing.T) {
	loader := dynamotest.NewJSONFilesystemReader("test_resources/")
	expectedResult := [][]byte{
		[]byte(`{
  "Name": "This is a Test A file"
}
`),
		[]byte(`{
  "Name": "This is a Test B file"
}
`),
	}
	actualResult, err := loader.ReadDefinitions()
	require.NoError(t, err)
	require.Equal(t, expectedResult, actualResult)
}

func TestJsonFilesystemLoaderWithNames(t *testing.T) {
	loader := dynamotest.NewJSONFilesystemReader("test_resources/")
	names := []string{"a", "nested/c"}
	expectedResult := [][]byte{
		[]byte(`{
  "Name": "This is a Test A file"
}
`),
		[]byte(`{
  "Name": "This is a Test nested/C file"
}
`),
	}
	actualResult, err := loader.ReadDefinitions(names...)
	require.NoError(t, err)
	require.Equal(t, expectedResult, actualResult)
}
