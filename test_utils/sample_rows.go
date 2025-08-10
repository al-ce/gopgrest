package test_utils

import (
	"log"

	"gopgrest/repository"
)

func GetSampleRows(repo repository.Repository) SampleRows {
	authorIDs := insertAuthors(repo)
	bookIDs := insertBooks(repo, authorIDs)
	authors := getSampleAuthors(repo, authorIDs)
	books := getSampleBooks(repo, bookIDs)
	return SampleRows{Authors: authors, Books: books}
}

func insertAuthors(repo repository.Repository) []int64 {
	authorIDs := []int64{}

	// Insert authors
	insertedAuthorRows, err := repo.DB.Query(`INSERT INTO authors (surname, forename)
		VALUES
			($1, $2),
			($3, $4),
			($5, $6)
		RETURNING id`,
		"Brontë", "Anne",
		"Carson", "Anne",
		"Woolf", "Virginia",
	)
	if err != nil {
		log.Fatalf("Sample Author Insert err: %v\n", err)
	}
	defer insertedAuthorRows.Close()
	// Scan ids into array
	for insertedAuthorRows.Next() {
		var id int64
		insertedAuthorRows.Scan(&id)
		authorIDs = append(authorIDs, id)
	}
	return authorIDs
}

func insertBooks(repo repository.Repository, authorIDs []int64) []int64 {
	bookIDs := []int64{}
	// Insert books
	insertedBookRows, err := repo.DB.Query(`INSERT INTO books (title, author_id)
		VALUES
			($1, $2),
			($3, $4),
			($5, $6)
		RETURNING id`,
		"The Tenant of Wildfell Hall", authorIDs[0],
		"Autobiography of Red", authorIDs[1],
		"Mrs. Dalloway", authorIDs[2],
	)
	if err != nil {
		log.Fatalf("Sample Book Insert err: %v\n", err)
	}
	defer insertedBookRows.Close()

	// Scan ids into array
	for insertedBookRows.Next() {
		var id int64
		insertedBookRows.Scan(&id)
		bookIDs = append(bookIDs, id)
	}
	return bookIDs
}

func getSampleAuthors(repo repository.Repository, authorIDs []int64) SampleAuthorsMap {
	authors := SampleAuthorsMap{}
	for _, id := range authorIDs {
		row := repo.DB.QueryRow("SELECT * FROM authors WHERE id=$1", id)
		author, err := ScanAuthor(row)
		if err != nil {
			log.Fatalf("Author scan err: %s\n", err)
		}
		authors[int64(id)] = author
	}
	return authors
}

func getSampleBooks(repo repository.Repository, bookIDs []int64) SampleBooksMap {
	books := SampleBooksMap{}
	for _, id := range bookIDs {
		row := repo.DB.QueryRow("SELECT * FROM books WHERE id=$1", id)
		book := SampleBook{}
		err := row.Scan(&book.ID, &book.Title, &book.AuthorID)
		if err != nil {
			log.Fatalf("Book scan err: %s\n", err)
		}
		books[id] = book
	}
	return books
}
