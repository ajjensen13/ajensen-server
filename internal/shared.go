package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type FileData struct {
	Name string
	Data []byte
}

func LoadFileData(l *log.Logger, dir string) ([]FileData, error) {
	l.Printf("loading directory: %s", dir)

	fis, err := ioutil.ReadDir(dir)
	switch {
	case os.IsNotExist(err):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("internal: error while reading directory: %w", err)
	}

	var result []FileData
	for _, fi := range fis {
		f := filepath.Join(dir, fi.Name())

		if fi.IsDir() {
			ps, err := LoadFileData(l, f)
			if err != nil {
				return nil, err
			}
			result = append(result, ps...)
			continue
		}

		d, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("internal: error while reading file %q: %w", f, err)
		}
		result = append(result, FileData{Name: f, Data: d})
		l.Printf("loaded file: %s (%d bytes)", f, len(f))
	}

	return result, nil
}

func ParseData(l *log.Logger, ds []FileData, i interface{}) ([]interface{}, error) {
	var result []interface{}

	t := reflect.TypeOf(i).Elem()
	for _, d := range ds {
		into := reflect.New(t).Interface()
		l.Printf("parsing file %q into type %T (%d bytes)", d.Name, into, len(d.Data))

		err := yaml.Unmarshal(d.Data, into)
		if err != nil {
			return nil, fmt.Errorf("internal: error while parsing file %q into type %T: %w", d.Name, into, err)
		}

		result = append(result, into)
	}

	return result, nil
}
