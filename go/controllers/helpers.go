package controllers

import (
	"fmt"
	"math/rand/v2"
	"sort"

	"bookclubbot.com/main/models"
)

// this layer does not mutate the underlying state objects, it operates without side effects.
func selectNextBook(books []models.BookEntry) (models.BookEntry, error) {
	if len(books) == 0 {
		return models.BookEntry{}, fmt.Errorf("No books provided.")
	}
	valid_books := []models.BookEntry{}

	for _, book := range books {
		if book.Read {
			continue
		}
		valid_books = append(valid_books, book)
	}

	if len(valid_books) == 0 {
		return models.BookEntry{}, fmt.Errorf("No valid books provided.")
	}

	sort.Slice(valid_books, func(i, j int) bool {
		return valid_books[i].Votes > valid_books[j].Votes
	})

	upper_bound := min(5, len(valid_books))
	i := rand.IntN(upper_bound)

	return valid_books[i], nil
}
