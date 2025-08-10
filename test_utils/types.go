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

// Authors matches the authors table in the test database so that rows
// can be scanned into fields of appropriate size
type Authors struct {
	ID       int    `json:"id"`
	Surname  string `json:"surname"`
	Forename string `json:"forename"`
}

// books matches the books table in the test database so that rows
// can be scanned into fields of appropriate size
type Books struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
}
