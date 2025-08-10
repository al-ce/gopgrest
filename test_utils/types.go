package test_utils

import (
	"database/sql"

	"gopgrest/repository"
)

type TestDB struct {
	DB     *sql.DB
	TX     *sql.Tx
	Tables []repository.Table
}

// SampleAuthor matches the authors table in the test database so that rows
// can be scanned into fields of appropriate size
type SampleAuthor struct {
	ID       int64  `json:"id"`
	Surname  string `json:"surname"`
	Forename string `json:"forename"`
}

// SampleBook matches the books table in the test database so that rows
// can be scanned into fields of appropriate size
type SampleBook struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
}

// SampleAuthorsMap is a map of authors by ID
type SampleAuthorsMap map[int64]SampleAuthor

// SampleBooksMap is a map of books by ID
type SampleBooksMap map[int64]SampleBook

// SampleRows holds arrays of inserted sample Author and Book structs
type SampleRows struct {
	Authors SampleAuthorsMap
	Books   SampleBooksMap
}
