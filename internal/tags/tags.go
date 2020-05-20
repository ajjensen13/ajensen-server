package tags

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/ajjensen13/ajensen-server/internal"
)

type tag struct {
	Id    uint64 `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`
}

func Init(l *log.Logger, r gin.IRoutes, dir string) error {
	l.Printf("initializing tags from directory: %s", dir)

	ds, err := internal.LoadFileData(l, dir)
	if err != nil {
		return err
	}

	ps, err := internal.ParseData(l, ds, new(tag))
	if err != nil {
		return err
	}

	r.GET("/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, ps)
	})

	return nil
}
