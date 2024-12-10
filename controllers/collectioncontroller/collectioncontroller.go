package collectioncontroller

import (
	"Authors/entities"
	"Authors/models/bookmodel"
	"net/http"
	"strconv"
	"time"
	"text/template"

	"github.com/gin-gonic/gin"
)

// Index displays a list of books in the collection
func Index(c *gin.Context) {
	books := bookmodel.GetAll()
	data := map[string]any{
		"books": books,
	}
	temp, err := template.ParseFiles("views/collection/index.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	temp.Execute(c.Writer, data)
}

// Show handles displaying details of a single book
func Show(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book := bookmodel.Detail(id)
	if book == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data := map[string]any{
		"book": book,
	}
	temp, err := template.ParseFiles("views/collection/show.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	temp.Execute(c.Writer, data)
}

// Add handles adding new books to the collection
func Add(c *gin.Context) {
	if c.Request.Method == "GET" {
		temp, err := template.ParseFiles("views/collection/create.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		temp.Execute(c.Writer, nil)
		return
	}

	if c.Request.Method == "POST" {
		var book entities.Book

		book.Title = c.PostForm("title")
		authorName := c.PostForm("author_name")

		// Find the author by name
		author, err := bookmodel.FindAuthorByName(authorName)
		if err != nil || author == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		book.Author = *author // Set the author field directly
		book.Genre = c.PostForm("genre")
		book.Description = c.PostForm("synopsis")

		// Convert release_date from string to time.Time
		releaseDateStr := c.PostForm("release_date")
		releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		book.ReleaseDate = releaseDate

		book.ImageURL = c.PostForm("image_url")

		// Create the book entry
		if ok := bookmodel.Create(book); !ok {
			temp, err := template.ParseFiles("views/collection/create.html")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			temp.Execute(c.Writer, nil)
			return
		}

		c.Redirect(http.StatusSeeOther, "/books")
	}
}


// Edit handles editing existing book information
func Edit(c *gin.Context) {
	if c.Request.Method == "GET" {
		idString := c.Query("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		book := bookmodel.Detail(id)
		if book == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Prepare the data for the template rendering
		data := map[string]any{
			"book": book,
		}

		temp, err := template.ParseFiles("views/collection/edit.html")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		temp.Execute(c.Writer, data)
		return
	}

	if c.Request.Method == "POST" {
		var book entities.Book

		idString := c.PostForm("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Populate the book fields
		book.Id = book.Id
		book.Title = c.PostForm("title")
		authorName := c.PostForm("author_name")

		// Find the author by name
		author, err := bookmodel.FindAuthorByName(authorName)
		if err != nil || author == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		book.Author = *author // Set the author field directly
		book.Genre = c.PostForm("genre")
		book.Description = c.PostForm("synopsis")

		// Convert release_date from string to time.Time
		releaseDateStr := c.PostForm("release_date")
		releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		book.ReleaseDate = releaseDate

		book.ImageURL = c.PostForm("image_url")

		if ok := bookmodel.Update(id, book); !ok {
			c.Redirect(http.StatusSeeOther, c.Request.Referer())
			return
		}

		c.Redirect(http.StatusSeeOther, "/books")
	}
}


// Delete handles deleting a book from the collection
func Delete(c *gin.Context) {
	idString := c.Query("id")
	id, err := strconv.Atoi(idString)
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
