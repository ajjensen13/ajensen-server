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
}

type webProject struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	ContentHtml string     `json:"contentHtml"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Tags        []string   `json:"tags"`
	Parent      string     `json:"parent"`
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
	}
}
