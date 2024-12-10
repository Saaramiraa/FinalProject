package bookmodel

import (
	"Authors/config"
	"Authors/entities"
	"errors"
)

// GetAll retrieves all books from the database
func GetAll() []entities.Book {
	rows, err := config.DB.Query(`
		SELECT 
			books.id,
			books.title,
			authors.name AS author_name,
			books.description,
			books.genre,
			books.updated_at,
			books.added_at
		FROM books
		JOIN authors ON books.author_id = authors.id
	`)
	if err != nil {
		// Log and return an empty slice
		panic(err)
	}
	defer rows.Close()

	var books []entities.Book
	for rows.Next() {
		var book entities.Book
		if err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.Author.Name,
			&book.Description,
			&book.Genre,
			&book.Updated_At,
			&book.Added_At,
		); err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return books
}

// Create adds a new book to the database
func Create(book entities.Book) bool {
	result, err := config.DB.Exec(`
		INSERT INTO books (
			title, author_id, genre, description, updated_at, added_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		book.Title,
		book.Author.Id,
		book.Genre,
		book.Description,
		book.Updated_At,
		book.Added_At,
	)
	if err != nil {
		panic(err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return lastInsertedId > 0
}

// Detail retrieves the details of a specific book by ID
func Detail(id int) *entities.Book {
	row := config.DB.QueryRow(`
		SELECT 
			books.id,
			books.title,
			authors.name AS author_name,
			books.description,
			books.genre,
			books.updated_at,
			books.added_at
		FROM books
		JOIN authors ON books.author_id = authors.id
		WHERE books.id = ?`,
		id,
	)

	var book entities.Book
	if err := row.Scan(
		&book.Id,
		&book.Title,
		&book.Author.Name,
		&book.Description,
		&book.Genre,
		&book.Updated_At,
		&book.Added_At,
	); err != nil {
		// If no record is found, return nil
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		panic(err) // Handle other errors
	}

	return &book // Return a pointer to the book
}


// Update modifies an existing book's details
func Update(id int, book entities.Book) bool {
	query, err := config.DB.Exec(`
		UPDATE books SET 
			title = ?, 
			author_id = ?, 
			genre = ?,
			description = ?,
			updated_at = ? 
		WHERE id = ?`,
		book.Title,
		book.Author.Id,
		book.Genre,
		book.Description,
		book.Updated_At,
		id,
	)
	if err != nil {
		panic(err)
	}

	result, err := query.RowsAffected()
	if err != nil {
		panic(err)
	}

	return result > 0
}

// FindAuthorByName searches for an author by name in the database
func FindAuthorByName(name string) (*entities.Author, error) {
	var author entities.Author
	err := config.DB.QueryRow(`
		SELECT id, name FROM authors WHERE name = ?`,
		name,
	).Scan(&author.Id, &author.Name)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("author not found")
		}
		return nil, err
	}
	return &author, nil
}

// Delete removes a book by ID
func Delete(id int) error {
	_, err := config.DB.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}
