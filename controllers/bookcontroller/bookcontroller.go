package bookcontroller

import (
	"Authors/entities"
	"Authors/models/authormodel"
	"Authors/models/bookmodel"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	books := bookmodel.GetAll()
	data := map[string]any{
		"books": books,
	}
	temp, err := template.ParseFiles("views/book/view.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	temp.Execute(c.Writer, data)
}

func Add(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/book/create.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		authors := authormodel.GetAll()
		data := map[string]any{
			"authors": authors,
		}
		temp.Execute(c.Writer, data)
		return
	}

	if c.Request.Method == http.MethodPost {
		var book entities.Book
		authorId, err := strconv.Atoi(c.PostForm("author_id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		book.Title = c.PostForm("title")
		book.Author.Id = authorId
		book.Genre = c.PostForm("genre")
		book.Description = c.PostForm("description")
		book.Updated_At = time.Now()
		book.Added_At = time.Now()

		if ok := bookmodel.Create(book); !ok {
			c.Redirect(http.StatusTemporaryRedirect, c.Request.Referer())
			return
		}
		c.Redirect(http.StatusSeeOther, "/books")
	}
}

func Edit(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/book/edit.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		book := bookmodel.Detail(id)
		authors := authormodel.GetAll()
		data := map[string]any{
			"authors": authors,
			"book":    book,
		}
		temp.Execute(c.Writer, data)
		return
	}

	if c.Request.Method == http.MethodPost {
		var book entities.Book

		id, err := strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		authorId, err := strconv.Atoi(c.PostForm("author_id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		book.Title = c.PostForm("title")
		book.Author.Id = authorId
		book.Genre = c.PostForm("genre")
		book.Description = c.PostForm("description")
		book.Updated_At = time.Now()

		if ok := bookmodel.Update(id, book); !ok {
			c.Redirect(http.StatusTemporaryRedirect, c.Request.Referer())
			return
		}
		c.Redirect(http.StatusSeeOther, "/books")
	}
}

func Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book := bookmodel.Detail(id)
	data := map[string]any{
		"book": book,
	}

	temp, err := template.ParseFiles("views/book/detail.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	temp.Execute(c.Writer, data)
}

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := bookmodel.Delete(id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/books")
}
