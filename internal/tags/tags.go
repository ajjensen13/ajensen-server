package tags

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/ajjensen13/ajensen-server/internal"
)

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing tags from directory: %s", dir)

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	is, err := internal.ParseFileData(l, ds, new([]*dataTag))
	if err != nil {
		return err
	}

	ws := transformFileData(l, is)

	r.GET("/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, ws)
	})

	return nil
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
}

type webTag struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Hyperlink string `json:"hyperlink"`
}

func (d *dataTag) webTag() *webTag {
	return &webTag{
		Id:        d.Id,
		Title:     d.Title,
		Hyperlink: d.Hyperlink,
	}
}
