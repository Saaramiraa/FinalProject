package authormodel

import (
	"Authors/config"
	"Authors/entities"
)

func GetAll() []entities.Author {
	rows, err := config.DB.Query(`SELECT * FROM authors`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var authors []entities.Author

	for rows.Next() {
		var author entities.Author
		if err := rows.Scan(&author.Id, &author.Name, &author.DoB, &author.Updated_At); err != nil {
			panic(err)
		}

		authors = append(authors, author)
	}
	return authors
}

func Create(author entities.Author) bool {
	result, err := config.DB.Exec(`
	INSERT INTO authors (name, date_of_birth)
	VALUE (?, ?)`,
		author.Name,
		author.DoB,
	)
	if err != nil {
		panic(err)
	}
	LastInsertedId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return LastInsertedId > 0
}

func Detail(id int) entities.Author {
	row := config.DB.QueryRow(`SELECT id, name, date_of_birth FROM authors WHERE id = ?`, id)

	var author entities.Author
	if err := row.Scan(&author.Id, &author.Name, &author.DoB); err != nil {
		panic(err.Error())
	}
	return author
}

func Update(id int, author entities.Author) bool {
	query, err := config.DB.Exec(`UPDATE authors SET name = ?, date_of_birth = ? WHERE id = ?`, author.Name, author.DoB, id)
	if err != nil {
		panic(err)
	}

	result, err := query.RowsAffected()
	if err != nil {
		panic(err)
	}

	return result > 0
}

func Delete(id int) error {
	_, err := config.DB.Exec(`DELETE FROM authors WHERE id = ?`, id)
	return err
}
