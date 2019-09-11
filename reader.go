package dynamotest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// DefinitionsLoader defines a struct able to read contents of tables definitions
type DefinitionsLoader interface {
	ReadDefinitions(names ...string) ([][]byte, error)
}

// FilesystemDirectoryLoader reads files from given directory filtering files by given extension
type FilesystemDirectoryLoader struct {
	dir       string
	extension string
}

func NewFilesystemDirectoryLoader(dir string, extension string) *FilesystemDirectoryLoader {
	return &FilesystemDirectoryLoader{dir: dir, extension: extension}
}

// NewJSONFilesystemReader creates FilesystemDirectoryLoader instance which reads JSON files
func NewJSONFilesystemReader(dir string) *FilesystemDirectoryLoader {
	return &FilesystemDirectoryLoader{dir: dir, extension: "json"}
}

func (r *FilesystemDirectoryLoader) ReadDefinitions(names ...string) ([][]byte, error) {
	var files []string
	if len(names) == 0 {
		var err error
		files, err = listFilesInDir(r.dir, r.extension)
		if err != nil {
			return nil, errors.Wrap(err, "migrate: cannot read definitions")
		}
	} else {
		files = combineNamesWithDirectory(names, r.dir, r.extension)
	}

	var result [][]byte
	for _, fileName := range files {
		contents, err := ioutil.ReadFile(filepath.Clean(fileName))
		if err != nil {
			return nil, errors.Wrapf(err, "migrate: cannot read file: '%s'", fileName)
		}
		result = append(result, contents)
	}

	return result, nil
}

func listFilesInDir(directory, extension string) ([]string, error) {
	extPattern := fmt.Sprintf("*.%s", extension)
	fullPath := filepath.Join(directory, extPattern)

	files, err := filepath.Glob(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "migrate: cannot load files from migrations path: '%s'", directory)
	}

	return files, nil
}

func combineNamesWithDirectory(names []string, directory, extension string) []string {
	var result []string
	for _, n := range names {
		fileName := fmt.Sprintf("%s.%s", n, extension)
		result = append(result, filepath.Join(directory, fileName))
	}

	return result
}
