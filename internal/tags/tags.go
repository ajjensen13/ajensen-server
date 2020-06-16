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

package tags

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"

	"github.com/ajjensen13/ajensen-server/internal"
)

var (
	lock    sync.RWMutex
	webTags []*webTag
)

func initWebTags(l *log.Logger, dir string) error {
	lock.Lock()
	defer lock.Unlock()

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(l, ds, new([]*dataTag))
	if err != nil {
		return err
	}

	webTags = transformFileData(l, is)
	return nil
}

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing tags from directory: %s", dir)
	err := initWebTags(l, dir)
	if err != nil {
		return err
	}

	r.GET("/tags", tagHandler(l, dir))

	return nil
}

func tagHandler(l *log.Logger, dir string) func(*gin.Context) {
	result := func(c *gin.Context) {
		lock.RLock()
		defer lock.RUnlock()

		c.JSON(http.StatusOK, webTags)
	}

	if gin.Mode() == gin.DebugMode {
		return func(c *gin.Context) {
			l.Printf("reloading tags from directory because we're in debug mode: %s", dir)
			err := initWebTags(l, dir)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			result(c)
		}
	}

	return result
}

func transformFileData(l *log.Logger, is []interface{}) []*webTag {
	ws := make([]*webTag, 0, len(is))
	for _, i := range is {
		ds := i.(*[]*dataTag)
		for _, d := range *ds {
			w := d.webTag()
			l.Printf("transformed: %s (%T -> %T)", w.Id, d, w)
			ws = append(ws, w)
		}
	}
	return ws
}

type dataTag struct {
	Id        string `yaml:"id"`
	Title     string `yaml:"title"`
	Hyperlink string `yaml:"hyperlink"`
	Color     string `yaml:"color"`
}

type webTag struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Hyperlink string `json:"hyperlink"`
	Color     string `json:"color"`
}

func (d *dataTag) webTag() *webTag {
	return &webTag{
		Id:        d.Id,
		Title:     d.Title,
		Hyperlink: d.Hyperlink,
		Color:     d.Color,
	}
}
