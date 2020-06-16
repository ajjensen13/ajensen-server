/*
Copyright © 2020 A. Jensen <jensen.aaro@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package internal

import (
	"fmt"
	"github.com/ajjensen13/gke"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type FileData struct {
	Name string
	Data []byte
}

// LoadDirData recursively loads the files from dir into a FileData slice
func LoadDirData(lg gke.Logger, dir string) ([]FileData, error) {
	lg.Default(gke.NewFmtMsgData("LoadDirData: %s", dir))

	fis, err := ioutil.ReadDir(dir)
	switch {
	case os.IsNotExist(err):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("error while reading directory: %w", err)
	}

	var result []FileData
	for _, fi := range fis {
		f := filepath.Join(dir, fi.Name())

		if fi.IsDir() {
			ps, err := LoadDirData(lg, f)
			if err != nil {
				return nil, err
			}
			result = append(result, ps...)
			continue
		}

		d, err := loadFileData(lg, f)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, nil
}

func loadFileData(lg gke.Logger, f string) (FileData, error) {
	d, err := ioutil.ReadFile(f)
	if err != nil {
		return FileData{}, fmt.Errorf("error while reading file %q: %w", f, err)
	}

	lg.Default(gke.NewFmtMsgData("loadFileData: %s (%d bytes)", f, len(f)))
	return FileData{Name: f, Data: d}, nil
}

// ParseFileData parses file data into a new slice by calling yaml.Unmarshal.
// Parameter i is the destination type.
func ParseFileData(lg gke.Logger, ds []FileData, i interface{}) ([]interface{}, error) {
	var result []interface{}

	t := reflect.TypeOf(i).Elem()
	for _, d := range ds {
		into := reflect.New(t).Interface()
		lg.Default(gke.NewFmtMsgData("parsing file %q into type %T (%d bytes)", d.Name, into, len(d.Data)))

		err := yaml.Unmarshal(d.Data, into)
		if err != nil {
			return nil, fmt.Errorf("error while parsing file %q into type %T: %w", d.Name, into, err)
		}

		result = append(result, into)
	}

	return result, nil
}
