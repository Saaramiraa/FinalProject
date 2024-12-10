package authorcontroller

import (
	"Authors/entities"
	"Authors/models/authormodel"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

// Index displays a list of authors
func Index(c *gin.Context) {
	authors := authormodel.GetAll()
	data := map[string]any{
		"authors": authors,
	}
	temp, err := template.ParseFiles("views/author/index.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	temp.Execute(c.Writer, data)
}

// Edit handles author editing
func Edit(c *gin.Context) {
	if c.Request.Method == "GET" {
		temp, err := template.ParseFiles("views/author/edit.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		idString := c.Query("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		author := authormodel.Detail(id)
		data := map[string]any{
			"author": author,
		}

		temp.Execute(c.Writer, data)
		return
	}

	if c.Request.Method == "POST" {
		var author entities.Author

		idString := c.PostForm("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		author.Name = c.PostForm("name")
		author.DoB = c.PostForm("DoB")

		if ok := authormodel.Update(id, author); !ok {
			c.Redirect(http.StatusSeeOther, c.Request.Referer())
			return
		}

		c.Redirect(http.StatusSeeOther, "/authors")
	}
}

// Add handles adding new authors
func Add(c *gin.Context) {
	if c.Request.Method == "GET" {
		temp, err := template.ParseFiles("views/author/create.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		temp.Execute(c.Writer, nil)
		return
	}

	if c.Request.Method == "POST" {
		var author entities.Author

		author.Name = c.PostForm("name")
		author.DoB = c.PostForm("DoB")

		if ok := authormodel.Create(author); !ok {
			temp, _ := template.ParseFiles("views/author/create.html")
			temp.Execute(c.Writer, nil)
			return
		}

		c.Redirect(http.StatusSeeOther, "/authors")
	}
}

// Delete handles deleting an author
func Delete(c *gin.Context) {
	idString := c.Query("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := authormodel.Delete(id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/authors")
}
