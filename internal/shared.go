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

func LoadFileData(lg gke.Logger, dir string) ([]FileData, error) {
	lg.Defaultf("loading directory: %s", dir)

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
			ps, err := LoadFileData(lg, f)
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
		lg.Defaultf("loaded file: %s (%d bytes)", f, len(f))
	}

	return result, nil
}

func ParseFileData(lg gke.Logger, ds []FileData, i interface{}) ([]interface{}, error) {
	var result []interface{}

	t := reflect.TypeOf(i).Elem()
	for _, d := range ds {
		into := reflect.New(t).Interface()
		lg.Defaultf("parsing file %q into type %T (%d bytes)", d.Name, into, len(d.Data))

		err := yaml.Unmarshal(d.Data, into)
		if err != nil {
			return nil, fmt.Errorf("internal: error while parsing file %q into type %T: %w", d.Name, into, err)
		}

		result = append(result, into)
	}

	return result, nil
}
