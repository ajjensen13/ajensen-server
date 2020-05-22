package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"log"
	"net/http"
	"time"

	"github.com/ajjensen13/ajensen-server/internal"
)

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing projects from directory: %s", dir)

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(l, ds, new([]*dataProject))
	if err != nil {
		return err
	}

	ws := transformFileData(l, is)

	r.GET("/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, ws)
	})

	return nil
}

func transformFileData(l *log.Logger, is []interface{}) []*webProject {
	ws := make([]*webProject, 0, len(is))
	for _, i := range is {
		ds := i.(*[]*dataProject)
		for _, d := range *ds {
			w := d.webProject()
			l.Printf("parsed markdown: %s (%d bytes)", w.Id, len(w.ContentHtml))
			ws = append(ws, w)
		}
	}
	return ws
}

type dataProject struct {
	Id              string     `yaml:"id"`
	Title           string     `yaml:"title"`
	ContentMarkdown string     `yaml:"content_markdown"`
	Tags            []uint64   `yaml:"tags"`
	StartDate       time.Time  `yaml:"start_date"`
	EndDate         *time.Time `yaml:"end_date"`
	Parent          string     `yaml:"parent"`
}

type webProject struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	ContentHtml string     `json:"content_html"`
	Tags        []uint64   `json:"tags"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Parent      string     `json:"parent"`
}

func (d *dataProject) webProject() *webProject {
	var result webProject
	result.ContentHtml = string(blackfriday.Run([]byte(d.ContentMarkdown)))
	return &result
}
