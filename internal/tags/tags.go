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
	"github.com/ajjensen13/gke"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"

	"github.com/ajjensen13/ajensen-server/internal"
)

var (
	lock    sync.RWMutex
	webTags []*webTag
)

func initWebTags(lg gke.Logger, dir string) error {
	lock.Lock()
	defer lock.Unlock()

	ds, err := internal.LoadDirData(lg, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(lg, ds, new([]*dataTag))
	if err != nil {
		return err
	}

	webTags = transformFileData(lg, is)
	return nil
}

// Init binds routes to r for serving tags.
func Init(lg gke.Logger, r gin.IRoutes, dir string) error {
	lg.Default(gke.NewFmtMsgData("initializing tags from directory: %s", dir))
	err := initWebTags(lg, dir)
	if err != nil {
		return err
	}

	r.GET("/tags", tagHandler(lg, dir))

	return nil
}

func tagHandler(lg gke.Logger, dir string) func(*gin.Context) {
	result := func(c *gin.Context) {
		lock.RLock()
		defer lock.RUnlock()

		c.JSON(http.StatusOK, webTags)
	}

	if gin.Mode() == gin.DebugMode {
		return func(c *gin.Context) {
			lg.Default(gke.NewFmtMsgData("reloading tags from directory because we're in debug mode: %s", dir))
			err := initWebTags(lg, dir)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			result(c)
		}
	}

	return result
}

func transformFileData(lg gke.Logger, is []interface{}) []*webTag {
	ws := make([]*webTag, 0, len(is))
	for _, i := range is {
		ds := i.(*[]*dataTag)
		for _, d := range *ds {
			w := d.webTag()
			lg.Default(gke.NewFmtMsgData("transformed: %s (%T -> %T)", w.Id, d, w))
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
	Color     string `json:"color,omitempty"`
}

func (d *dataTag) webTag() *webTag {
	return &webTag{
		Id:        d.Id,
		Title:     d.Title,
		Hyperlink: d.Hyperlink,
		Color:     d.Color,
	}
}
