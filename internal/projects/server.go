package projects

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/ajjensen13/ajensen-server/internal"
)

type project struct {
	Id              uint64     `json:"id" yaml:"id"`
	Title           string     `json:"title" yaml:"title"`
	ContentHtml     string     `json:"content_html" yaml:"content_html"`
	ContentMarkdown string     `json:"content_markdown" yaml:"content_markdown"`
	Tags            []uint64   `json:"tags" yaml:"tags"`
	StartDate       time.Time  `json:"start_date" yaml:"start_date"`
	EndDate         *time.Time `json:"end_date" yaml:"end_date"`
	Children        []project  `json:"children" yaml:"children"`
}

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing projects from directory: %s", dir)

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	ps, err := internal.ParseData(l, ds, new(project))
	if err != nil {
		return err
	}

	r.GET("/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, ps)
	})

	return nil
}
