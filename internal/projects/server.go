package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"log"
	"net/http"
	"time"

	"github.com/ajjensen13/ajensen-server/internal"
)

type project struct {
	Id              uint64     `json:"id" yaml:"id"`
	Title           string     `json:"title" yaml:"title"`
	ContentHtml     string     `json:"content_html" yaml:"-"`
	ContentMarkdown string     `json:"-" yaml:"content_markdown"`
	Tags            []uint64   `json:"tags" yaml:"tags"`
	StartDate       time.Time  `json:"start_date" yaml:"start_date"`
	EndDate         *time.Time `json:"end_date" yaml:"end_date"`
	Children        []*project `json:"children" yaml:"children"`
}

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing projects from directory: %s", dir)

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseData(l, ds, new(project))
	if err != nil {
		return err
	}

	for _, i := range is {
		p := i.(*project)
		p.initHtmlFromMarkdown(l)
	}

	r.GET("/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, is)
	})

	return nil
}

func (p *project) initHtmlFromMarkdown(l *log.Logger) {
	p.ContentHtml = string(blackfriday.Run([]byte(p.ContentMarkdown)))
	l.Printf("parsed markdown: %s [%d] (%d bytes)", p.Title, p.Id, len(p.ContentHtml))

	for _, c := range p.Children {
		c.initHtmlFromMarkdown(l)
	}
}
