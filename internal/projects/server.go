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
	"github.com/ajjensen13/gke"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"net/http"
	"sync"
	"time"

	"github.com/ajjensen13/ajensen-server/internal"
)

var (
	lock        sync.RWMutex
	webProjects []*webProject
)

func initWebProjects(lg gke.Logger, dir string) error {
	lock.Lock()
	defer lock.Unlock()

	ds, err := internal.LoadDirData(lg, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(lg, ds, new([]*dataProject))
	if err != nil {
		return err
	}

	webProjects = transformFileData(lg, is)
	return nil
}

// Init binds routes to r for serving projects.
func Init(lg gke.Logger, r gin.IRoutes, dir string) error {
	lg.Default(gke.NewFmtMsgData("initializing projects from directory: %s", dir))

	err := initWebProjects(lg, dir)
	if err != nil {
		return err
	}

	r.GET("/projects", projectHandler(lg, dir))

	return nil
}

func projectHandler(lg gke.Logger, dir string) func(*gin.Context) {
	result := func(c *gin.Context) {
		lock.RLock()
		defer lock.RUnlock()

		c.JSON(http.StatusOK, webProjects)
	}

	if gin.Mode() == gin.DebugMode {
		return func(c *gin.Context) {
			lg.Default(gke.NewFmtMsgData("reloading projects from directory because we're in debug mode: %s", dir))
			err := initWebProjects(lg, dir)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			result(c)
		}
	}

	return result
}

func transformFileData(lg gke.Logger, is []interface{}) []*webProject {
	ws := make([]*webProject, 0, len(is))
	for _, i := range is {
		ds := i.(*[]*dataProject)
		for _, d := range *ds {
			w := d.webProject()
			lg.Default(gke.NewFmtMsgData("transformed: %s (%T -> %T)", w.Id, d, w))
			ws = append(ws, w)
		}
	}
	return ws
}

type dataProject struct {
	Id              string     `yaml:"id"`
	Title           string     `yaml:"title"`
	Summary         string     `yaml:"summary"`
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
	Summary     string     `json:"summary"`
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
		Summary:     d.Summary,
		ContentHtml: string(blackfriday.Run([]byte(d.ContentMarkdown))),
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Tags:        d.Tags,
		Parent:      d.Parent,
		Color:       d.Color,
	}
}
