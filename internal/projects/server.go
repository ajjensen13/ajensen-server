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

package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ajjensen13/ajensen-server/internal"
)

var (
	lock        sync.RWMutex
	webProjects []*webProject
)

func initWebProjects(l *log.Logger, dir string) error {
	lock.Lock()
	defer lock.Unlock()

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(l, ds, new([]*dataProject))
	if err != nil {
		return err
	}

	webProjects = transformFileData(l, is)
	return nil
}

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing projects from directory: %s", dir)

	err := initWebProjects(l, dir)
	if err != nil {
		return err
	}

	r.GET("/projects", projectHandler(l, dir))

	return nil
}

func projectHandler(l *log.Logger, dir string) func(*gin.Context) {
	result := func(c *gin.Context) {
		lock.RLock()
		defer lock.RUnlock()

		c.JSON(http.StatusOK, webProjects)
	}

	if gin.Mode() == gin.DebugMode {
		return func(c *gin.Context) {
			l.Printf("reloading projects from directory because we're in debug mode: %s", dir)
			err := initWebProjects(l, dir)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			result(c)
		}
	}

	return result
}

func transformFileData(l *log.Logger, is []interface{}) []*webProject {
	ws := make([]*webProject, 0, len(is))
	for _, i := range is {
		ds := i.(*[]*dataProject)
		for _, d := range *ds {
			w := d.webProject()
			l.Printf("transformed: %s (%T -> %T)", w.Id, d, w)
			ws = append(ws, w)
		}
	}
	return ws
}

type dataProject struct {
	Id              string     `yaml:"id"`
	Title           string     `yaml:"title"`
	ContentMarkdown string     `yaml:"contentMarkdown"`
	StartDate       time.Time  `yaml:"startDate"`
	EndDate         *time.Time `yaml:"endDate"`
	Tags            []string   `yaml:"tags"`
	Parent          string     `yaml:"parent"`
	Color           string     `yaml:"color"`
}

type webProject struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	ContentHtml string     `json:"contentHtml"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Parent      string     `json:"parent,omitempty"`
	Color       string     `json:"color,omitempty"`
}

func (d *dataProject) webProject() *webProject {
	return &webProject{
		Id:          d.Id,
		Title:       d.Title,
		ContentHtml: string(blackfriday.Run([]byte(d.ContentMarkdown))),
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Tags:        d.Tags,
		Parent:      d.Parent,
		Color:       d.Color,
	}
}
