package homecontroller

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	temp, err := template.ParseFiles("views/index.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	temp.Execute(c.Writer, nil)
}
